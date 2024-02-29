package main

import (
	"URLShortenePetPrpoject/internal/config"
	"URLShortenePetPrpoject/internal/http-server/handlers/url/delete"
	"URLShortenePetPrpoject/internal/http-server/handlers/url/redirect"
	"URLShortenePetPrpoject/internal/http-server/handlers/url/save"
	"URLShortenePetPrpoject/internal/http-server/middleware/mvLogger"
	"URLShortenePetPrpoject/internal/lib/logger/sl"
	"URLShortenePetPrpoject/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env), slog.String("version", "1"))
	log.Debug("Debug is enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mvLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/url/{alias}", redirect.New(log, storage))
	router.Delete("/url/{alias}", delete.New(log, storage))

	log.Info("Starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

	log.Error("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
