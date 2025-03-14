package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/config"
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/handler"
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/repository"
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/service"
)

// App represents the application
type App struct {
	config     *config.Config
	db         *sql.DB
	router     *mux.Router
	httpServer *http.Server
	logger     *log.Logger
}

// New creates a new App instance
func New(cfg *config.Config, db *sql.DB, logger *log.Logger) *App {
	app := &App{
		config: cfg,
		db:     db,
		router: mux.NewRouter(),
		logger: logger,
	}

	app.setupRoutes()
	return app
}

// setupRoutes initializes all the routes for the application
func (a *App) setupRoutes() {
	// Create repositories
	userRepo := repository.NewPostgresUserRepository(a.db)

	// Create services
	userService := service.NewUserService(userRepo, a.logger)

	// Create handlers
	userHandler := handler.NewUserHandler(userService)

	// Register routes
	userHandler.RegisterRoutes(a.router)

	// Add middleware
	a.router.Use(loggingMiddleware(a.logger))

	// Health check endpoint
	a.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// Home page
	a.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Welcome to Go PostgreSQL CRUD API"}`))
	}).Methods(http.MethodGet)
}

// Start starts the HTTP server
func (a *App) Start() error {
	addr := fmt.Sprintf(":%s", a.config.Server.Port)

	a.httpServer = &http.Server{
		Addr:         addr,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}

	return a.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}

// loggingMiddleware logs all requests
func loggingMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Printf(
				"%s %s %s %s",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				time.Since(start),
			)
		})
	}
}
