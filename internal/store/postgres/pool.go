package postgres

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool wraps a pgx connection pool to keep the infrastructure boundary explicit.
type Pool struct {
	pool *pgxpool.Pool
}

// NewPool creates a PostgreSQL pool when a database URL is configured.
// It returns nil when the service is running without a database dependency.
func NewPool(ctx context.Context, databaseURL string) (*Pool, error) {
	if databaseURL == "" {
		return nil, nil
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("open postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Pool{pool: pool}, nil
}

// Close releases the underlying connection pool.
func (pool *Pool) Close() {
	if pool == nil || pool.pool == nil {
		return
	}

	pool.pool.Close()
}

// Check verifies readiness for the current request.
func (pool *Pool) Check(r *http.Request) error {
	if pool == nil || pool.pool == nil {
		return nil
	}

	return pool.pool.Ping(r.Context())
}
