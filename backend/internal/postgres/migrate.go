package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"who-can-search-ygo/backend/migrations"
)

const migrationLockID int64 = 2026051901

// RunMigrations applies every pending database migration.
func RunMigrations(ctx context.Context, databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open migration database: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping migration database: %w", err)
	}

	if _, err := db.ExecContext(ctx, "SELECT pg_advisory_lock($1)", migrationLockID); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	defer func() {
		_, _ = db.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", migrationLockID)
	}()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}
	if err := validateMigrationVersion(ctx, db); err != nil {
		return err
	}
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

func validateMigrationVersion(ctx context.Context, db *sql.DB) error {
	availableMigrations, err := goose.CollectMigrations(".", 0, goose.MaxVersion)
	if err != nil {
		return fmt.Errorf("collect migrations: %w", err)
	}

	currentVersion, err := goose.GetDBVersionContext(ctx, db)
	if err != nil {
		return fmt.Errorf("get migration version: %w", err)
	}

	var latestAvailableVersion int64
	for _, migration := range availableMigrations {
		if migration.Version > latestAvailableVersion {
			latestAvailableVersion = migration.Version
		}
	}

	if currentVersion > latestAvailableVersion {
		return fmt.Errorf(
			"database migration version %d is newer than application migration version %d",
			currentVersion,
			latestAvailableVersion,
		)
	}
	return nil
}
