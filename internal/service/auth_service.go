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
}

// TokenIssuer signs and verifies user auth tokens.
type TokenIssuer interface {
	Sign(ctx context.Context, userID int64, email string, ttl time.Duration) (string, error)
	Verify(ctx context.Context, token string) (domain.AuthClaims, error)
}

// AuthResult groups the authenticated user and its bearer token.
type AuthResult struct {
	User  domain.User
	Token string
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

	token, err := service.tokens.Sign(ctx, user.ID, user.Email, service.tokenTTL)
	if err != nil {
		return AuthResult{}, fmt.Errorf("sign token: %w", err)
	}

	return AuthResult{User: user, Token: token}, nil
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

	token, err := service.tokens.Sign(ctx, user.ID, user.Email, service.tokenTTL)
	if err != nil {
		return AuthResult{}, fmt.Errorf("sign token: %w", err)
	}

	return AuthResult{User: user, Token: token}, nil
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
