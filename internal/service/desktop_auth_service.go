package service

import (
	"context"
	"fmt"
	"strings"

	"ai-go-service/internal/integration/goadmin"
)

// GoAdminAuthClient defines the upstream auth behavior required by DesktopAuthService.
type GoAdminAuthClient interface {
	Login(ctx context.Context, username string, password string) (goadmin.LoginEnvelope, error)
	Refresh(ctx context.Context, refreshToken string) (goadmin.LoginResult, error)
	GetCurrentUser(ctx context.Context, accessToken string) (goadmin.ProfileResult, error)
	GetCurrentMenus(ctx context.Context, accessToken string) ([]goadmin.MenuResult, error)
	UpdateCurrentPassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error
	UpdateCurrentProfile(ctx context.Context, accessToken string, username string, avatar string, sex string) (goadmin.ProfileResult, error)
	UploadCurrentAvatar(ctx context.Context, accessToken string, fileName string, fileContent []byte) (string, error)
}

// DesktopAuthService adapts go-admin auth responses to the desktop-facing API.
type DesktopAuthService struct {
	client GoAdminAuthClient
}

// DesktopUser is the stable desktop-facing user DTO.
type DesktopUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone,omitempty"`
	Remark    string `json:"remark,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// DesktopTokenPair is the stable desktop-facing token DTO.
type DesktopTokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// DesktopLoginResponse is returned when the desktop signs in.
type DesktopLoginResponse struct {
	User   DesktopUser      `json:"user"`
	Tokens DesktopTokenPair `json:"tokens"`
}

// DesktopRefreshResponse is returned when the desktop refreshes tokens.
type DesktopRefreshResponse struct {
	Tokens DesktopTokenPair `json:"tokens"`
}

// DesktopPasswordUpdateRequest is the stable desktop-facing password update DTO.
type DesktopPasswordUpdateRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// DesktopProfileUpdateRequest is the stable desktop-facing profile update DTO.
type DesktopProfileUpdateRequest struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
}

// NewDesktopAuthService creates a desktop auth adapter service.
func NewDesktopAuthService(client GoAdminAuthClient) *DesktopAuthService {
	return &DesktopAuthService{client: client}
}

// Login authenticates through go-admin and reshapes the result.
func (service *DesktopAuthService) Login(ctx context.Context, username string, password string) (DesktopLoginResponse, error) {
	if service == nil || service.client == nil {
		return DesktopLoginResponse{}, ErrUnavailable
	}

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return DesktopLoginResponse{}, fmt.Errorf("username and password are required: %w", ErrInvalidInput)
	}

	payload, err := service.client.Login(ctx, username, password)
	if err != nil {
		return DesktopLoginResponse{}, err
	}

	return DesktopLoginResponse{
		User: mapLoginUser(payload.UserInfo),
		Tokens: DesktopTokenPair{
			AccessToken:  payload.Result.Token,
			RefreshToken: payload.Result.RefreshToken,
		},
	}, nil
}

// Refresh exchanges the refresh token through go-admin.
func (service *DesktopAuthService) Refresh(ctx context.Context, refreshToken string) (DesktopRefreshResponse, error) {
	if service == nil || service.client == nil {
		return DesktopRefreshResponse{}, ErrUnavailable
	}

	if strings.TrimSpace(refreshToken) == "" {
		return DesktopRefreshResponse{}, fmt.Errorf("refreshToken is required: %w", ErrInvalidInput)
	}

	payload, err := service.client.Refresh(ctx, refreshToken)
	if err != nil {
		return DesktopRefreshResponse{}, err
	}

	return DesktopRefreshResponse{
		Tokens: DesktopTokenPair{
			AccessToken:  payload.Token,
			RefreshToken: payload.RefreshToken,
		},
	}, nil
}

// CurrentUser resolves the current desktop user through go-admin.
func (service *DesktopAuthService) CurrentUser(ctx context.Context, accessToken string) (DesktopUser, error) {
	if service == nil || service.client == nil {
		return DesktopUser{}, ErrUnavailable
	}

	if strings.TrimSpace(accessToken) == "" {
		return DesktopUser{}, fmt.Errorf("access token is required: %w", ErrInvalidInput)
	}

	payload, err := service.client.GetCurrentUser(ctx, accessToken)
	if err != nil {
		return DesktopUser{}, err
	}

	return mapProfileUser(payload), nil
}

// UpdateCurrentPassword updates the current user's password through go-admin.
func (service *DesktopAuthService) UpdateCurrentPassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error {
	if service == nil || service.client == nil {
		return ErrUnavailable
	}

	if strings.TrimSpace(accessToken) == "" {
		return fmt.Errorf("access token is required: %w", ErrInvalidInput)
	}

	if strings.TrimSpace(oldPassword) == "" || strings.TrimSpace(newPassword) == "" {
		return fmt.Errorf("oldPassword and newPassword are required: %w", ErrInvalidInput)
	}

	return service.client.UpdateCurrentPassword(ctx, accessToken, oldPassword, newPassword)
}

// UpdateCurrentProfile updates the current user's profile through go-admin.
func (service *DesktopAuthService) UpdateCurrentProfile(ctx context.Context, accessToken string, input DesktopProfileUpdateRequest) (DesktopUser, error) {
	if service == nil || service.client == nil {
		return DesktopUser{}, ErrUnavailable
	}
	if strings.TrimSpace(accessToken) == "" {
		return DesktopUser{}, fmt.Errorf("access token is required: %w", ErrInvalidInput)
	}
	if strings.TrimSpace(input.Username) == "" {
		return DesktopUser{}, fmt.Errorf("username is required: %w", ErrInvalidInput)
	}
	result, err := service.client.UpdateCurrentProfile(ctx, accessToken, input.Username, input.Avatar, input.Sex)
	if err != nil {
		return DesktopUser{}, err
	}
	return mapProfileUser(result), nil
	}

// UploadCurrentAvatar uploads the current user's avatar through go-admin.
func (service *DesktopAuthService) UploadCurrentAvatar(ctx context.Context, accessToken string, fileName string, fileContent []byte) (string, error) {
	if service == nil || service.client == nil {
		return "", ErrUnavailable
	}
	if strings.TrimSpace(accessToken) == "" {
		return "", fmt.Errorf("access token is required: %w", ErrInvalidInput)
	}
	if strings.TrimSpace(fileName) == "" || len(fileContent) == 0 {
		return "", fmt.Errorf("avatar file is required: %w", ErrInvalidInput)
	}
	return service.client.UploadCurrentAvatar(ctx, accessToken, fileName, fileContent)
}

func mapLoginUser(input goadmin.LoginUserInfo) DesktopUser {
	return DesktopUser{
		ID:        input.ID,
		Username:  input.Username,
		Email:     input.Email,
		Avatar:    input.Avatar,
		Sex:       input.Sex,
		Phone:     input.Phone,
		Remark:    input.Remark,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}

func mapProfileUser(input goadmin.ProfileResult) DesktopUser {
	return DesktopUser{
		ID:        input.ID,
		Username:  input.Username,
		Email:     input.Email,
		Avatar:    input.Avatar,
		Sex:       input.Sex,
		CreatedAt: input.CreatedAt,
	}
}
