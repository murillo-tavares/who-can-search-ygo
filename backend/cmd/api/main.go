package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"who-can-search-ygo/backend/internal/config"
	"who-can-search-ygo/backend/internal/httpapi"
	"who-can-search-ygo/backend/internal/postgres"
	"who-can-search-ygo/backend/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.FromEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo, cleanup, err := openRepository(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	svc := service.New(repo)
	handler := httpapi.NewHandler(svc)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("api listening", "addr", cfg.Addr)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func openRepository(ctx context.Context, cfg config.Config) (service.Repository, func(), error) {
	if cfg.DatabaseURL == "" {
		return nil, func() {}, errors.New("DATABASE_URL is required")
	}
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, func() {}, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, func() {}, err
	}
	return postgres.New(pool), pool.Close, nil
}
