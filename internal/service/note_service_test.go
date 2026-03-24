package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/service"
)

type noteRepositoryStub struct {
	createFn  func(ctx context.Context, title string, body string) (domain.Note, error)
	listFn    func(ctx context.Context) ([]domain.Note, error)
	getByIDFn func(ctx context.Context, id int64) (domain.Note, error)
	updateFn  func(ctx context.Context, id int64, title string, body string) (domain.Note, error)
	deleteFn  func(ctx context.Context, id int64) error
}

func (stub noteRepositoryStub) Create(ctx context.Context, title string, body string) (domain.Note, error) {
	return stub.createFn(ctx, title, body)
}

func (stub noteRepositoryStub) List(ctx context.Context) ([]domain.Note, error) {
	return stub.listFn(ctx)
}

func (stub noteRepositoryStub) GetByID(ctx context.Context, id int64) (domain.Note, error) {
	return stub.getByIDFn(ctx, id)
}

func (stub noteRepositoryStub) Update(ctx context.Context, id int64, title string, body string) (domain.Note, error) {
	return stub.updateFn(ctx, id, title, body)
}

func (stub noteRepositoryStub) Delete(ctx context.Context, id int64) error {
	return stub.deleteFn(ctx, id)
}

func TestCreateNoteTrimsInput(t *testing.T) {
	t.Parallel()

	noteService := service.NewNoteService(noteRepositoryStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{Title: title, Body: body}, nil
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return nil, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string, _ string) (domain.Note, error) {
			return domain.Note{}, nil
		},
		getByIDFn: func(_ context.Context, _ int64) (domain.Note, error) {
			return domain.Note{}, nil
		},
		deleteFn: func(_ context.Context, _ int64) error {
			return nil
		},
	})

	note, err := noteService.CreateNote(context.Background(), "  hello  ", "  world  ")
	if err != nil {
		t.Fatalf("create note: %v", err)
	}

	if note.Title != "hello" {
		t.Fatalf("expected trimmed title, got %q", note.Title)
	}

	if note.Body != "world" {
		t.Fatalf("expected trimmed body, got %q", note.Body)
	}
}

func TestCreateNoteRejectsBlankTitle(t *testing.T) {
	t.Parallel()

	noteService := service.NewNoteService(noteRepositoryStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{}, nil
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return nil, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string, _ string) (domain.Note, error) {
			return domain.Note{}, nil
		},
		getByIDFn: func(_ context.Context, _ int64) (domain.Note, error) {
			return domain.Note{}, nil
		},
		deleteFn: func(_ context.Context, _ int64) error {
			return nil
		},
	})

	_, err := noteService.CreateNote(context.Background(), "   ", "body")
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestListNotesReturnsUnavailableWithoutRepository(t *testing.T) {
	t.Parallel()

	noteService := service.NewNoteService(nil)

	_, err := noteService.ListNotes(context.Background())
	if !errors.Is(err, service.ErrUnavailable) {
		t.Fatalf("expected unavailable error, got %v", err)
	}
}

func TestListNotesReturnsRepositoryData(t *testing.T) {
	t.Parallel()

	createdAt := time.Now().UTC()
	noteService := service.NewNoteService(noteRepositoryStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{}, nil
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return []domain.Note{{ID: 1, Title: "demo", Body: "body", CreatedAt: createdAt, UpdatedAt: createdAt}}, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string, _ string) (domain.Note, error) {
			return domain.Note{}, nil
		},
		getByIDFn: func(_ context.Context, _ int64) (domain.Note, error) {
			return domain.Note{}, nil
		},
		deleteFn: func(_ context.Context, _ int64) error {
			return nil
		},
	})

	notes, err := noteService.ListNotes(context.Background())
	if err != nil {
		t.Fatalf("list notes: %v", err)
	}

	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
}

func TestGetNoteReturnsNotFound(t *testing.T) {
	t.Parallel()

	noteService := service.NewNoteService(noteRepositoryStub{
		createFn: func(_ context.Context, _ string, _ string) (domain.Note, error) { return domain.Note{}, nil },
		listFn:   func(_ context.Context) ([]domain.Note, error) { return nil, nil },
		getByIDFn: func(_ context.Context, _ int64) (domain.Note, error) {
			return domain.Note{}, domain.ErrNotFound
		},
		updateFn: func(_ context.Context, _ int64, _ string, _ string) (domain.Note, error) { return domain.Note{}, nil },
		deleteFn: func(_ context.Context, _ int64) error { return nil },
	})

	_, err := noteService.GetNote(context.Background(), 1)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestUpdateNoteReturnsNotFound(t *testing.T) {
	t.Parallel()

	noteService := service.NewNoteService(noteRepositoryStub{
		createFn: func(_ context.Context, _ string, _ string) (domain.Note, error) { return domain.Note{}, nil },
		listFn:   func(_ context.Context) ([]domain.Note, error) { return nil, nil },
		updateFn: func(_ context.Context, _ int64, _ string, _ string) (domain.Note, error) {
			return domain.Note{}, domain.ErrNotFound
		},
		deleteFn: func(_ context.Context, _ int64) error { return nil },
	})

	_, err := noteService.UpdateNote(context.Background(), 1, "title", "body")
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestDeleteNoteRejectsInvalidID(t *testing.T) {
	t.Parallel()

	noteService := service.NewNoteService(noteRepositoryStub{
		createFn: func(_ context.Context, _ string, _ string) (domain.Note, error) { return domain.Note{}, nil },
		listFn:   func(_ context.Context) ([]domain.Note, error) { return nil, nil },
		updateFn: func(_ context.Context, _ int64, _ string, _ string) (domain.Note, error) { return domain.Note{}, nil },
		deleteFn: func(_ context.Context, _ int64) error { return nil },
	})

	err := noteService.DeleteNote(context.Background(), 0)
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
