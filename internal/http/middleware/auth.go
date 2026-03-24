package middleware

import (
	"context"
	"net/http"
	"strings"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/http/response"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"
const bearerPrefix = "Bearer "

// Authenticator resolves a user from a bearer token.
type Authenticator interface {
	GetUserByToken(ctx context.Context, token string) (domain.User, error)
}

// RequireAuth protects routes that must have an authenticated user.
func RequireAuth(authenticator Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authenticator == nil {
				response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "auth service is not ready")
				return
			}

			authorizationHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
				response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authorization bearer token is required")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, bearerPrefix))
			if token == "" {
				response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authorization bearer token is required")
				return
			}

			user, err := authenticator.GetUserByToken(r.Context(), token)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
		})
	}
}

// CurrentUserFromContext returns the authenticated user when auth middleware has populated it.
func CurrentUserFromContext(ctx context.Context) (domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(domain.User)
	return user, ok
}
