package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/krekio/urlShortener.git/internal/config"
	"github.com/krekio/urlShortener.git/internal/httpserver/handlers/redirect"
	"github.com/krekio/urlShortener.git/internal/httpserver/handlers/url/save"
	"github.com/krekio/urlShortener.git/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
)

const (
	envlocal = "local"
	envdev   = "dev"
)

func main() {
	cfg, httpCfg := config.NewConfig()
	log := setupLogger(cfg.Env)
	log.Info("start", slog.String("env", cfg.Env))
	log.Debug("debug enabled")
	storage, err := postgres.NewStorage(*cfg)
	if err != nil {
		log.Error("error creating storage", err)
		os.Exit(1)
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			httpCfg.User: httpCfg.Password,
		}))
		r.Post("/", save.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))
	log.Info("starting server")
	srv := http.Server{
		Addr:         httpCfg.Address,
		Handler:      router,
		ReadTimeout:  httpCfg.Timeout,
		WriteTimeout: httpCfg.Timeout,
		IdleTimeout:  httpCfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("error starting server", err)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envlocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envdev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
