package integration_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mygameplatform/internal/auth"
	"mygameplatform/internal/lobby_service"
)

type nullHandler struct{}

func (nullHandler) Enabled(_ context.Context, _ slog.Level) bool { return false }
func (nullHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (nullHandler) WithAttrs(_ []slog.Attr) slog.Handler          { return nullHandler{} }
func (nullHandler) WithGroup(string) slog.Handler                { return nullHandler{} }

func TestLobbyServiceGames_Integration(t *testing.T) {
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}
	logger := slog.New(nullHandler{})
	h := lobby_service.NewHandler(signer, logger)

	token, err := signer.Sign("user-123", time.Now().UTC())
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/games", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Games []struct {
			GameID string `json:"game_id"`
		} `json:"games"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(body.Games) < 2 {
		t.Fatalf("expected at least 2 games, got %d", len(body.Games))
	}
}
