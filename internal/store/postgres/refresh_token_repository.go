package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/store/sqlcdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// RefreshTokenRepository stores refresh token state in PostgreSQL.
type RefreshTokenRepository struct {
	queries *sqlcdb.Queries
}

// NewRefreshTokenRepository returns a repository when a PostgreSQL pool is available.
func NewRefreshTokenRepository(pool *Pool) *RefreshTokenRepository {
	if pool == nil || pool.pool == nil {
		return nil
	}

	return &RefreshTokenRepository{queries: sqlcdb.New(pool.pool)}
}

// Create persists a refresh token hash for later rotation or revocation.
func (repository *RefreshTokenRepository) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	if repository == nil || repository.queries == nil {
		return fmt.Errorf("refresh token repository is unavailable")
	}

	_, err := repository.queries.CreateRefreshToken(ctx, sqlcdb.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}

	return nil
}

// GetActive verifies that a refresh token exists, is active, and has not expired.
func (repository *RefreshTokenRepository) GetActive(ctx context.Context, userID int64, tokenHash string, now time.Time) error {
	if repository == nil || repository.queries == nil {
		return fmt.Errorf("refresh token repository is unavailable")
	}

	_, err := repository.queries.GetActiveRefreshToken(ctx, sqlcdb.GetActiveRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}

		return fmt.Errorf("get active refresh token: %w", err)
	}

	return nil
}

// Revoke marks a single refresh token as revoked.
func (repository *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	if repository == nil || repository.queries == nil {
		return fmt.Errorf("refresh token repository is unavailable")
	}

	_, err := repository.queries.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}

// RevokeByUser marks all refresh tokens for a user as revoked.
func (repository *RefreshTokenRepository) RevokeByUser(ctx context.Context, userID int64) error {
	if repository == nil || repository.queries == nil {
		return fmt.Errorf("refresh token repository is unavailable")
	}

	_, err := repository.queries.RevokeRefreshTokensByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("revoke refresh tokens by user: %w", err)
	}

	return nil
}
