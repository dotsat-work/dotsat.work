package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	// Set PostgreSQL dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Get migrations subdirectory from embed.FS
	migrationsDir, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations directory: %w", err)
	}

	// Set base filesystem for migrations
	goose.SetBaseFS(migrationsDir)

	// Run migrations
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("migrations completed successfully")
	return nil
}

func MigrateDown(db *sql.DB) error {
	// Set PostgreSQL dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Get migrations subdirectory from embed.FS
	migrationsDir, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations directory: %w", err)
	}

	// Set base filesystem for migrations
	goose.SetBaseFS(migrationsDir)

	// Rollback one migration
	if err := goose.Down(db, "."); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	slog.Info("rolled back one migration")
	return nil
}
