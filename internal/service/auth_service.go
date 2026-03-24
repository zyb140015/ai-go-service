package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ai-go-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUnauthorized indicates authentication failed.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrConflict indicates a resource already exists.
	ErrConflict = errors.New("conflict")
)

const bcryptCost = 12

// UserRepository defines the persistence behavior required by AuthService.
type UserRepository interface {
	Create(ctx context.Context, email string, displayName string, passwordHash string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
	UpdatePassword(ctx context.Context, id int64, passwordHash string) (domain.User, error)
}

// TokenIssuer signs and verifies user auth tokens.
type TokenIssuer interface {
	Sign(ctx context.Context, userID int64, email string, ttl time.Duration) (string, error)
	Verify(ctx context.Context, token string) (domain.AuthClaims, error)
}

// AuthResult groups the authenticated user and its bearer token.
type AuthResult struct {
	User         domain.User
	AccessToken  string
	RefreshToken string
}

// AuthService provides registration, login, and token-backed user lookup.
type AuthService struct {
	users    UserRepository
	tokens   TokenIssuer
	tokenTTL time.Duration
}

// NewAuthService creates an auth service from its dependencies.
func NewAuthService(users UserRepository, tokens TokenIssuer, tokenTTL time.Duration) *AuthService {
	return &AuthService{users: users, tokens: tokens, tokenTTL: tokenTTL}
}

// Register creates a user account and returns a signed token.
func (service *AuthService) Register(ctx context.Context, email string, displayName string, password string) (AuthResult, error) {
	if service == nil || service.users == nil || service.tokens == nil {
		return AuthResult{}, ErrUnavailable
	}

	normalizedEmail, normalizedName, normalizedPassword, err := normalizeCredentials(email, displayName, password)
	if err != nil {
		return AuthResult{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(normalizedPassword), bcryptCost)
	if err != nil {
		return AuthResult{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := service.users.Create(ctx, normalizedEmail, normalizedName, string(passwordHash))
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return AuthResult{}, ErrConflict
		}

		return AuthResult{}, fmt.Errorf("create user: %w", err)
	}

	accessToken, refreshToken, err := service.issueTokenPair(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Login verifies credentials and returns a signed token.
func (service *AuthService) Login(ctx context.Context, email string, password string) (AuthResult, error) {
	if service == nil || service.users == nil || service.tokens == nil {
		return AuthResult{}, ErrUnavailable
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	trimmedPassword := strings.TrimSpace(password)
	if normalizedEmail == "" || trimmedPassword == "" {
		return AuthResult{}, fmt.Errorf("email and password are required: %w", ErrInvalidInput)
	}

	user, err := service.users.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, ErrUnauthorized
		}

		return AuthResult{}, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(trimmedPassword)); err != nil {
		return AuthResult{}, ErrUnauthorized
	}

	accessToken, refreshToken, err := service.issueTokenPair(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Refresh exchanges a refresh token for a fresh token pair.
func (service *AuthService) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	if service == nil || service.users == nil || service.tokens == nil {
		return AuthResult{}, ErrUnavailable
	}

	claims, err := service.tokens.Verify(ctx, strings.TrimSpace(refreshToken))
	if err != nil {
		return AuthResult{}, ErrUnauthorized
	}

	user, err := service.users.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, ErrUnauthorized
		}

		return AuthResult{}, fmt.Errorf("get user for refresh: %w", err)
	}

	accessToken, nextRefreshToken, err := service.issueTokenPair(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, AccessToken: accessToken, RefreshToken: nextRefreshToken}, nil
}

// Logout verifies the bearer token so callers can invalidate local credentials safely.
func (service *AuthService) Logout(ctx context.Context, token string) error {
	if service == nil || service.tokens == nil {
		return ErrUnavailable
	}

	if _, err := service.tokens.Verify(ctx, strings.TrimSpace(token)); err != nil {
		return ErrUnauthorized
	}

	return nil
}

// ChangePassword verifies the current password, updates the stored hash, and issues fresh tokens.
func (service *AuthService) ChangePassword(ctx context.Context, token string, currentPassword string, newPassword string) (AuthResult, error) {
	if service == nil || service.users == nil || service.tokens == nil {
		return AuthResult{}, ErrUnavailable
	}

	claims, err := service.tokens.Verify(ctx, strings.TrimSpace(token))
	if err != nil {
		return AuthResult{}, ErrUnauthorized
	}

	trimmedCurrentPassword := strings.TrimSpace(currentPassword)
	trimmedNewPassword := strings.TrimSpace(newPassword)
	if trimmedCurrentPassword == "" || trimmedNewPassword == "" {
		return AuthResult{}, fmt.Errorf("currentPassword and newPassword are required: %w", ErrInvalidInput)
	}

	if len(trimmedNewPassword) < 8 {
		return AuthResult{}, fmt.Errorf("newPassword must be at least 8 characters: %w", ErrInvalidInput)
	}

	user, err := service.users.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, ErrUnauthorized
		}

		return AuthResult{}, fmt.Errorf("get user for password change: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(trimmedCurrentPassword)); err != nil {
		return AuthResult{}, ErrUnauthorized
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(trimmedNewPassword), bcryptCost)
	if err != nil {
		return AuthResult{}, fmt.Errorf("hash new password: %w", err)
	}

	updatedUser, err := service.users.UpdatePassword(ctx, user.ID, string(passwordHash))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, ErrUnauthorized
		}

		return AuthResult{}, fmt.Errorf("update password: %w", err)
	}

	accessToken, refreshToken, err := service.issueTokenPair(ctx, updatedUser)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: updatedUser, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// GetUserByToken resolves the current user from a bearer token.
func (service *AuthService) GetUserByToken(ctx context.Context, token string) (domain.User, error) {
	if service == nil || service.users == nil || service.tokens == nil {
		return domain.User{}, ErrUnavailable
	}

	claims, err := service.tokens.Verify(ctx, strings.TrimSpace(token))
	if err != nil {
		return domain.User{}, ErrUnauthorized
	}

	user, err := service.users.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.User{}, ErrUnauthorized
		}

		return domain.User{}, fmt.Errorf("get user by token: %w", err)
	}

	return user, nil
}

func normalizeCredentials(email string, displayName string, password string) (string, string, string, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	normalizedName := strings.TrimSpace(displayName)
	normalizedPassword := strings.TrimSpace(password)

	if normalizedEmail == "" {
		return "", "", "", fmt.Errorf("email is required: %w", ErrInvalidInput)
	}

	if normalizedName == "" {
		return "", "", "", fmt.Errorf("displayName is required: %w", ErrInvalidInput)
	}

	if normalizedPassword == "" {
		return "", "", "", fmt.Errorf("password is required: %w", ErrInvalidInput)
	}

	if !strings.Contains(normalizedEmail, "@") {
		return "", "", "", fmt.Errorf("email must be valid: %w", ErrInvalidInput)
	}

	if len(normalizedPassword) < 8 {
		return "", "", "", fmt.Errorf("password must be at least 8 characters: %w", ErrInvalidInput)
	}

	return normalizedEmail, normalizedName, normalizedPassword, nil
}

func (service *AuthService) issueTokenPair(ctx context.Context, user domain.User) (string, string, error) {
	accessToken, err := service.tokens.Sign(ctx, user.ID, user.Email, service.tokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err := service.tokens.Sign(ctx, user.ID, user.Email, service.tokenTTL*7)
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
