package auth_service

import (
	"log/slog"
	"net/http"
	"time"

	"mygameplatform/internal/db"
	"mygameplatform/internal/errorsx"
)

type TokenSigner interface {
	Sign(sub string, now time.Time) (string, error)
}

type Handler struct {
	users  db.UserStore
	signer TokenSigner
	logger *slog.Logger
}

func NewHandler(users db.UserStore, signer TokenSigner, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{users: users, signer: signer, logger: logger}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/signup" {
		h.handleSignup(w, r)
		return
	}
	errorsx.Write(w, http.StatusNotFound, "not_found", "Not found")
}

