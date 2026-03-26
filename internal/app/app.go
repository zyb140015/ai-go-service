package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"ai-go-service/internal/auth"
	"ai-go-service/internal/config"
	httpserver "ai-go-service/internal/http"
	"ai-go-service/internal/integration/goadmin"
	"ai-go-service/internal/observability"
	"ai-go-service/internal/service"
	"ai-go-service/internal/store/postgres"
)

// App wires configuration, logging, and the HTTP server together.
type App struct {
	config config.Config
	logger *slog.Logger
	server *http.Server
	db     *postgres.Pool
}

// New creates a new application instance from the current environment.
func New(ctx context.Context) (*App, error) {
	appConfig, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger, err := observability.NewLogger(appConfig.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	db, err := postgres.NewPool(ctx, appConfig.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	noteRepository := postgres.NewNoteRepository(db)
	noteService := service.NewNoteService(noteRepository)
	desktopDataRepository := postgres.NewDesktopDataRepository(db)
	desktopDataService := service.NewDesktopDataService(desktopDataRepository)
	userRepository := postgres.NewUserRepository(db)
	refreshTokenRepository := postgres.NewRefreshTokenRepository(db)
	goAdminClient := goadmin.NewClient(appConfig.GoAdminBaseURL, appConfig.GoAdminTimeout)
	desktopAuthService := service.NewDesktopAuthService(goAdminClient)
	desktopMenuService := service.NewDesktopMenuService(goAdminClient, desktopDataRepository)
	if appConfig.EnableDesktopSeed {
		if err := desktopDataService.EnsureSeedData(ctx); err != nil {
			return nil, fmt.Errorf("ensure desktop seed data: %w", err)
		}
	}

	var authService *service.AuthService
	if userRepository != nil && refreshTokenRepository != nil && appConfig.AuthTokenSecret != "" {
		tokenManager, err := auth.NewTokenManager(appConfig.AuthTokenSecret, nil)
		if err != nil {
			return nil, fmt.Errorf("create token manager: %w", err)
		}

		authService = service.NewAuthService(userRepository, refreshTokenRepository, tokenManager, appConfig.AuthTokenTTL)
	}

	server := &http.Server{
		Addr:              appConfig.HTTPAddr,
		Handler:           httpserver.NewRouter(logger, db, noteService, authService, desktopAuthService, desktopDataService, desktopMenuService),
		ReadTimeout:       appConfig.ReadTimeout,
		ReadHeaderTimeout: appConfig.ReadHeaderTimeout,
		WriteTimeout:      appConfig.WriteTimeout,
		IdleTimeout:       appConfig.IdleTimeout,
	}

	return &App{
		config: appConfig,
		logger: logger,
		server: server,
		db:     db,
	}, nil
}

// Run starts the HTTP server and shuts it down gracefully when the context is canceled.
func (app *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		app.logger.Info("starting http server", "addr", app.server.Addr)

		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("listen and serve: %w", err)
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), app.config.ShutdownTimeout)
		defer cancel()
		defer app.db.Close()

		app.logger.Info("shutting down http server")

		if err := app.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		return nil
	case err := <-errCh:
		app.db.Close()
		return err
	}
}
