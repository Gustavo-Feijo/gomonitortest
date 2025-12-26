package config

import "fmt"

// Database configuration.
type DatabaseConfig struct {
	Database     string
	Host         string
	MaxOpenConns int
	MaxIdleConns int
	Password     string
	Port         string
	User         string
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		d.Host, d.User, d.Password, d.Database, d.Port,
	)
}

// getDatabaseConfig gets all necessary environment variables.
func getDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Database:     getEnv("SQL_DATABASE", "default"),
		Host:         getEnv("SQL_HOST", "postgresql"),
		MaxOpenConns: getIntEnv("SQL_MAX_OPEN_CONNS", 400),
		MaxIdleConns: getIntEnv("SQL_MAX_IDLE_CONNS", 10),
		Password:     getEnv("SQL_PASSWORD", ""),
		Port:         getEnv("SQL_PORT", "5432"),
		User:         getEnv("SQL_USER", "root"),
	}
}
