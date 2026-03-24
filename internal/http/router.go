package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"ai-go-service/internal/http/handlers"
	projectmiddleware "ai-go-service/internal/http/middleware"
	"ai-go-service/internal/http/response"
)

const requestTimeout = 15 * time.Second

// NewRouter builds the HTTP router and registers all public endpoints.
func NewRouter(logger *slog.Logger, readinessChecker handlers.ReadinessChecker, noteService handlers.NoteService) http.Handler {
	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(projectmiddleware.RequestLogger(logger))
	router.Use(projectmiddleware.Recoverer(logger))
	router.Use(chimiddleware.Timeout(requestTimeout))

	router.Get("/healthz", handlers.Healthz)
	router.Get("/readyz", handlers.ReadyzHandler(readinessChecker))
	router.Route("/notes", func(notesRouter chi.Router) {
		notesRouter.Get("/", handlers.ListNotesHandler(noteService))
		notesRouter.Post("/", handlers.CreateNoteHandler(noteService))
		notesRouter.Put("/{noteID}", handlers.UpdateNoteHandler(noteService))
		notesRouter.Delete("/{noteID}", handlers.DeleteNoteHandler(noteService))
	})

	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "resource not found")
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		response.Error(w, http.StatusMethodNotAllowed, response.CodeMethodNotAllowed, "method not allowed")
	})

	return router
}
