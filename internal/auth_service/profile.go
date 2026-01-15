package auth_service

import (
	"errors"
	"net/http"
	"strings"

	"mygameplatform/internal/db"
	"mygameplatform/internal/errorsx"
	"mygameplatform/internal/httpx"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 32
)

// Username uniqueness is not enforced in MVP.

type profileResponse struct {
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
}

type profileUpdateRequest struct {
	Username  *string `json:"username"`
	AvatarURL *string `json:"avatar_url"`
}

func (h *Handler) handleMeGet(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		h.logRequest(r, http.StatusInternalServerError, "internal", errors.New("missing user id in context"))
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	u, err := h.users.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			h.logRequest(r, http.StatusUnauthorized, "forbidden", nil)
			errorsx.Write(w, http.StatusUnauthorized, "forbidden", "Unauthorized")
			return
		}
		h.logRequest(r, http.StatusInternalServerError, "internal", err)
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	h.logRequest(r, http.StatusOK, "", nil)
	httpx.WriteJSON(w, http.StatusOK, profileResponse{Username: u.Username, AvatarURL: u.AvatarURL})
}

func (h *Handler) handleMePatch(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		h.logRequest(r, http.StatusInternalServerError, "internal", errors.New("missing user id in context"))
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "" && !httpx.IsJSONContentType(r) {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "Content-Type must be application/json")
		return
	}

	var req profileUpdateRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "Invalid JSON body")
		return
	}

	if req.Username == nil {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "username: required")
		return
	}
	username, err := normalizeUsername(*req.Username)
	if err != nil {
		h.logRequest(r, http.StatusBadRequest, "validation_failed", nil)
		errorsx.Write(w, http.StatusBadRequest, "validation_failed", "username: invalid")
		return
	}

	avatarURL := normalizeAvatarURL(req.AvatarURL)

	u, err := h.users.UpdateUserProfile(r.Context(), userID, username, avatarURL)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			h.logRequest(r, http.StatusUnauthorized, "forbidden", nil)
			errorsx.Write(w, http.StatusUnauthorized, "forbidden", "Unauthorized")
			return
		}
		h.logRequest(r, http.StatusInternalServerError, "internal", err)
		errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
		return
	}

	h.logRequest(r, http.StatusOK, "", nil)
	httpx.WriteJSON(w, http.StatusOK, profileResponse{Username: u.Username, AvatarURL: u.AvatarURL})
}

func normalizeUsername(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("username required")
	}
	trimmed = strings.ToLower(trimmed)
	if len(trimmed) < minUsernameLength || len(trimmed) > maxUsernameLength {
		return "", errors.New("username length invalid")
	}
	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' || r == '+' {
			continue
		}
		return "", errors.New("username contains invalid characters")
	}
	return trimmed, nil
}

func normalizeAvatarURL(raw *string) *string {
	if raw == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func defaultUsernameFromEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	local := strings.TrimSpace(parts[0])
	if local == "" {
		local = "user"
	}
	local = strings.ToLower(local)
	var b strings.Builder
	for _, r := range local {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' || r == '+' {
			b.WriteRune(r)
			continue
		}
		b.WriteRune('_')
	}
	username := strings.Trim(b.String(), "._-+")
	if len(username) < minUsernameLength {
		username = "user_" + username
	}
	if len(username) > maxUsernameLength {
		username = username[:maxUsernameLength]
	}
	if username == "" {
		username = "user"
	}
	return username
}
