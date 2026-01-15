package auth_service

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"mygameplatform/internal/auth"
	"mygameplatform/internal/db"
)

type captureHandler struct {
	records []string
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool { return true }
func (h *captureHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	b.WriteString(r.Message)
	r.Attrs(func(a slog.Attr) bool {
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString("=")
		b.WriteString(a.Value.String())
		return true
	})
	h.records = append(h.records, b.String())
	return nil
}
func (h *captureHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *captureHandler) WithGroup(string) slog.Handler      { return h }

type fakeUserStore struct {
	lastEmail        string
	lastPasswordHash string
	usersByEmail     map[string]userRecord
}

type userRecord struct {
	user         db.User
	passwordHash string
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{usersByEmail: map[string]userRecord{}}
}

func (s *fakeUserStore) CreateUser(_ context.Context, email, passwordHash, username string, avatarURL *string) (db.User, error) {
	s.lastEmail = email
	s.lastPasswordHash = passwordHash

	if _, ok := s.usersByEmail[email]; ok {
		return db.User{}, db.ErrEmailConflict
	}
	u := db.User{ID: "00000000-0000-0000-0000-000000000001", Email: email, Username: username, AvatarURL: avatarURL}
	s.usersByEmail[email] = userRecord{user: u, passwordHash: passwordHash}
	return u, nil
}

func (s *fakeUserStore) FindUserByEmail(_ context.Context, email string) (db.UserWithPassword, error) {
	rec, ok := s.usersByEmail[email]
	if !ok {
		return db.UserWithPassword{}, db.ErrUserNotFound
	}
	return db.UserWithPassword{ID: rec.user.ID, Email: rec.user.Email, PasswordHash: rec.passwordHash}, nil
}

func (s *fakeUserStore) GetUserByID(_ context.Context, id string) (db.User, error) {
	for _, rec := range s.usersByEmail {
		if rec.user.ID == id {
			return rec.user, nil
		}
	}
	return db.User{}, db.ErrUserNotFound
}

func (s *fakeUserStore) UpdateUserProfile(_ context.Context, id, username string, avatarURL *string) (db.User, error) {
	for key, rec := range s.usersByEmail {
		if rec.user.ID == id {
			rec.user.Username = username
			rec.user.AvatarURL = avatarURL
			s.usersByEmail[key] = rec
			return rec.user, nil
		}
	}
	return db.User{}, db.ErrUserNotFound
}

func TestSignup_Success_ReturnsTokenAndPersistsUser(t *testing.T) {
	users := newFakeUserStore()
	logSink := &captureHandler{}
	logger := slog.New(logSink)
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	h := NewHandler(users, signer, logger)

	body := map[string]any{
		"email":    "  USER@EXAMPLE.COM ",
		"password": "password123",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(b))
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

	if users.lastEmail != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", users.lastEmail)
	}
	if users.lastPasswordHash == "" || users.lastPasswordHash == body["password"] {
		t.Fatalf("expected bcrypt hash")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(users.lastPasswordHash), []byte(body["password"].(string))); err != nil {
		t.Fatalf("expected bcrypt hash to match password: %v", err)
	}

	raw := strings.Join(logSink.records, "\n")
	if strings.Contains(raw, body["password"].(string)) {
		t.Fatalf("log leak: contains password")
	}
	if strings.Contains(raw, resp.AccessToken) {
		t.Fatalf("log leak: contains jwt")
	}
}

func TestSignup_DuplicateEmail_409Conflict(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	h := NewHandler(users, signer, logger)

	_, _ = users.CreateUser(context.Background(), "user@example.com", "x", "user_example", nil)

	body := map[string]any{
		"email":    "user@example.com",
		"password": "password123",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"error":"conflict"`) {
		t.Fatalf("expected conflict error code, got %s", rec.Body.String())
	}
}

func TestSignup_InvalidInput_400ValidationFailed(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	h := NewHandler(users, signer, logger)

	body := map[string]any{
		"email":    "not-an-email",
		"password": "short",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"error":"validation_failed"`) {
		t.Fatalf("expected validation_failed error code, got %s", rec.Body.String())
	}
}

func TestSignup_TooLongPassword_400ValidationFailed(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	h := NewHandler(users, signer, logger)

	body := map[string]any{
		"email":    "user@example.com",
		"password": strings.Repeat("a", maxPasswordLength+1),
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"error":"validation_failed"`) {
		t.Fatalf("expected validation_failed error code, got %s", rec.Body.String())
	}
}

func TestSignup_InvalidContentType_400ValidationFailed(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	h := NewHandler(users, signer, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"error":"validation_failed"`) {
		t.Fatalf("expected validation_failed error code, got %s", rec.Body.String())
	}
}
