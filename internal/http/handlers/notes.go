package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"ai-go-service/internal/domain"
	projectmiddleware "ai-go-service/internal/http/middleware"
	"ai-go-service/internal/http/response"
	"ai-go-service/internal/service"
	"github.com/go-chi/chi/v5"
)

// NoteService defines the note use cases exposed to the HTTP layer.
type NoteService interface {
	CreateNote(ctx context.Context, userID int64, title string, body string) (domain.Note, error)
	ListNotes(ctx context.Context, userID int64, options service.NoteListOptions) (service.NoteListResult, error)
	GetNote(ctx context.Context, userID int64, id int64) (domain.Note, error)
	UpdateNote(ctx context.Context, userID int64, id int64, title string, body string) (domain.Note, error)
	DeleteNote(ctx context.Context, userID int64, id int64) error
}

// CreateNoteRequest is the expected JSON payload for note creation.
type CreateNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// UpdateNoteRequest is the expected JSON payload for note updates.
type UpdateNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// NoteResponse is the public JSON shape returned for note resources.
type NoteResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// NotesResponse wraps the note list response in a stable top-level field.
type NotesResponse struct {
	Items []NoteResponse `json:"items"`
	Meta  NotesMeta      `json:"meta"`
}

// NotesMeta describes list-note pagination and sorting metadata.
type NotesMeta struct {
	Page          int32  `json:"page"`
	PageSize      int32  `json:"pageSize"`
	Total         int64  `json:"total"`
	SortField     string `json:"sort"`
	SortDirection string `json:"order"`
	Query         string `json:"q,omitempty"`
}

// CreateNoteHandler returns an HTTP handler that creates notes.
func CreateNoteHandler(noteService NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if noteService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "service is not ready")
			return
		}

		var request CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}

		user, ok := projectmiddleware.CurrentUserFromContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
			return
		}

		note, err := noteService.CreateNote(r.Context(), user.ID, request.Title, request.Body)
		if err != nil {
			handleNoteError(w, err)
			return
		}

		response.JSON(w, http.StatusCreated, toNoteResponse(note))
	}
}

// ListNotesHandler returns an HTTP handler that lists notes.
func ListNotesHandler(noteService NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if noteService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "service is not ready")
			return
		}

		options, ok := parseNoteListOptions(w, r)
		if !ok {
			return
		}

		user, ok := projectmiddleware.CurrentUserFromContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
			return
		}

		result, err := noteService.ListNotes(r.Context(), user.ID, options)
		if err != nil {
			handleNoteError(w, err)
			return
		}

		items := make([]NoteResponse, 0, len(result.Items))
		for _, note := range result.Items {
			items = append(items, toNoteResponse(note))
		}

		response.JSON(w, http.StatusOK, NotesResponse{
			Items: items,
			Meta: NotesMeta{
				Page:          result.Page,
				PageSize:      result.PageSize,
				Total:         result.Total,
				SortField:     result.SortField,
				SortDirection: result.SortDirection,
				Query:         result.Query,
			},
		})
	}
}

// GetNoteHandler returns an HTTP handler that retrieves a single note.
func GetNoteHandler(noteService NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if noteService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "service is not ready")
			return
		}

		noteID, ok := parseNoteID(w, r)
		if !ok {
			return
		}

		user, ok := projectmiddleware.CurrentUserFromContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
			return
		}

		note, err := noteService.GetNote(r.Context(), user.ID, noteID)
		if err != nil {
			handleNoteError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, toNoteResponse(note))
	}
}

// UpdateNoteHandler returns an HTTP handler that updates a note.
func UpdateNoteHandler(noteService NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if noteService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "service is not ready")
			return
		}

		noteID, ok := parseNoteID(w, r)
		if !ok {
			return
		}

		var request UpdateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}

		user, ok := projectmiddleware.CurrentUserFromContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
			return
		}

		note, err := noteService.UpdateNote(r.Context(), user.ID, noteID, request.Title, request.Body)
		if err != nil {
			handleNoteError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, toNoteResponse(note))
	}
}

// DeleteNoteHandler returns an HTTP handler that deletes a note.
func DeleteNoteHandler(noteService NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if noteService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "service is not ready")
			return
		}

		noteID, ok := parseNoteID(w, r)
		if !ok {
			return
		}

		user, ok := projectmiddleware.CurrentUserFromContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication failed")
			return
		}

		if err := noteService.DeleteNote(r.Context(), user.ID, noteID); err != nil {
			handleNoteError(w, err)
			return
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}

func handleNoteError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrInvalidInput) {
		response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrUnavailable) {
		response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "service is not ready")
		return
	}

	if errors.Is(err, service.ErrNotFound) {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "note not found")
		return
	}

	response.Error(w, http.StatusInternalServerError, response.CodeInternal, "internal server error")
}

func parseNoteID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	noteIDValue := chi.URLParam(r, "noteID")
	noteID, err := strconv.ParseInt(noteIDValue, 10, 64)
	if err != nil || noteID <= 0 {
		response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "note id must be a positive integer")
		return 0, false
	}

	return noteID, true
}

func parseNoteListOptions(w http.ResponseWriter, r *http.Request) (service.NoteListOptions, bool) {
	page, ok := parsePositiveInt32Query(w, r, "page")
	if !ok {
		return service.NoteListOptions{}, false
	}

	pageSize, ok := parsePositiveInt32Query(w, r, "pageSize")
	if !ok {
		return service.NoteListOptions{}, false
	}

	return service.NoteListOptions{
		Page:          page,
		PageSize:      pageSize,
		Query:         r.URL.Query().Get("q"),
		SortField:     r.URL.Query().Get("sort"),
		SortDirection: r.URL.Query().Get("order"),
	}, true
}

func parsePositiveInt32Query(w http.ResponseWriter, r *http.Request, key string) (int32, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, true
	}

	parsedValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil || parsedValue < 1 {
		response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, key+" must be a positive integer")
		return 0, false
	}

	return int32(parsedValue), true
}

func toNoteResponse(note domain.Note) NoteResponse {
	return NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Body:      note.Body,
		CreatedAt: note.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
