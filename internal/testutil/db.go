package testutil

import (
	"context"
	"gomonitor/internal/config"
	databaseinfra "gomonitor/internal/infra/database"
	"log"
	"os"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

const (
	TestPostgresUser     = "test"
	TestPostgresPassword = "test"
	TestPostgresDB       = "testdb"
)

// GetConnections is a singleton implementaudo ssion of the database.
// Return the connection pool.
func NewTestConnection() (*gorm.DB, func()) {

	ctx := context.Background()

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "postgres:18-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     TestPostgresUser,
				"POSTGRES_PASSWORD": TestPostgresPassword,
				"POSTGRES_DB":       TestPostgresDB,
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("Failed to start postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)

	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get container port: %v", err)
	}

	_ = os.Setenv("SQL_HOST", host)
	_ = os.Setenv("SQL_PORT", port.Port())
	_ = os.Setenv("SQL_USER", TestPostgresUser)
	_ = os.Setenv("SQL_PASSWORD", TestPostgresPassword)
	_ = os.Setenv("SQL_DATABASE", TestPostgresDB)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := databaseinfra.New(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}

	if err := databaseinfra.RunMigrations(ctx, cfg.Database, db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	sqlDB, _ := db.DB()
	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}
