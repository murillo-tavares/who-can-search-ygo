package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"who-can-search-ygo/backend/internal/config"
	"who-can-search-ygo/backend/internal/postgres"
)

func main() {
	if err := run(); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.FromEnv()
	if cfg.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := postgres.RunMigrations(ctx, cfg.DatabaseURL); err != nil {
		return err
	}
	slog.Info("database migrations applied")
	return nil
}
