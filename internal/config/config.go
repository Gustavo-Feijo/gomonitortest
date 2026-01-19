package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Admin          *AdminConfig
	Auth           *AuthConfig
	CircuitBreaker *CircuitBreakerConfig
	Database       *DatabaseConfig
	HTTP           *HTTPConfig
	Logging        *LoggingConfig
	Redis          *RedisConfig
	Tracing        *TracingConfig
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
		Admin:          adminConfig,
		Auth:           authConfig,
		CircuitBreaker: getCircuitBreakerConfig(),
		Database:       getDatabaseConfig(),
		HTTP:           getHTTPConfig(),
		Logging:        getLoggingConfig(),
		Redis:          getRedisConfig(),
		Tracing:        getTracingConfig(),
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
		projectRoot := findProjectRoot()
		if projectRoot == "" {
			log.Fatal("Error finding project root")
		}

		return loadMissingEnv(filepath.Join(projectRoot, ".env"))
	}

	return nil
}

// loadMissing loads only missing envs, allowing to override some configs (Test containers ports for example)
func loadMissingEnv(filenames ...string) error {
	envs, err := godotenv.Read(filenames...)
	if err != nil {
		return err
	}

	for k, v := range envs {
		if _, exists := os.LookupEnv(k); !exists {
			_ = os.Setenv(k, v)
		}
	}

	return nil
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
