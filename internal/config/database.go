package config

import (
	"fmt"
	"log"
	"path/filepath"
)

// Database configuration.
type DatabaseConfig struct {
	Database       string
	Host           string
	MaxOpenConns   int
	MaxIdleConns   int
	MigrationsPath string
	Password       string
	Port           string
	User           string
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		d.Host, d.User, d.Password, d.Database, d.Port,
	)
}

// getDatabaseConfig gets all necessary environment variables.
func getDatabaseConfig() *DatabaseConfig {
	dbConfig := &DatabaseConfig{
		Database:       getEnv("SQL_DATABASE", "default"),
		Host:           getEnv("SQL_HOST", "postgresql"),
		MaxOpenConns:   getIntEnv("SQL_MAX_OPEN_CONNS", 400),
		MaxIdleConns:   getIntEnv("SQL_MAX_IDLE_CONNS", 10),
		MigrationsPath: getEnv("SQL_MIGRATION_PATH", "migrations"),
		Password:       getEnv("SQL_PASSWORD", ""),
		Port:           getEnv("SQL_PORT", "5432"),
		User:           getEnv("SQL_USER", "root"),
	}

	// Handle local testing outside docker.
	if !IsProduction() {
		projectRoot := FindProjectRoot()
		if projectRoot == "" {
			log.Fatal("Error finding project root")
		}
		dbConfig.MigrationsPath = filepath.Join(projectRoot, dbConfig.MigrationsPath)
	}

	return dbConfig
}
