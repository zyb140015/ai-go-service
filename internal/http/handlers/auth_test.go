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
	refreshFn  func(ctx context.Context, refreshToken string) (service.AuthResult, error)
	logoutFn   func(ctx context.Context, token string) error
	changeFn   func(ctx context.Context, token string, currentPassword string, newPassword string) (service.AuthResult, error)
	currentFn  func(ctx context.Context, token string) (domain.User, error)
}

func (stub authServiceStub) Register(ctx context.Context, email string, displayName string, password string) (service.AuthResult, error) {
	return stub.registerFn(ctx, email, displayName, password)
}

func (stub authServiceStub) Login(ctx context.Context, email string, password string) (service.AuthResult, error) {
	return stub.loginFn(ctx, email, password)
}

func (stub authServiceStub) Refresh(ctx context.Context, refreshToken string) (service.AuthResult, error) {
	return stub.refreshFn(ctx, refreshToken)
}

func (stub authServiceStub) Logout(ctx context.Context, token string) error {
	return stub.logoutFn(ctx, token)
}

func (stub authServiceStub) ChangePassword(ctx context.Context, token string, currentPassword string, newPassword string) (service.AuthResult, error) {
	return stub.changeFn(ctx, token, currentPassword, newPassword)
}

func (stub authServiceStub) GetUserByToken(ctx context.Context, token string) (domain.User, error) {
	return stub.currentFn(ctx, token)
}

func TestRegisterHandlerReturnsTokenPair(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 24, 14, 0, 0, 0, time.UTC)
	router := httpserver.NewRouter(newTestLogger(), nil, nil, authServiceStub{
		registerFn: func(_ context.Context, email string, displayName string, password string) (service.AuthResult, error) {
			return service.AuthResult{User: domain.User{ID: 1, Email: email, DisplayName: displayName, CreatedAt: now, UpdatedAt: now}, AccessToken: "access-123", RefreshToken: "refresh-123"}, nil
		},
		loginFn: func(_ context.Context, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		refreshFn: func(_ context.Context, _ string) (service.AuthResult, error) { return service.AuthResult{}, nil },
		logoutFn:  func(_ context.Context, _ string) error { return nil },
		changeFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
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

	if body["accessToken"] != "access-123" || body["refreshToken"] != "refresh-123" {
		t.Fatalf("expected token pair, got %#v", body)
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
		refreshFn: func(_ context.Context, _ string) (service.AuthResult, error) { return service.AuthResult{}, nil },
		logoutFn:  func(_ context.Context, _ string) error { return nil },
		changeFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
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

func TestRefreshHandlerReturnsNewTokens(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 24, 15, 0, 0, 0, time.UTC)
	router := httpserver.NewRouter(newTestLogger(), nil, nil, authServiceStub{
		registerFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		loginFn: func(_ context.Context, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		refreshFn: func(_ context.Context, refreshToken string) (service.AuthResult, error) {
			if refreshToken != "refresh-123" {
				t.Fatalf("unexpected refresh token %q", refreshToken)
			}
			return service.AuthResult{User: domain.User{ID: 1, Email: "user@example.com", DisplayName: "Demo", CreatedAt: now, UpdatedAt: now}, AccessToken: "access-456", RefreshToken: "refresh-456"}, nil
		},
		logoutFn: func(_ context.Context, _ string) error { return nil },
		changeFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		currentFn: func(_ context.Context, _ string) (domain.User, error) { return domain.User{}, nil },
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(`{"refreshToken":"refresh-123"}`))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestLogoutHandlerRequiresBearerToken(t *testing.T) {
	t.Parallel()

	router := httpserver.NewRouter(newTestLogger(), nil, nil, authServiceStub{
		registerFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		loginFn: func(_ context.Context, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		refreshFn: func(_ context.Context, _ string) (service.AuthResult, error) { return service.AuthResult{}, nil },
		logoutFn:  func(_ context.Context, _ string) error { return nil },
		changeFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		currentFn: func(_ context.Context, _ string) (domain.User, error) { return domain.User{}, nil },
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestChangePasswordReturnsNewTokens(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 24, 16, 0, 0, 0, time.UTC)
	router := httpserver.NewRouter(newTestLogger(), nil, nil, authServiceStub{
		registerFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		loginFn: func(_ context.Context, _ string, _ string) (service.AuthResult, error) {
			return service.AuthResult{}, nil
		},
		refreshFn: func(_ context.Context, _ string) (service.AuthResult, error) { return service.AuthResult{}, nil },
		logoutFn:  func(_ context.Context, _ string) error { return nil },
		changeFn: func(_ context.Context, token string, currentPassword string, newPassword string) (service.AuthResult, error) {
			if token != "valid-token" || currentPassword != "oldpass123" || newPassword != "newpass123" {
				t.Fatalf("unexpected change password arguments")
			}
			return service.AuthResult{User: domain.User{ID: 1, Email: "user@example.com", DisplayName: "Demo", CreatedAt: now, UpdatedAt: now}, AccessToken: "access-789", RefreshToken: "refresh-789"}, nil
		},
		currentFn: func(_ context.Context, _ string) (domain.User, error) { return domain.User{}, nil },
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/change-password", strings.NewReader(`{"currentPassword":"oldpass123","newPassword":"newpass123"}`))
	request.Header.Set("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
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
		refreshFn: func(_ context.Context, _ string) (service.AuthResult, error) { return service.AuthResult{}, nil },
		logoutFn:  func(_ context.Context, _ string) error { return nil },
		changeFn: func(_ context.Context, _ string, _ string, _ string) (service.AuthResult, error) {
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
