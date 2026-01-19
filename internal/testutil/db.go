package testutil

import (
	"context"
	"gomonitor/internal/config"
	databaseinfra "gomonitor/internal/infra/database"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

const (
	TestPostgresUser     = "test"
	TestPostgresPassword = "test"
	TestPostgresDB       = "testdb"
)

// NewTestDBConnection creates a testcontainer, sets OS envs and return a gorm Connection.
// Should only be used if need to reuse container in multiple tests (Like in TestMain), else, use StartTestDB and connect as normal.
func NewTestDBConnection() (*gorm.DB, func(), error) {
	ctx := context.Background()

	_, host, port, containerCleanup, err := startDB(ctx)

	_ = os.Setenv("SQL_HOST", host)
	_ = os.Setenv("SQL_PORT", port)
	_ = os.Setenv("SQL_USER", TestPostgresUser)
	_ = os.Setenv("SQL_PASSWORD", TestPostgresPassword)
	_ = os.Setenv("SQL_DATABASE", TestPostgresDB)

	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}

	db, err := databaseinfra.New(ctx, cfg.Database)
	if err != nil {
		return nil, nil, err
	}

	sqlDB, _ := db.DB()
	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
		if err := containerCleanup(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}

	if err := databaseinfra.RunMigrations(ctx, cfg.Database, db); err != nil {
		return nil, cleanup, err
	}

	return db, cleanup, nil
}

// startRedis creates a new redis test container.
// Returns the container, host, port, the cleanup function and any error.
func startDB(ctx context.Context) (*postgres.PostgresContainer, string, string, func(ctx context.Context) error, error) {
	container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase(TestPostgresDB),
		postgres.WithUsername(TestPostgresUser),
		postgres.WithPassword(TestPostgresPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("Failed to start postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", "", nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", "", nil, err
	}

	return container, host, port.Port(), func(ctx context.Context) error {
		return container.Terminate(ctx)
	}, nil
}

// StartTestDB creates a new redis test container, set up cleanup and environment variables.
func StartTestDB(t *testing.T) *postgres.PostgresContainer {
	ctx := context.Background()

	container, host, port, cleanup, err := startDB(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanup(t.Context())
	})

	setupTestDBEnv(t, host, port)

	return container
}

func setupTestDBEnv(t *testing.T, host string, port string) {
	t.Setenv("SQL_HOST", host)
	t.Setenv("SQL_PORT", port)
	t.Setenv("SQL_USER", TestPostgresUser)
	t.Setenv("SQL_PASSWORD", TestPostgresPassword)
	t.Setenv("SQL_DATABASE", TestPostgresDB)
}
