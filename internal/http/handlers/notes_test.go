package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai-go-service/internal/domain"
	httpserver "ai-go-service/internal/http"
	"ai-go-service/internal/service"
)

type noteServiceStub struct {
	createFn func(ctx context.Context, title string, body string) (domain.Note, error)
	listFn   func(ctx context.Context) ([]domain.Note, error)
}

func (stub noteServiceStub) CreateNote(ctx context.Context, title string, body string) (domain.Note, error) {
	return stub.createFn(ctx, title, body)
}

func (stub noteServiceStub) ListNotes(ctx context.Context) ([]domain.Note, error) {
	return stub.listFn(ctx)
}

func TestCreateNoteReturnsCreatedNote(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, time.March, 24, 10, 0, 0, 0, time.UTC)
	router := httpserver.NewRouter(newTestLogger(), nil, noteServiceStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{ID: 1, Title: title, Body: body, CreatedAt: createdAt, UpdatedAt: createdAt}, nil
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return nil, nil
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/notes/", strings.NewReader(`{"title":"hello","body":"world"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["title"] != "hello" {
		t.Fatalf("expected title %q, got %#v", "hello", body["title"])
	}
}

func TestCreateNoteRejectsInvalidRequest(t *testing.T) {
	t.Parallel()

	router := httpserver.NewRouter(newTestLogger(), nil, noteServiceStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{}, service.ErrInvalidInput
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return nil, nil
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/notes/", strings.NewReader(`{"title":"","body":"world"}`))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestListNotesReturnsItems(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, time.March, 24, 10, 0, 0, 0, time.UTC)
	router := httpserver.NewRouter(newTestLogger(), nil, noteServiceStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{}, nil
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return []domain.Note{{ID: 1, Title: "hello", Body: "world", CreatedAt: createdAt, UpdatedAt: createdAt}}, nil
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/notes/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(body.Items))
	}
}

func TestListNotesReturnsServiceUnavailable(t *testing.T) {
	t.Parallel()

	router := httpserver.NewRouter(newTestLogger(), nil, noteServiceStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{}, nil
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return nil, service.ErrUnavailable
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/notes/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, recorder.Code)
	}
}

func TestCreateNoteReturnsInternalErrorOnUnexpectedFailure(t *testing.T) {
	t.Parallel()

	router := httpserver.NewRouter(newTestLogger(), nil, noteServiceStub{
		createFn: func(_ context.Context, title string, body string) (domain.Note, error) {
			return domain.Note{}, errors.New("unexpected")
		},
		listFn: func(_ context.Context) ([]domain.Note, error) {
			return nil, nil
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/notes/", strings.NewReader(`{"title":"hello","body":"world"}`))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
}
