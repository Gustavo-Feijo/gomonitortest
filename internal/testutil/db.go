package testutil

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	TestPostgresUser     = "test"
	TestPostgresPassword = "test"
	TestPostgresDB       = "testdb"
)

// startRedis creates a new redis test container.
// Returns the container, host, port, the cleanup function and any error.
func StartDB(ctx context.Context) (*postgres.PostgresContainer, string, string, func(ctx context.Context) error, error) {
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

	container, host, port, cleanup, err := StartDB(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = cleanup(ctx)
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
