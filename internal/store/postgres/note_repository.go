package postgres

import (
	"context"
	"fmt"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/store/sqlcdb"
)

// NoteRepository provides note persistence backed by PostgreSQL.
type NoteRepository struct {
	queries *sqlcdb.Queries
}

// NewNoteRepository returns a repository when a PostgreSQL pool is available.
// A nil repository is returned when the service is running without a database.
func NewNoteRepository(pool *Pool) *NoteRepository {
	if pool == nil || pool.pool == nil {
		return nil
	}

	return &NoteRepository{queries: sqlcdb.New(pool.pool)}
}

// Create stores a note and returns the persisted record.
func (repository *NoteRepository) Create(ctx context.Context, title string, body string) (domain.Note, error) {
	if repository == nil || repository.queries == nil {
		return domain.Note{}, fmt.Errorf("note repository is unavailable")
	}

	record, err := repository.queries.CreateNote(ctx, sqlcdb.CreateNoteParams{
		Title: title,
		Body:  body,
	})
	if err != nil {
		return domain.Note{}, fmt.Errorf("create note: %w", err)
	}

	return mapNote(record), nil
}

// List returns all notes ordered from newest to oldest.
func (repository *NoteRepository) List(ctx context.Context) ([]domain.Note, error) {
	if repository == nil || repository.queries == nil {
		return nil, fmt.Errorf("note repository is unavailable")
	}

	records, err := repository.queries.ListNotes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}

	notes := make([]domain.Note, 0, len(records))
	for _, record := range records {
		notes = append(notes, mapNote(record))
	}

	return notes, nil
}

func mapNote(record sqlcdb.AppNote) domain.Note {
	return domain.Note{
		ID:        record.ID,
		Title:     record.Title,
		Body:      record.Body,
		CreatedAt: record.CreatedAt.Time,
		UpdatedAt: record.UpdatedAt.Time,
	}
}
