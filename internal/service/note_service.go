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
	List(ctx context.Context) ([]domain.Note, error)
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

// ListNotes returns all available notes ordered by the repository.
func (service *NoteService) ListNotes(ctx context.Context) ([]domain.Note, error) {
	if service == nil || service.repository == nil {
		return nil, ErrUnavailable
	}

	notes, err := service.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}

	return notes, nil
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
