package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"ai-go-service/internal/http/response"
	"ai-go-service/internal/integration/goadmin"
	"ai-go-service/internal/service"
)

const authBearerPrefix = "Bearer "

// DesktopAuthUseCase defines the auth use cases exposed by the desktop BFF.
type DesktopAuthUseCase interface {
	Login(ctx context.Context, username string, password string) (service.DesktopLoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (service.DesktopRefreshResponse, error)
	CurrentUser(ctx context.Context, accessToken string) (service.DesktopUser, error)
	UpdateCurrentPassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error
	UpdateCurrentProfile(ctx context.Context, accessToken string, input service.DesktopProfileUpdateRequest) (service.DesktopUser, error)
	UploadCurrentAvatar(ctx context.Context, accessToken string, fileName string, fileContent []byte) (string, error)
}

// DesktopLoginRequest is the expected JSON payload for desktop login.
type DesktopLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DesktopRefreshRequest is the expected JSON payload for desktop token refresh.
type DesktopRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// DesktopPasswordUpdateRequest is the expected JSON payload for desktop password update.
type DesktopPasswordUpdateRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// DesktopProfileUpdateRequest is the expected JSON payload for desktop profile update.
type DesktopProfileUpdateRequest struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
}

// DesktopLoginHandler returns an HTTP handler that logs the desktop user in via go-admin.
func DesktopLoginHandler(authService DesktopAuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop auth service is not ready")
			return
		}

		var request DesktopLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}

		result, err := authService.Login(r.Context(), request.Username, request.Password)
		if err != nil {
			handleDesktopAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopRefreshHandler returns an HTTP handler that refreshes desktop tokens via go-admin.
func DesktopRefreshHandler(authService DesktopAuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop auth service is not ready")
			return
		}

		var request DesktopRefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}

		result, err := authService.Refresh(r.Context(), request.RefreshToken)
		if err != nil {
			handleDesktopAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopCurrentUserHandler returns an HTTP handler that resolves the desktop current user via go-admin.
func DesktopCurrentUserHandler(authService DesktopAuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop auth service is not ready")
			return
		}

		accessToken, err := parseDesktopBearerToken(r)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}

		result, err := authService.CurrentUser(r.Context(), accessToken)
		if err != nil {
			handleDesktopAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopUpdatePasswordHandler returns an HTTP handler that updates the desktop current user's password via go-admin.
func DesktopUpdatePasswordHandler(authService DesktopAuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop auth service is not ready")
			return
		}

		accessToken, err := parseDesktopBearerToken(r)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}

		var request DesktopPasswordUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}

		if err := authService.UpdateCurrentPassword(r.Context(), accessToken, request.OldPassword, request.NewPassword); err != nil {
			handleDesktopAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]string{"message": "密码修改成功"})
	}
}

// DesktopUpdateProfileHandler returns an HTTP handler that updates the desktop current user's profile via go-admin.
func DesktopUpdateProfileHandler(authService DesktopAuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop auth service is not ready")
			return
		}
		accessToken, err := parseDesktopBearerToken(r)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request DesktopProfileUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		result, err := authService.UpdateCurrentProfile(r.Context(), accessToken, service.DesktopProfileUpdateRequest{
			Username: request.Username,
			Avatar:   request.Avatar,
			Sex:      request.Sex,
		})
		if err != nil {
			handleDesktopAuthError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopUploadAvatarHandler returns an HTTP handler that uploads the desktop current user's avatar via go-admin.
func DesktopUploadAvatarHandler(authService DesktopAuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop auth service is not ready")
			return
		}
		accessToken, err := parseDesktopBearerToken(r)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "avatar file is required")
			return
		}
		defer file.Close()
		var buffer bytes.Buffer
		if _, err := io.Copy(&buffer, file); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "failed to read avatar file")
			return
		}
		avatar, err := authService.UploadCurrentAvatar(r.Context(), accessToken, header.Filename, buffer.Bytes())
		if err != nil {
			handleDesktopAuthError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"avatar": avatar})
	}
}

func parseDesktopBearerToken(r *http.Request) (string, error) {
	authorizationHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(authorizationHeader, authBearerPrefix) {
		return "", errors.New("authorization bearer token is required")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, authBearerPrefix))
	if token == "" {
		return "", errors.New("authorization bearer token is required")
	}

	return token, nil
}

func handleDesktopAuthError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrInvalidInput) {
		response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, err.Error())
		return
	}

	if errors.Is(err, goadmin.ErrUnauthorized) {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
		return
	}

	if errors.Is(err, goadmin.ErrUpstreamFailure) || errors.Is(err, service.ErrUnavailable) {
		response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "go-admin service is unavailable")
		return
	}

	if errors.Is(err, goadmin.ErrInvalidResponse) {
		response.Error(w, http.StatusBadGateway, response.CodeInternal, "go-admin returned an invalid response")
		return
	}

	var businessError *goadmin.BusinessError
	if errors.As(err, &businessError) {
		response.Error(w, http.StatusBadGateway, response.CodeServiceUnavailable, businessError.Error())
		return
	}

	response.Error(w, http.StatusInternalServerError, response.CodeInternal, "internal server error")
}
