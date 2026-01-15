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

type TokenVerifier interface {
	Verify(token string, now time.Time) (string, error)
}

type Handler struct {
	users             db.UserStore
	signer            TokenSigner
	verifier          TokenVerifier
	logger            *slog.Logger
	limiter           *rateLimiter
	trustProxyHeaders bool
}

func NewHandler(users db.UserStore, signer TokenSigner, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	var verifier TokenVerifier
	if v, ok := signer.(TokenVerifier); ok {
		verifier = v
	}
	return &Handler{
		users:    users,
		signer:   signer,
		verifier: verifier,
		logger:   logger,
		limiter:  newRateLimiter(10, time.Minute),
	}
}

func (h *Handler) SetLimiter(limiter *rateLimiter) {
	if limiter != nil {
		h.limiter = limiter
	}
}

func (h *Handler) SetVerifier(verifier TokenVerifier) {
	if verifier != nil {
		h.verifier = verifier
	}
}

func (h *Handler) SetTrustProxyHeaders(v bool) {
	h.trustProxyHeaders = v
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/signup" {
		h.handleSignup(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/login" {
		h.handleLogin(w, r)
		return
	}
	if r.URL.Path == "/api/v1/me" {
		if r.Method == http.MethodGet {
			h.authMiddleware(http.HandlerFunc(h.handleMeGet)).ServeHTTP(w, r)
			return
		}
		if r.Method == http.MethodPatch {
			h.authMiddleware(http.HandlerFunc(h.handleMePatch)).ServeHTTP(w, r)
			return
		}
	}
	errorsx.Write(w, http.StatusNotFound, "not_found", "Not found")
}
