package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Admin    *AdminConfig
	Auth     *AuthConfig
	Database *DatabaseConfig
	HTTP     *HTTPConfig
	Logging  *LoggingConfig
	Redis    *RedisConfig
	Tracing  *TracingConfig
}

// Load get all necessary configuration values.
func Load() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}

	adminConfig, err := getAdminConfig()
	if err != nil {
		return nil, err
	}

	authConfig, err := getAuthConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Admin:    adminConfig,
		Auth:     authConfig,
		Database: getDatabaseConfig(),
		HTTP:     getHTTPConfig(),
		Logging:  getLoggingConfig(),
		Redis:    getRedisConfig(),
		Tracing:  getTracingConfig(),
	}, nil
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
	if os.Getenv("ENVIRONMENT") == "" {
		return godotenv.Load()
	}
	return nil
}
