package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"ai-go-service/internal/http/response"
)

// Recoverer converts panics into a stable JSON error while logging the stack trace.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error(
						"panic recovered",
						"panic", fmt.Sprint(recovered),
						"method", r.Method,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)

					response.Error(w, http.StatusInternalServerError, response.CodeInternal, "internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
