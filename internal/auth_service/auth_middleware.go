package auth_service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"mygameplatform/internal/errorsx"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.verifier == nil {
			h.logRequest(r, http.StatusInternalServerError, "internal", errors.New("token verifier not configured"))
			errorsx.Write(w, http.StatusInternalServerError, "internal", "Internal error")
			return
		}
		token := bearerToken(r.Header.Get("Authorization"))
		if token == "" {
			h.logRequest(r, http.StatusUnauthorized, "forbidden", nil)
			errorsx.Write(w, http.StatusUnauthorized, "forbidden", "Unauthorized")
			return
		}
		userID, err := h.verifier.Verify(token, time.Now().UTC())
		if err != nil {
			h.logRequest(r, http.StatusUnauthorized, "forbidden", nil)
			errorsx.Write(w, http.StatusUnauthorized, "forbidden", "Unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}

func bearerToken(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
