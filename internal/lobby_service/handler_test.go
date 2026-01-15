package lobby_service

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mygameplatform/internal/auth"
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

func TestGames_RequiresAuth(t *testing.T) {
	logger := slog.New(&captureHandler{})
	verifier, err := auth.NewHS256Signer([]byte("test-secret"), nil)
	if err != nil {
		t.Fatalf("NewHS256Signer: %v", err)
	}

	h := NewHandler(verifier, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/games", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"error":"forbidden"`) {
		t.Fatalf("expected forbidden error code, got %s", rec.Body.String())
	}
}
