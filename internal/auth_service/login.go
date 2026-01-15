package auth_service

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"mygameplatform/internal/db"
	"mygameplatform/internal/errorsx"
	"mygameplatform/internal/httpx"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

var dummyPasswordHash = func() []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)
	if err != nil {
		return []byte("$2a$10$invalidinvalidinvalidinvalidinvalidinvalidinval")
	}
	return hash
}()

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && !httpx.IsJSONContentType(r) {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "Content-Type must be application/json")
		return
	}

	var req loginRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "Invalid JSON body")
		return
	}

	email, err := normalizeAndValidateEmail(req.Email)
	if err != nil {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "email: invalid")
		return
	}
	if req.Password == "" {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "password: required")
		return
	}

	ip := clientIP(r, h.trustProxyHeaders)
	if !h.limiter.Allow(ip, time.Now().UTC()) {
		h.logRequest(r, http.StatusTooManyRequests, "rate_limited", nil)
		errorsx.Write(w, http.StatusTooManyRequests, "rate_limited", "Too many requests")
		return
	}

	user, err := h.users.FindUserByEmail(r.Context(), email)
	passwordHash := dummyPasswordHash
	if err == nil {
		passwordHash = []byte(user.PasswordHash)
	} else if !errors.Is(err, db.ErrUserNotFound) {
		h.logRequest(r, http.StatusInternalServerError, "internal", err)
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	if bcrypt.CompareHashAndPassword(passwordHash, []byte(req.Password)) != nil || err != nil {
		h.logRequest(r, http.StatusUnauthorized, "invalid_credentials", nil)
		errorsx.Write(w, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password")
		return
	}

	token, err := h.signer.Sign(user.ID, time.Now().UTC())
	if err != nil {
		h.logRequest(r, http.StatusInternalServerError, "internal", err)
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	h.logRequest(r, http.StatusOK, "", nil)
	httpx.WriteJSON(w, http.StatusOK, loginResponse{AccessToken: token})
}

func clientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if ip := firstForwardedFor(r.Header.Get("X-Forwarded-For")); ip != "" {
			return ip
		}
		if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" && net.ParseIP(ip) != nil {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	if r.RemoteAddr != "" {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return "unknown"
}

func firstForwardedFor(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	if len(parts) == 0 {
		return ""
	}
	ip := strings.TrimSpace(parts[0])
	if ip == "" {
		return ""
	}
	if net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}
