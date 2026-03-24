package postgres

import (
	"context"
	"errors"
	"fmt"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/service"
	"ai-go-service/internal/store/sqlcdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

// List returns one filtered page of notes and the matching total count.
func (repository *NoteRepository) List(ctx context.Context, options service.NoteListOptions) ([]domain.Note, int64, error) {
	if repository == nil || repository.queries == nil {
		return nil, 0, fmt.Errorf("note repository is unavailable")
	}

	queryText := pgtype.Text{}
	if options.Query != "" {
		queryText.String = options.Query
		queryText.Valid = true
	}

	total, err := repository.queries.CountNotes(ctx, queryText)
	if err != nil {
		return nil, 0, fmt.Errorf("count notes: %w", err)
	}

	offsetCount := (options.Page - 1) * options.PageSize
	records, err := repository.queries.ListNotes(ctx, sqlcdb.ListNotesParams{
		QueryText:     queryText,
		SortField:     options.SortField,
		SortDirection: options.SortDirection,
		OffsetCount:   offsetCount,
		LimitCount:    options.PageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list notes: %w", err)
	}

	notes := make([]domain.Note, 0, len(records))
	for _, record := range records {
		notes = append(notes, mapNote(record))
	}

	return notes, total, nil
}

// GetByID returns one note by its identifier.
func (repository *NoteRepository) GetByID(ctx context.Context, id int64) (domain.Note, error) {
	if repository == nil || repository.queries == nil {
		return domain.Note{}, fmt.Errorf("note repository is unavailable")
	}

	record, err := repository.queries.GetNoteByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Note{}, domain.ErrNotFound
		}

		return domain.Note{}, fmt.Errorf("get note by id: %w", err)
	}

	return mapNote(record), nil
}

// Update changes a note and returns the updated record.
func (repository *NoteRepository) Update(ctx context.Context, id int64, title string, body string) (domain.Note, error) {
	if repository == nil || repository.queries == nil {
		return domain.Note{}, fmt.Errorf("note repository is unavailable")
	}

	record, err := repository.queries.UpdateNote(ctx, sqlcdb.UpdateNoteParams{
		ID:    id,
		Title: title,
		Body:  body,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Note{}, domain.ErrNotFound
		}

		return domain.Note{}, fmt.Errorf("update note: %w", err)
	}

	return mapNote(record), nil
}

// Delete removes a note by ID.
func (repository *NoteRepository) Delete(ctx context.Context, id int64) error {
	if repository == nil || repository.queries == nil {
		return fmt.Errorf("note repository is unavailable")
	}

	rowsAffected, err := repository.queries.DeleteNote(ctx, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
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
