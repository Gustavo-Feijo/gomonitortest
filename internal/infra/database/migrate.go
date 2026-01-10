package databaseinfra

import (
	"context"
	"fmt"
	"gomonitor/internal/config"
	"net/url"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations applies all pending migrations to the database.
func RunMigrations(ctx context.Context, cfg *config.DatabaseConfig, db *gorm.DB) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return fmt.Errorf("error getting the sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	absPath, _ := filepath.Abs(cfg.MigrationsPath)
	u := url.URL{
		Scheme: "file",
		Path:   absPath,
	}

	m, err := migrate.NewWithDatabaseInstance(
		u.String(),
		cfg.Database,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}
