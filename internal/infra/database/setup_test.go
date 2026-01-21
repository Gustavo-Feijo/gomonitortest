package databaseinfra_test

import (
	"context"
	"gomonitor/internal/config"
	"gomonitor/internal/testutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	testDbCfg *config.DatabaseConfig
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	_, host, port, containerCleanup, _ := testutil.StartDB(ctx)
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

	code := m.Run()
	_ = containerCleanup(ctx)
	os.Exit(code)
}
