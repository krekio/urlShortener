package main

import (
	"github.com/krekio/urlShortener.git/internal/config"
	"log/slog"
	"os"
)

const (
	envlocal = "local"
	envdev   = "dev"
)

func main() {
	cfg := config.NewConfig()
	log := setupLogger(cfg.Env)
	log.Info("start", slog.String("env", cfg.Env))
	log.Debug("debug enabled")
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
