package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ai-go-service/internal/domain"
)

type tokenPayload struct {
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
	Expiry int64  `json:"exp"`
}

// TokenManager signs and verifies bearer tokens for authenticated users.
type TokenManager struct {
	secret []byte
	nowFn  func() time.Time
}

// NewTokenManager creates a token manager using the provided secret.
func NewTokenManager(secret string, nowFn func() time.Time) (*TokenManager, error) {
	trimmedSecret := strings.TrimSpace(secret)
	if trimmedSecret == "" {
		return nil, fmt.Errorf("auth token secret must not be empty")
	}

	if nowFn == nil {
		nowFn = time.Now
	}

	return &TokenManager{secret: []byte(trimmedSecret), nowFn: nowFn}, nil
}

// Sign returns a signed bearer token for the given user identity.
func (manager *TokenManager) Sign(_ context.Context, userID int64, email string, ttl time.Duration) (string, error) {
	if manager == nil {
		return "", fmt.Errorf("token manager is unavailable")
	}

	if userID <= 0 {
		return "", fmt.Errorf("user id must be positive")
	}

	if ttl <= 0 {
		return "", fmt.Errorf("token ttl must be positive")
	}

	payloadBytes, err := json.Marshal(tokenPayload{
		UserID: userID,
		Email:  email,
		Expiry: manager.nowFn().Add(ttl).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("marshal token payload: %w", err)
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := manager.sign(encodedPayload)

	return encodedPayload + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

// Verify validates a bearer token and returns its claims.
func (manager *TokenManager) Verify(_ context.Context, token string) (domain.AuthClaims, error) {
	if manager == nil {
		return domain.AuthClaims{}, fmt.Errorf("token manager is unavailable")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return domain.AuthClaims{}, fmt.Errorf("token format is invalid")
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return domain.AuthClaims{}, fmt.Errorf("decode token signature: %w", err)
	}

	expectedSignature := manager.sign(parts[0])
	if subtle.ConstantTimeCompare(signature, expectedSignature) != 1 {
		return domain.AuthClaims{}, fmt.Errorf("token signature is invalid")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return domain.AuthClaims{}, fmt.Errorf("decode token payload: %w", err)
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return domain.AuthClaims{}, fmt.Errorf("unmarshal token payload: %w", err)
	}

	expiry := time.Unix(payload.Expiry, 0).UTC()
	if !expiry.After(manager.nowFn().UTC()) {
		return domain.AuthClaims{}, fmt.Errorf("token has expired")
	}

	return domain.AuthClaims{UserID: payload.UserID, Email: payload.Email, Expiry: expiry}, nil
}

func (manager *TokenManager) sign(value string) []byte {
	hasher := hmac.New(sha256.New, manager.secret)
	_, _ = hasher.Write([]byte(value))
	return hasher.Sum(nil)
}
