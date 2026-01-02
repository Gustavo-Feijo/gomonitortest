package databaseinfra

import (
	"errors"
	"fmt"
	"gomonitor/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// New opens a gorm connection pool to postgres.
func New(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// Create the database instance.
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get the SQL database itself.
	sqlDb, sqlErr := db.DB()

	// Verify if could get the connection.
	if sqlErr != nil {
		return nil, fmt.Errorf("failed to get the sql connection: %v", err)
	}

	// Set the pool values.
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Hour)
	sqlDb.SetConnMaxIdleTime(time.Hour)

	// Test the connection
	if err := sqlDb.Ping(); err != nil {
		dbCloseErr := sqlDb.Close()
		if dbCloseErr != nil {
			err = errors.Join(err, dbCloseErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		dbCloseErr := sqlDb.Close()
		if dbCloseErr != nil {
			err = errors.Join(err, dbCloseErr)
		}
		return nil, fmt.Errorf("failed to use tracing: %w", err)
	}

	return db, err
}
