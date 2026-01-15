package lobby_service

import (
	"log/slog"
	"net/http"
	"time"

	"mygameplatform/internal/errorsx"
)

type TokenVerifier interface {
	Verify(token string, now time.Time) (string, error)
}

type Handler struct {
	verifier TokenVerifier
	logger   *slog.Logger
}

func NewHandler(verifier TokenVerifier, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{verifier: verifier, logger: logger}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/api/v1/games" {
		h.authMiddleware(http.HandlerFunc(h.handleGames)).ServeHTTP(w, r)
		return
	}

	errorsx.Write(w, http.StatusNotFound, "not_found", "Not found")
}
