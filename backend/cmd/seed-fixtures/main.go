package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"who-can-search-ygo/backend/internal/config"
	"who-can-search-ygo/backend/internal/postgres"
)

const fixtureDataDir = "testdata/fixtures"

func main() {
	if err := run(); err != nil {
		slog.Error("fixture seed failed", "error", err)
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

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	if err := postgres.SeedFixtures(ctx, pool, fixtureDataDir); err != nil {
		return err
	}
	slog.Info("database fixtures synced", "dir", fixtureDataDir)
	return nil
}
