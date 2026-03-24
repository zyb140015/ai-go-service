package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const unknownStatusCode = 0

// RequestLogger writes one structured log entry per HTTP request.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			wrappedWriter := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(wrappedWriter, r)

			statusCode := wrappedWriter.Status()
			if statusCode == unknownStatusCode {
				statusCode = http.StatusOK
			}

			logger.Info(
				"http request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", statusCode,
				"bytes", wrappedWriter.BytesWritten(),
				"duration", time.Since(startedAt),
				"request_id", chimiddleware.GetReqID(r.Context()),
			)
		})
	}
}
