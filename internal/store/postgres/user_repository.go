package postgres

import (
	"context"
	"errors"
	"fmt"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/store/sqlcdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

// UserRepository provides user persistence backed by PostgreSQL.
type UserRepository struct {
	queries *sqlcdb.Queries
}

// NewUserRepository returns a repository when a PostgreSQL pool is available.
func NewUserRepository(pool *Pool) *UserRepository {
	if pool == nil || pool.pool == nil {
		return nil
	}

	return &UserRepository{queries: sqlcdb.New(pool.pool)}
}

// Create stores a user and returns the persisted record.
func (repository *UserRepository) Create(ctx context.Context, email string, displayName string, passwordHash string) (domain.User, error) {
	if repository == nil || repository.queries == nil {
		return domain.User{}, fmt.Errorf("user repository is unavailable")
	}

	record, err := repository.queries.CreateUser(ctx, sqlcdb.CreateUserParams{
		Email:        email,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
	})
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == uniqueViolationCode {
			return domain.User{}, domain.ErrConflict
		}

		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return mapUser(record), nil
}

// GetByEmail returns one user by email address.
func (repository *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if repository == nil || repository.queries == nil {
		return domain.User{}, fmt.Errorf("user repository is unavailable")
	}

	record, err := repository.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return mapUser(record), nil
}

// GetByID returns one user by identifier.
func (repository *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	if repository == nil || repository.queries == nil {
		return domain.User{}, fmt.Errorf("user repository is unavailable")
	}

	record, err := repository.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return mapUser(record), nil
}

func mapUser(record sqlcdb.AppUser) domain.User {
	return domain.User{
		ID:           record.ID,
		Email:        record.Email,
		DisplayName:  record.DisplayName,
		PasswordHash: record.PasswordHash,
		CreatedAt:    record.CreatedAt.Time,
		UpdatedAt:    record.UpdatedAt.Time,
	}
}
