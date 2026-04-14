package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rishavtarway/taskflow/internal/config"
	"github.com/rishavtarway/taskflow/internal/db"
	"github.com/rishavtarway/taskflow/internal/handlers"
	"github.com/rishavtarway/taskflow/internal/middleware"
	"github.com/rishavtarway/taskflow/internal/models"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.Close()

	userModel := models.NewUserModel(database)
	projectModel := models.NewProjectModel(database)
	taskModel := models.NewTaskModel(database)

	authHandler := handlers.NewAuthHandler(userModel, cfg)
	projectHandler := handlers.NewProjectHandler(projectModel)
	taskHandler := handlers.NewTaskHandler(taskModel, projectModel)

	r := chi.NewRouter()

	allowedOrigins := "*"
	if cfg.Environment == "production" {
		allowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	}
	r.Use(middleware.CORS(allowedOrigins))
	r.Use(middleware.Logging(logger))

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Route("/projects", func(r chi.Router) {
		r.Use(middleware.Auth(cfg))
		r.Get("/", projectHandler.ListProjects)
		r.Post("/", projectHandler.CreateProject)

		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", projectHandler.GetProject)
			r.Patch("/", projectHandler.UpdateProject)
			r.Delete("/", projectHandler.DeleteProject)

			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", taskHandler.ListTasks)
				r.Post("/", taskHandler.CreateTask)

				r.Route("/{taskID}", func(r chi.Router) {
					r.Patch("/", taskHandler.UpdateTask)
					r.Delete("/", taskHandler.DeleteTask)
				})
			})
		})
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.Auth(cfg))
		r.Route("/{taskID}", func(r chi.Router) {
			r.Patch("/", taskHandler.UpdateTask)
			r.Delete("/", taskHandler.DeleteTask)
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("server starting", slog.String("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	logger.Info("server exited")
}
