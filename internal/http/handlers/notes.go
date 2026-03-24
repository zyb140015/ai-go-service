package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"ai-go-service/internal/domain"
	"ai-go-service/internal/http/response"
	"ai-go-service/internal/service"
)

// NoteService defines the note use cases exposed to the HTTP layer.
type NoteService interface {
	CreateNote(ctx context.Context, title string, body string) (domain.Note, error)
	ListNotes(ctx context.Context) ([]domain.Note, error)
}

// CreateNoteRequest is the expected JSON payload for note creation.
type CreateNoteRequest struct {
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

		note, err := noteService.CreateNote(r.Context(), request.Title, request.Body)
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

		notes, err := noteService.ListNotes(r.Context())
		if err != nil {
			handleNoteError(w, err)
			return
		}

		items := make([]NoteResponse, 0, len(notes))
		for _, note := range notes {
			items = append(items, toNoteResponse(note))
		}

		response.JSON(w, http.StatusOK, NotesResponse{Items: items})
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

	response.Error(w, http.StatusInternalServerError, response.CodeInternal, "internal server error")
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
