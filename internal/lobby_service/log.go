package lobby_service

import (
	"fmt"
	"net/http"
)

func (h *Handler) logRequest(r *http.Request, status int, errCode string, err error) {
	attrs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
	}
	if errCode != "" {
		attrs = append(attrs, "error_code", errCode)
	}
	if err != nil && status >= 500 {
		attrs = append(attrs, "err_kind", fmt.Sprintf("%T", err))
		h.logger.Error("request", attrs...)
		return
	}
	h.logger.Info("request", attrs...)
}
