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
func NewRouter(logger *slog.Logger, readinessChecker handlers.ReadinessChecker, noteService handlers.NoteService, authService handlers.AuthService, desktopAuthService handlers.DesktopAuthUseCase, desktopDataService handlers.DesktopDataUseCase, desktopMenuServices ...handlers.DesktopMenuUseCase) http.Handler {
	router := chi.NewRouter()

	var desktopMenuService handlers.DesktopMenuUseCase
	if len(desktopMenuServices) > 0 {
		desktopMenuService = desktopMenuServices[0]
	}

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(projectmiddleware.CORS)
	router.Use(projectmiddleware.RequestLogger(logger))
	router.Use(projectmiddleware.Recoverer(logger))
	router.Use(chimiddleware.Timeout(requestTimeout))

	router.Get("/healthz", handlers.Healthz)
	router.Get("/readyz", handlers.ReadyzHandler(readinessChecker))
	router.Get("/openapi.yaml", handlers.OpenAPIHandler)
	router.Get("/docs", handlers.SwaggerUIHandler)
	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/register", handlers.RegisterHandler(authService))
		authRouter.Post("/login", handlers.LoginHandler(authService))
		authRouter.Post("/refresh", handlers.RefreshHandler(authService))
		authRouter.Post("/logout", handlers.LogoutHandler(authService))
		authRouter.Post("/change-password", handlers.ChangePasswordHandler(authService))
		authRouter.Get("/me", handlers.CurrentUserHandler(authService))
	})
	router.Route("/desktop", func(desktopRouter chi.Router) {
		desktopRouter.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/login", handlers.DesktopLoginHandler(desktopAuthService))
			authRouter.Post("/refresh", handlers.DesktopRefreshHandler(desktopAuthService))
			authRouter.Get("/me", handlers.DesktopCurrentUserHandler(desktopAuthService))
			authRouter.Put("/profile", handlers.DesktopUpdateProfileHandler(desktopAuthService))
			authRouter.Post("/avatar", handlers.DesktopUploadAvatarHandler(desktopAuthService))
			authRouter.Put("/password", handlers.DesktopUpdatePasswordHandler(desktopAuthService))
		})
		desktopRouter.Route("/menus", func(menuRouter chi.Router) {
			menuRouter.Get("/current", handlers.DesktopCurrentMenusHandler(desktopMenuService))
			menuRouter.Get("/manage", handlers.DesktopManagedMenusHandler(desktopMenuService))
			menuRouter.Post("/manage", handlers.DesktopCreateManagedMenuHandler(desktopMenuService))
			menuRouter.Put("/manage", handlers.DesktopUpdateManagedMenuHandler(desktopMenuService))
			menuRouter.Post("/manage/delete", handlers.DesktopDeleteManagedMenuHandler(desktopMenuService))
		})
		desktopRouter.Get("/messages", handlers.DesktopMessagesHandler(desktopDataService))
		desktopRouter.Post("/messages/read", handlers.DesktopMarkMessageReadHandler(desktopDataService))
		desktopRouter.Post("/messages/read-batch", handlers.DesktopMarkMessagesReadHandler(desktopDataService))
		desktopRouter.Get("/announcements", handlers.DesktopAnnouncementsHandler(desktopDataService))
		desktopRouter.Post("/announcements", handlers.DesktopCreateAnnouncementHandler(desktopDataService))
		desktopRouter.Put("/announcements", handlers.DesktopUpdateAnnouncementHandler(desktopDataService))
		desktopRouter.Post("/announcements/publish", handlers.DesktopPublishAnnouncementHandler(desktopDataService))
		desktopRouter.Post("/announcements/delete", handlers.DesktopDeleteAnnouncementHandler(desktopDataService))
		desktopRouter.Get("/dicts", handlers.DesktopDictsHandler(desktopDataService))
		desktopRouter.Post("/dicts", handlers.DesktopCreateDictHandler(desktopDataService))
		desktopRouter.Put("/dicts", handlers.DesktopUpdateDictHandler(desktopDataService))
		desktopRouter.Post("/dicts/delete", handlers.DesktopDeleteDictHandler(desktopDataService))
		desktopRouter.Post("/dicts/items", handlers.DesktopCreateDictItemHandler(desktopDataService))
		desktopRouter.Put("/dicts/items", handlers.DesktopUpdateDictItemHandler(desktopDataService))
		desktopRouter.Post("/dicts/items/delete", handlers.DesktopDeleteDictItemHandler(desktopDataService))
		desktopRouter.Get("/orgs", handlers.DesktopOrgsHandler(desktopDataService))
		desktopRouter.Post("/orgs", handlers.DesktopCreateOrgHandler(desktopDataService))
		desktopRouter.Put("/orgs", handlers.DesktopUpdateOrgHandler(desktopDataService))
		desktopRouter.Post("/orgs/delete", handlers.DesktopDeleteOrgHandler(desktopDataService))
		desktopRouter.Get("/menu-templates", handlers.DesktopMenuTemplatesHandler(desktopDataService))
		desktopRouter.Post("/menu-templates", handlers.DesktopCreateMenuTemplateHandler(desktopDataService))
		desktopRouter.Put("/menu-templates", handlers.DesktopUpdateMenuTemplateHandler(desktopDataService))
		desktopRouter.Post("/menu-templates/delete", handlers.DesktopDeleteMenuTemplateHandler(desktopDataService))
		desktopRouter.Get("/tenants", handlers.DesktopTenantsHandler(desktopDataService))
		desktopRouter.Post("/tenants", handlers.DesktopCreateTenantHandler(desktopDataService))
		desktopRouter.Put("/tenants", handlers.DesktopUpdateTenantHandler(desktopDataService))
		desktopRouter.Post("/tenants/delete", handlers.DesktopDeleteTenantHandler(desktopDataService))
		desktopRouter.Get("/users", handlers.DesktopUsersHandler(desktopDataService))
		desktopRouter.Post("/users", handlers.DesktopCreateUserHandler(desktopDataService))
		desktopRouter.Put("/users", handlers.DesktopUpdateUserHandler(desktopDataService))
		desktopRouter.Post("/users/delete", handlers.DesktopDeleteUserHandler(desktopDataService))
		desktopRouter.Post("/users/reset-password", handlers.DesktopResetUserPasswordHandler(desktopDataService))
		desktopRouter.Get("/roles", handlers.DesktopRolesHandler(desktopDataService))
		desktopRouter.Post("/roles", handlers.DesktopCreateRoleHandler(desktopDataService))
		desktopRouter.Put("/roles", handlers.DesktopUpdateRoleHandler(desktopDataService))
		desktopRouter.Post("/roles/delete", handlers.DesktopDeleteRoleHandler(desktopDataService))
		desktopRouter.Get("/roles/data-permission", handlers.DesktopRoleDataPermissionHandler(desktopDataService))
		desktopRouter.Put("/roles/data-permission", handlers.DesktopSaveRoleDataPermissionHandler(desktopDataService))
		desktopRouter.Get("/settings", handlers.DesktopSystemSettingsHandler(desktopDataService))
		desktopRouter.Put("/settings", handlers.DesktopSaveSystemSettingsHandler(desktopDataService))
		desktopRouter.Get("/logs/login", handlers.DesktopLoginLogsHandler(desktopDataService))
		desktopRouter.Get("/logs/operation", handlers.DesktopOperationLogsHandler(desktopDataService))
		desktopRouter.Get("/monitor", handlers.DesktopMonitorHandler(desktopDataService))
		desktopRouter.Post("/monitor/status", handlers.DesktopUpdateMonitorStatusHandler(desktopDataService))
		desktopRouter.Get("/stats/system", handlers.DesktopSystemStatsHandler(desktopDataService))
		desktopRouter.Get("/stats/usage", handlers.DesktopUsageStatsHandler(desktopDataService))
	})
	router.Route("/notes", func(notesRouter chi.Router) {
		notesRouter.Use(projectmiddleware.RequireAuth(authService))
		notesRouter.Get("/", handlers.ListNotesHandler(noteService))
		notesRouter.Post("/", handlers.CreateNoteHandler(noteService))
		notesRouter.Get("/{noteID}", handlers.GetNoteHandler(noteService))
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
