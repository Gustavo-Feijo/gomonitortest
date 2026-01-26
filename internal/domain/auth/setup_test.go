package auth_test

import (
	"context"
	"gomonitor/internal/config"
	databaseinfra "gomonitor/internal/infra/database"
	"gomonitor/internal/testutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gorm"
)

var (
	testDbCfg *config.DatabaseConfig
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	_, host, port, containerCleanup, err := testutil.StartDB(ctx)
	if err != nil {
		log.Fatalf("error starting database container: %v", err)
	}
	testDbCfg = &config.DatabaseConfig{
		Database:       testutil.TestPostgresDB,
		Password:       testutil.TestPostgresPassword,
		User:           testutil.TestPostgresUser,
		Host:           host,
		Port:           port,
		MigrationsPath: "migrations",
	}
	if !config.IsProduction() {
		projectRoot := config.FindProjectRoot()
		if projectRoot == "" {
			log.Fatal("Error finding project root")
		}
		testDbCfg.MigrationsPath = filepath.Join(projectRoot, "migrations")
	}

	dbConn, err := databaseinfra.New(ctx, testDbCfg)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}

	if err := databaseinfra.RunMigrations(ctx, testDbCfg, dbConn); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}

	code := m.Run()
	_ = containerCleanup(ctx)
	os.Exit(code)
}

func setupTx(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })
	return tx
}
