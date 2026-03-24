package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai-go-service/internal/domain"
	httpserver "ai-go-service/internal/http"
	"ai-go-service/internal/service"
)

type authServiceStub struct {
	registerFn func(ctx context.Context, email string, displayName string, password string) (service.AuthResult, error)
	loginFn    func(ctx context.Context, email string, password string) (service.AuthResult, error)
	currentFn  func(ctx context.Context, token string) (domain.User, error)
}

func (stub authServiceStub) Register(ctx context.Context, email string, displayName string, password string) (service.AuthResult, error) {
	return stub.registerFn(ctx, email, displayName, password)
}

func (stub authServiceStub) Login(ctx context.Context, email string, password string) (service.AuthResult, error) {
	return stub.loginFn(ctx, email, password)
}

func (stub authServiceStub) GetUserByToken(ctx context.Context, token string) (domain.User, error) {
	return stub.currentFn(ctx, token)
}

func TestRegisterHandlerReturnsToken(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 24, 14, 0, 0, 0, time.UTC)
	router := httpserver.NewRouter(newTestLogger(), nil, nil, authServiceStub{
		registerFn: func(_ context.Context, email string, displayName string, password string) (service.AuthResult, error) {
			return service.AuthResult{User: domain.User{ID: 1, Email: email, DisplayName: displayName, CreatedAt: now, UpdatedAt: now}, Token: "token-123"}, nil
		},
		loginFn: func(_ context.Context, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		currentFn: func(_ context.Context, _ string) (domain.User, error) { return domain.User{}, nil },
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"user@example.com","displayName":"Demo","password":"password123"}`))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["token"] != "token-123" {
		t.Fatalf("expected token, got %#v", body["token"])
	}
}

func TestLoginHandlerRejectsInvalidCredentials(t *testing.T) {
	t.Parallel()

	router := httpserver.NewRouter(newTestLogger(), nil, nil, authServiceStub{
		registerFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		loginFn: func(_ context.Context, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, service.ErrUnauthorized
		},
		currentFn: func(_ context.Context, _ string) (domain.User, error) { return domain.User{}, nil },
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"wrong"}`))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestCurrentUserHandlerReturnsUser(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 24, 14, 0, 0, 0, time.UTC)
	router := httpserver.NewRouter(newTestLogger(), nil, nil, authServiceStub{
		registerFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		loginFn: func(_ context.Context, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		currentFn: func(_ context.Context, token string) (domain.User, error) {
			if token != "valid-token" {
				t.Fatalf("unexpected token %q", token)
			}
			return domain.User{ID: 1, Email: "user@example.com", DisplayName: "Demo", CreatedAt: now, UpdatedAt: now}, nil
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}
