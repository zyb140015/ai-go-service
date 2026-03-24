package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ai-go-service/internal/domain"
)

var (
	// ErrUnavailable indicates a required dependency is not configured.
	ErrUnavailable = errors.New("service unavailable")
	// ErrInvalidInput indicates request data failed validation.
	ErrInvalidInput = errors.New("invalid input")
	// ErrNotFound indicates the requested note does not exist.
	ErrNotFound = errors.New("note not found")
)

// NoteRepository defines the persistence behavior required by NoteService.
type NoteRepository interface {
	Create(ctx context.Context, title string, body string) (domain.Note, error)
	List(ctx context.Context, options NoteListOptions) ([]domain.Note, int64, error)
	GetByID(ctx context.Context, id int64) (domain.Note, error)
	Update(ctx context.Context, id int64, title string, body string) (domain.Note, error)
	Delete(ctx context.Context, id int64) error
}

// NoteService validates note inputs and delegates persistence to a repository.
type NoteService struct {
	repository NoteRepository
}

// NewNoteService creates a note service with the provided repository.
func NewNoteService(repository NoteRepository) *NoteService {
	return &NoteService{repository: repository}
}

// CreateNote validates and creates a new note.
func (service *NoteService) CreateNote(ctx context.Context, title string, body string) (domain.Note, error) {
	if service == nil || service.repository == nil {
		return domain.Note{}, ErrUnavailable
	}

	trimmedTitle := strings.TrimSpace(title)
	trimmedBody := strings.TrimSpace(body)
	if trimmedTitle == "" {
		return domain.Note{}, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}

	if trimmedBody == "" {
		return domain.Note{}, fmt.Errorf("body is required: %w", ErrInvalidInput)
	}

	note, err := service.repository.Create(ctx, trimmedTitle, trimmedBody)
	if err != nil {
		return domain.Note{}, fmt.Errorf("create note: %w", err)
	}

	return note, nil
}

// ListNotes returns a filtered page of notes ordered by the requested sort.
func (service *NoteService) ListNotes(ctx context.Context, options NoteListOptions) (NoteListResult, error) {
	if service == nil || service.repository == nil {
		return NoteListResult{}, ErrUnavailable
	}

	normalizedOptions, err := normalizeListOptions(options)
	if err != nil {
		return NoteListResult{}, err
	}

	notes, total, err := service.repository.List(ctx, normalizedOptions)
	if err != nil {
		return NoteListResult{}, fmt.Errorf("list notes: %w", err)
	}

	return NoteListResult{
		Items:         notes,
		Total:         total,
		Page:          normalizedOptions.Page,
		PageSize:      normalizedOptions.PageSize,
		Query:         normalizedOptions.Query,
		SortField:     normalizedOptions.SortField,
		SortDirection: normalizedOptions.SortDirection,
	}, nil
}

// GetNote returns one note by ID.
func (service *NoteService) GetNote(ctx context.Context, id int64) (domain.Note, error) {
	if service == nil || service.repository == nil {
		return domain.Note{}, ErrUnavailable
	}

	if id <= 0 {
		return domain.Note{}, fmt.Errorf("id must be positive: %w", ErrInvalidInput)
	}

	note, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Note{}, ErrNotFound
		}

		return domain.Note{}, fmt.Errorf("get note: %w", err)
	}

	return note, nil
}

// UpdateNote validates and updates an existing note.
func (service *NoteService) UpdateNote(ctx context.Context, id int64, title string, body string) (domain.Note, error) {
	if service == nil || service.repository == nil {
		return domain.Note{}, ErrUnavailable
	}

	if id <= 0 {
		return domain.Note{}, fmt.Errorf("id must be positive: %w", ErrInvalidInput)
	}

	trimmedTitle := strings.TrimSpace(title)
	trimmedBody := strings.TrimSpace(body)
	if trimmedTitle == "" {
		return domain.Note{}, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}

	if trimmedBody == "" {
		return domain.Note{}, fmt.Errorf("body is required: %w", ErrInvalidInput)
	}

	note, err := service.repository.Update(ctx, id, trimmedTitle, trimmedBody)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Note{}, ErrNotFound
		}

		return domain.Note{}, fmt.Errorf("update note: %w", err)
	}

	return note, nil
}

// DeleteNote removes an existing note.
func (service *NoteService) DeleteNote(ctx context.Context, id int64) error {
	if service == nil || service.repository == nil {
		return ErrUnavailable
	}

	if id <= 0 {
		return fmt.Errorf("id must be positive: %w", ErrInvalidInput)
	}

	err := service.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("delete note: %w", err)
	}

	return nil
}

func normalizeListOptions(options NoteListOptions) (NoteListOptions, error) {
	normalized := NoteListOptions{
		Page:          options.Page,
		PageSize:      options.PageSize,
		Query:         strings.TrimSpace(options.Query),
		SortField:     strings.TrimSpace(options.SortField),
		SortDirection: strings.TrimSpace(options.SortDirection),
	}

	if normalized.Page == 0 {
		normalized.Page = DefaultNotePage
	}

	if normalized.Page < 1 {
		return NoteListOptions{}, fmt.Errorf("page must be at least 1: %w", ErrInvalidInput)
	}

	if normalized.PageSize == 0 {
		normalized.PageSize = DefaultNotePageSize
	}

	if normalized.PageSize < 1 || normalized.PageSize > MaxNotePageSize {
		return NoteListOptions{}, fmt.Errorf("pageSize must be between 1 and %d: %w", MaxNotePageSize, ErrInvalidInput)
	}

	if normalized.SortField == "" {
		normalized.SortField = NoteSortFieldCreatedAt
	}

	switch normalized.SortField {
	case NoteSortFieldCreatedAt, NoteSortFieldTitle:
	default:
		return NoteListOptions{}, fmt.Errorf("sort must be %q or %q: %w", NoteSortFieldCreatedAt, NoteSortFieldTitle, ErrInvalidInput)
	}

	if normalized.SortDirection == "" {
		normalized.SortDirection = NoteSortDirectionDesc
	}

	switch normalized.SortDirection {
	case NoteSortDirectionAsc, NoteSortDirectionDesc:
	default:
		return NoteListOptions{}, fmt.Errorf("order must be %q or %q: %w", NoteSortDirectionAsc, NoteSortDirectionDesc, ErrInvalidInput)
	}

	return normalized, nil
}
