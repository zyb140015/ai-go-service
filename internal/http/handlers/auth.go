package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/http/response"
	"ai-go-service/internal/service"
)

const bearerPrefix = "Bearer "

// AuthService defines auth use cases exposed to the HTTP layer.
type AuthService interface {
	Register(ctx context.Context, email string, displayName string, password string) (service.AuthResult, error)
	Login(ctx context.Context, email string, password string) (service.AuthResult, error)
	GetUserByToken(ctx context.Context, token string) (domain.User, error)
}

// RegisterRequest is the expected JSON payload for user registration.
type RegisterRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

// LoginRequest is the expected JSON payload for login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthUserResponse is the public JSON shape returned for authenticated users.
type AuthUserResponse struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// AuthResponse returns a user and its bearer token.
type AuthResponse struct {
	User  AuthUserResponse `json:"user"`
	Token string           `json:"token"`
}

// RegisterHandler returns an HTTP handler that creates an account.
func RegisterHandler(authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "auth service is not ready")
			return
		}

		var request RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}

		result, err := authService.Register(r.Context(), request.Email, request.DisplayName, request.Password)
		if err != nil {
			handleAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusCreated, AuthResponse{User: toAuthUserResponse(result.User), Token: result.Token})
	}
}

// LoginHandler returns an HTTP handler that verifies user credentials.
func LoginHandler(authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "auth service is not ready")
			return
		}

		var request LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}

		result, err := authService.Login(r.Context(), request.Email, request.Password)
		if err != nil {
			handleAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, AuthResponse{User: toAuthUserResponse(result.User), Token: result.Token})
	}
}

// CurrentUserHandler returns the authenticated user from a bearer token.
func CurrentUserHandler(authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "auth service is not ready")
			return
		}

		token, ok := parseBearerToken(w, r)
		if !ok {
			return
		}

		user, err := authService.GetUserByToken(r.Context(), token)
		if err != nil {
			handleAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, toAuthUserResponse(user))
	}
}

func handleAuthError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrInvalidInput) {
		response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrConflict) {
		response.Error(w, http.StatusConflict, response.CodeConflict, "resource already exists")
		return
	}

	if errors.Is(err, service.ErrUnauthorized) {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
		return
	}

	if errors.Is(err, service.ErrUnavailable) {
		response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "auth service is not ready")
		return
	}

	response.Error(w, http.StatusInternalServerError, response.CodeInternal, "internal server error")
}

func parseBearerToken(w http.ResponseWriter, r *http.Request) (string, bool) {
	authorizationHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authorization bearer token is required")
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, bearerPrefix))
	if token == "" {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authorization bearer token is required")
		return "", false
	}

	return token, true
}

func toAuthUserResponse(user domain.User) AuthUserResponse {
	return AuthUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
