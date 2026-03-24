package postgres_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"ai-go-service/internal/service"
	"ai-go-service/internal/store/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestNoteRepositoryCRUDAndPagination(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pool, cleanup := newTestPool(t, ctx)
	defer cleanup()

	repository := postgres.NewNoteRepository(pool)
	if repository == nil {
		t.Fatal("expected note repository")
	}

	createdOne, err := repository.Create(ctx, "hello", "world")
	if err != nil {
		t.Fatalf("create first note: %v", err)
	}

	createdTwo, err := repository.Create(ctx, "alpha", "beta")
	if err != nil {
		t.Fatalf("create second note: %v", err)
	}

	page, total, err := repository.List(ctx, service.NoteListOptions{
		Page:          1,
		PageSize:      1,
		SortField:     service.NoteSortFieldTitle,
		SortDirection: service.NoteSortDirectionAsc,
		Query:         "a",
	})
	if err != nil {
		t.Fatalf("list notes: %v", err)
	}

	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	if len(page) != 1 || page[0].ID != createdTwo.ID {
		t.Fatalf("expected filtered note %d, got %#v", createdTwo.ID, page)
	}

	loaded, err := repository.GetByID(ctx, createdOne.ID)
	if err != nil {
		t.Fatalf("get note by id: %v", err)
	}

	if loaded.Title != createdOne.Title {
		t.Fatalf("expected title %q, got %q", createdOne.Title, loaded.Title)
	}

	updated, err := repository.Update(ctx, createdOne.ID, "hello-updated", "world-updated")
	if err != nil {
		t.Fatalf("update note: %v", err)
	}

	if updated.Title != "hello-updated" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	if err := repository.Delete(ctx, createdTwo.ID); err != nil {
		t.Fatalf("delete note: %v", err)
	}

	_, err = repository.GetByID(ctx, createdTwo.ID)
	if err == nil {
		t.Fatal("expected deleted note lookup to fail")
	}
}

func newTestPool(t *testing.T, ctx context.Context) (*postgres.Pool, func()) {
	t.Helper()
	testcontainers.SkipIfProviderIsNotHealthy(t)

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("ai_go_service_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
	)
	if err != nil {
		t.Skipf("skip integration test because postgres container could not start: %v", err)
	}

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("build connection string: %v", err)
	}

	pgxPool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("open pgx pool: %v", err)
	}

	if err := applyMigrations(ctx, pgxPool); err != nil {
		pgxPool.Close()
		_ = container.Terminate(ctx)
		t.Fatalf("apply migrations: %v", err)
	}

	pool, err := postgres.NewPool(ctx, connectionString)
	if err != nil {
		pgxPool.Close()
		_ = container.Terminate(ctx)
		t.Fatalf("create repository pool: %v", err)
	}

	cleanup := func() {
		pool.Close()
		pgxPool.Close()
		_ = container.Terminate(context.Background())
	}

	return pool, cleanup
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationPaths, err := filepath.Glob(filepath.Join("..", "..", "..", "db", "migrations", "*.up.sql"))
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}

	sort.Strings(migrationPaths)
	for _, migrationPath := range migrationPaths {
		contents, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", migrationPath, err)
		}

		if _, err := pool.Exec(ctx, string(contents)); err != nil {
			return fmt.Errorf("exec migration %s: %w", migrationPath, err)
		}
	}

	return nil
}
