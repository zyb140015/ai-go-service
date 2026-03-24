package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/service"
)

type userRepositoryStub struct {
	createFn func(ctx context.Context, email string, displayName string, passwordHash string) (domain.User, error)
	byEmail  func(ctx context.Context, email string) (domain.User, error)
	byID     func(ctx context.Context, id int64) (domain.User, error)
}

func (stub userRepositoryStub) Create(ctx context.Context, email string, displayName string, passwordHash string) (domain.User, error) {
	return stub.createFn(ctx, email, displayName, passwordHash)
}

func (stub userRepositoryStub) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return stub.byEmail(ctx, email)
}

func (stub userRepositoryStub) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return stub.byID(ctx, id)
}

type tokenIssuerStub struct {
	signFn   func(ctx context.Context, userID int64, email string, ttl time.Duration) (string, error)
	verifyFn func(ctx context.Context, token string) (domain.AuthClaims, error)
}

func (stub tokenIssuerStub) Sign(ctx context.Context, userID int64, email string, ttl time.Duration) (string, error) {
	return stub.signFn(ctx, userID, email, ttl)
}

func (stub tokenIssuerStub) Verify(ctx context.Context, token string) (domain.AuthClaims, error) {
	return stub.verifyFn(ctx, token)
}

func TestAuthServiceRegisterNormalizesInput(t *testing.T) {
	t.Parallel()

	authService := service.NewAuthService(userRepositoryStub{
		createFn: func(_ context.Context, email string, displayName string, passwordHash string) (domain.User, error) {
			if email != "user@example.com" {
				t.Fatalf("expected normalized email, got %q", email)
			}
			if displayName != "Demo User" {
				t.Fatalf("expected normalized display name, got %q", displayName)
			}
			if passwordHash == "password123" || passwordHash == "" {
				t.Fatalf("expected hashed password, got %q", passwordHash)
			}
			return domain.User{ID: 1, Email: email, DisplayName: displayName}, nil
		},
		byEmail: func(_ context.Context, _ string) (domain.User, error) { return domain.User{}, nil },
		byID:    func(_ context.Context, _ int64) (domain.User, error) { return domain.User{}, nil },
	}, tokenIssuerStub{
		signFn: func(_ context.Context, userID int64, email string, ttl time.Duration) (string, error) {
			if userID != 1 || email != "user@example.com" || ttl <= 0 {
				t.Fatalf("unexpected token sign args: %d %s %s", userID, email, ttl)
			}
			return "signed-token", nil
		},
		verifyFn: func(_ context.Context, _ string) (domain.AuthClaims, error) {
			return domain.AuthClaims{}, nil
		},
	}, time.Hour)

	result, err := authService.Register(context.Background(), " User@Example.com ", " Demo User ", " password123 ")
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if result.Token != "signed-token" {
		t.Fatalf("expected token, got %q", result.Token)
	}
}

func TestAuthServiceLoginRejectsBadPassword(t *testing.T) {
	t.Parallel()

	authService := service.NewAuthService(userRepositoryStub{
		createFn: func(_ context.Context, _ string, _ string, _ string) (domain.User, error) { return domain.User{}, nil },
		byEmail: func(_ context.Context, _ string) (domain.User, error) {
			return domain.User{ID: 1, Email: "user@example.com", PasswordHash: "$2a$12$7A2yOmKB5tQ3fQxQv8W/L.0z5Y4xXHzkwmo7aX6ixkmKuuNHYsYAG"}, nil
		},
		byID: func(_ context.Context, _ int64) (domain.User, error) { return domain.User{}, nil },
	}, tokenIssuerStub{
		signFn: func(_ context.Context, _ int64, _ string, _ time.Duration) (string, error) { return "", nil },
		verifyFn: func(_ context.Context, _ string) (domain.AuthClaims, error) {
			return domain.AuthClaims{}, nil
		},
	}, time.Hour)

	_, err := authService.Login(context.Background(), "user@example.com", "wrong-password")
	if !errors.Is(err, service.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}

func TestAuthServiceGetUserByTokenResolvesUser(t *testing.T) {
	t.Parallel()

	authService := service.NewAuthService(userRepositoryStub{
		createFn: func(_ context.Context, _ string, _ string, _ string) (domain.User, error) { return domain.User{}, nil },
		byEmail:  func(_ context.Context, _ string) (domain.User, error) { return domain.User{}, nil },
		byID: func(_ context.Context, id int64) (domain.User, error) {
			return domain.User{ID: id, Email: "user@example.com", DisplayName: "Demo User"}, nil
		},
	}, tokenIssuerStub{
		signFn: func(_ context.Context, _ int64, _ string, _ time.Duration) (string, error) { return "", nil },
		verifyFn: func(_ context.Context, token string) (domain.AuthClaims, error) {
			if token != "valid-token" {
				t.Fatalf("unexpected token %q", token)
			}
			return domain.AuthClaims{UserID: 42, Email: "user@example.com"}, nil
		},
	}, time.Hour)

	user, err := authService.GetUserByToken(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("get user by token: %v", err)
	}

	if user.ID != 42 {
		t.Fatalf("expected user id 42, got %d", user.ID)
	}
}
