package lobby_service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mygameplatform/internal/auth"
)

func TestGames_ReturnsCatalog(t *testing.T) {
	logger := slog.New(&captureHandler{})
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}

	token, err := signer.Sign("user-123", time.Now().UTC())
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	h := NewHandler(signer, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/games", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	gamesRaw, ok := body["games"]
	if !ok {
		t.Fatalf("expected games key in response")
	}
	games, ok := gamesRaw.([]any)
	if !ok || len(games) < 2 {
		t.Fatalf("expected at least 2 games, got %T (%d)", gamesRaw, len(games))
	}

	found := map[string]bool{}
	for _, item := range games {
		game, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("expected game object, got %T", item)
		}
		id, _ := game["game_id"].(string)
		if id != "" {
			found[id] = true
		}
		if _, ok := game["game_id"]; !ok {
			t.Fatalf("expected game_id field")
		}
		if _, ok := game["min_players"]; !ok {
			t.Fatalf("expected min_players field")
		}
		if _, ok := game["max_players"]; !ok {
			t.Fatalf("expected max_players field")
		}
		if _, ok := game["rules_summary"]; !ok {
			t.Fatalf("expected rules_summary field")
		}
	}

	if !found["connect4"] {
		t.Fatalf("expected connect4 in games list")
	}
	if !found["draughts_10x10"] {
		t.Fatalf("expected draughts_10x10 in games list")
	}
}

func TestGames_DoesNotLogAuthorizationHeaderOrTokenizedURL(t *testing.T) {
	logSink := &captureHandler{}
	logger := slog.New(logSink)
	signer, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}

	token, err := signer.Sign("user-123", time.Now().UTC())
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	h := NewHandler(signer, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/games?access_token=leak-me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	logged := strings.Join(logSink.records, "\n")
	if strings.Contains(logged, token) {
		t.Fatalf("log leak: contains jwt")
	}
	if strings.Contains(logged, "leak-me") {
		t.Fatalf("log leak: contains tokenized URL")
	}
}
