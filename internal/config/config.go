package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database *DatabaseConfig
	HTTP     *HTTPConfig
	Redis    *RedisConfig
}

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

type HTTPConfig struct {
	Address string
}

// Redis configuration.
type RedisConfig struct {
	Addr         string
	Database     int
	MaxRetries   int
	MinIdleConns int
	Password     string
	PoolSize     int
	PoolTimeout  time.Duration
}

// Load get all necessary configuration values.
func Load() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}

	return &Config{
		Database: getDatabaseConfig(),
		HTTP:     getHTTPConfig(),
		Redis:    getRedisConfig(),
	}, nil
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

// getHTTPConfig gets all necessary environment variables.
func getHTTPConfig() *HTTPConfig {
	port := getEnv("APP_PORT", "8080")
	return &HTTPConfig{
		Address: ":" + port,
	}
}

func getRedisConfig() *RedisConfig {
	host := getEnv("REDIS_HOST", "redis")
	port := getEnv("REDIS_PORT", "6379")

	return &RedisConfig{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     getEnv("REDIS_PASSWORD", ""),
		Database:     getIntEnv("REDIS_DATABASE", 0),
		MaxRetries:   getIntEnv("REDIS_MAX_RETRIES", 3),
		MinIdleConns: getIntEnv("REDIS_MIN_IDLE_CONNS", 10),
		PoolSize:     getIntEnv("REDIS_POOL_SIZE", 100),
		PoolTimeout:  time.Duration(getIntEnv("REDIS_POOL_TIMEOUT_SECONDS", 30)) * time.Second,
	}
}

// getEnv returns a env variable, if empty returns a default value.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getIntEnv returns a integer from the env variables, if invalid returns a default value.
func getIntEnv(env string, defaultVal int) int {
	val := os.Getenv(env)
	if val == "" {
		log.Printf("Using default %d for env: %s", defaultVal, env)
		return defaultVal
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Using default %d for env: %s", defaultVal, env)
		return defaultVal
	}
	return i
}

// loadEnv loads the enviromental values if running outside docker.
func loadEnv() error {
	if os.Getenv("ENVIRONMENT") != "docker" {
		return godotenv.Load()
	}
	return nil
}
