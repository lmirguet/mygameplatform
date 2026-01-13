package auth_service

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"mygameplatform/internal/auth"
	"mygameplatform/internal/db"
	"mygameplatform/internal/errorsx"
	"mygameplatform/internal/httpx"
)

const minPasswordLength = 8

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) handleSignup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", err)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "Invalid request body")
		return
	}

	email, err := normalizeAndValidateEmail(req.Email)
	if err != nil {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", err)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "Invalid email")
		return
	}
	if len(req.Password) < minPasswordLength {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", errors.New("password too short"))
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "Invalid password")
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logRequest(r, http.StatusInternalServerError, "internal", err)
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	u, err := h.users.CreateUser(r.Context(), email, passwordHash)
	if err != nil {
		if errors.Is(err, db.ErrEmailConflict) {
			h.logRequest(r, http.StatusConflict, "conflict", err)
			errorsx.Write(w, http.StatusConflict, "conflict", "Email already registered")
			return
		}
		h.logRequest(r, http.StatusInternalServerError, "internal", err)
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	token, err := h.signer.Sign(u.ID, time.Now().UTC())
	if err != nil {
		h.logRequest(r, http.StatusInternalServerError, "internal", err)
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	h.logRequest(r, http.StatusOK, "", nil)
	httpx.WriteJSON(w, http.StatusOK, signupResponse{AccessToken: token})
}

func normalizeAndValidateEmail(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("email required")
	}
	if strings.Contains(trimmed, " ") || strings.Contains(trimmed, "\n") || strings.Contains(trimmed, "\t") {
		return "", errors.New("email contains whitespace")
	}
	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", err
	}
	if addr.Address != trimmed {
		return "", errors.New("email must be addr-spec only")
	}
	return strings.ToLower(trimmed), nil
}

func (h *Handler) logRequest(r *http.Request, status int, errCode string, err error) {
	attrs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
	}
	if errCode != "" {
		attrs = append(attrs, "error_code", errCode)
	}
	if err != nil {
		attrs = append(attrs, "err", err.Error())
	}
	h.logger.Info("request", attrs...)
}

