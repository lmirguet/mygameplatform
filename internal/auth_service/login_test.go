package auth_service

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mygameplatform/internal/auth"
	"mygameplatform/internal/db"
)

func seedUser(t *testing.T, users *fakeUserStore, email, password string) db.UserWithPassword {
	t.Helper()
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	u := db.User{ID: "00000000-0000-0000-0000-000000000099", Email: email, Username: "player_one"}
	users.usersByEmail[email] = userRecord{user: u, passwordHash: hash}
	return db.UserWithPassword{ID: u.ID, Email: u.Email, PasswordHash: hash}
}

func TestLogin_Success_ReturnsToken(t *testing.T) {
	users := newFakeUserStore()
	logSink := &captureHandler{}
	logger := slog.New(logSink)
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	_ = seedUser(t, users, "user@example.com", "password123")

	h := NewHandler(users, signer, logger)

	body := map[string]any{
		"email":    "  USER@EXAMPLE.COM ",
		"password": "password123",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatalf("expected access_token")
	}
	if parts := strings.Split(resp.AccessToken, "."); len(parts) != 3 {
		t.Fatalf("expected jwt with 3 parts, got %d", len(parts))
	}

	raw := strings.Join(logSink.records, "\n")
	if strings.Contains(raw, body["password"].(string)) {
		t.Fatalf("log leak: contains password")
	}
	if strings.Contains(raw, resp.AccessToken) {
		t.Fatalf("log leak: contains jwt")
	}
}

func TestLogin_InvalidCredentials_SameResponse(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	_ = seedUser(t, users, "user@example.com", "password123")

	h := NewHandler(users, signer, logger)
	h.limiter = newRateLimiter(100, time.Minute)

	missingBody := map[string]any{"email": "missing@example.com", "password": "password123"}
	wrongBody := map[string]any{"email": "user@example.com", "password": "wrongpassword"}

	missingReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(mustJSON(t, missingBody)))
	missingReq.RemoteAddr = "10.0.0.1:1111"
	missingRec := httptest.NewRecorder()
	h.ServeHTTP(missingRec, missingReq)

	wrongReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(mustJSON(t, wrongBody)))
	wrongReq.RemoteAddr = "10.0.0.2:2222"
	wrongRec := httptest.NewRecorder()
	h.ServeHTTP(wrongRec, wrongReq)

	if missingRec.Code != http.StatusUnauthorized || wrongRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for both, got %d and %d", missingRec.Code, wrongRec.Code)
	}
	if missingRec.Body.String() != wrongRec.Body.String() {
		t.Fatalf("expected same error response for invalid credentials")
	}
	if !strings.Contains(missingRec.Body.String(), `"error":"invalid_credentials"`) {
		t.Fatalf("expected invalid_credentials error code, got %s", missingRec.Body.String())
	}
}

func TestLogin_RateLimit_Returns429(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	_ = seedUser(t, users, "user@example.com", "password123")

	h := NewHandler(users, signer, logger)
	h.limiter = newRateLimiter(1, time.Minute)

	body := map[string]any{"email": "user@example.com", "password": "password123"}
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(mustJSON(t, body)))
	req1.RemoteAddr = "192.168.0.10:1234"
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(mustJSON(t, body)))
	req2.RemoteAddr = "192.168.0.10:9999"
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d: %s", rec2.Code, rec2.Body.String())
	}
	if !strings.Contains(rec2.Body.String(), `"error":"rate_limited"`) {
		t.Fatalf("expected rate_limited error code, got %s", rec2.Body.String())
	}
}

func TestLogin_RateLimit_TrustProxyHeadersUsesForwardedFor(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	_ = seedUser(t, users, "user@example.com", "password123")

	h := NewHandler(users, signer, logger)
	h.SetLimiter(newRateLimiter(1, time.Minute))
	h.SetTrustProxyHeaders(true)

	body := map[string]any{"email": "user@example.com", "password": "password123"}
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(mustJSON(t, body)))
	req1.RemoteAddr = "10.0.0.1:1111"
	req1.Header.Set("X-Forwarded-For", "203.0.113.10")
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(mustJSON(t, body)))
	req2.RemoteAddr = "10.0.0.2:2222"
	req2.Header.Set("X-Forwarded-For", "203.0.113.10")
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d: %s", rec2.Code, rec2.Body.String())
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return b
}
