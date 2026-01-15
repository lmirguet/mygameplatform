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

func seedProfileUser(users *fakeUserStore, id, email, username string, avatarURL *string) {
	users.usersByEmail[email] = userRecord{
		user: db.User{
			ID:        id,
			Email:     email,
			Username:  username,
			AvatarURL: avatarURL,
		},
		passwordHash: "hash",
	}
}

func authHeader(token string) string {
	return "Bearer " + token
}

func TestMe_Get_Unauthorized(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	h := NewHandler(users, signer, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestMe_Get_Profile(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	userID := "00000000-0000-0000-0000-000000000777"
	avatar := "https://example.com/a.png"
	seedProfileUser(users, userID, "user@example.com", "player_one", &avatar)
	h := NewHandler(users, signer, logger)

	token, err := signer.Sign(userID, time.Now().UTC())
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Authorization", authHeader(token))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp profileResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Username != "player_one" {
		t.Fatalf("expected username player_one, got %q", resp.Username)
	}
	if resp.AvatarURL == nil || *resp.AvatarURL != avatar {
		t.Fatalf("expected avatar_url %q, got %v", avatar, resp.AvatarURL)
	}
}

func TestMe_Patch_UpdatesProfile(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	userID := "00000000-0000-0000-0000-000000000888"
	seedProfileUser(users, userID, "user@example.com", "player_one", nil)
	h := NewHandler(users, signer, logger)

	token, err := signer.Sign(userID, time.Now().UTC())
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	body := map[string]any{
		"username":   "player_two",
		"avatar_url": "https://example.com/b.png",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", bytes.NewReader(b))
	req.Header.Set("Authorization", authHeader(token))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	reqGet.Header.Set("Authorization", authHeader(token))
	recGet := httptest.NewRecorder()
	h.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recGet.Code, recGet.Body.String())
	}
	var resp profileResponse
	if err := json.Unmarshal(recGet.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Username != "player_two" {
		t.Fatalf("expected username player_two, got %q", resp.Username)
	}
	if resp.AvatarURL == nil || *resp.AvatarURL != "https://example.com/b.png" {
		t.Fatalf("expected avatar_url to be updated, got %v", resp.AvatarURL)
	}
}

func TestMe_Patch_InvalidUsername(t *testing.T) {
	users := newFakeUserStore()
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	userID := "00000000-0000-0000-0000-000000000999"
	seedProfileUser(users, userID, "user@example.com", "player_one", nil)
	h := NewHandler(users, signer, logger)

	token, err := signer.Sign(userID, time.Now().UTC())
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	body := map[string]any{"username": " "}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", bytes.NewReader(b))
	req.Header.Set("Authorization", authHeader(token))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"error":"validation_failed"`) {
		t.Fatalf("expected validation_failed error code, got %s", rec.Body.String())
	}
}
