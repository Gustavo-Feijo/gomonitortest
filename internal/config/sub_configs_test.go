package config

import (
	"maps"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAuthConfig(t *testing.T) {
	baseEnv := map[string]string{
		"AUTH_ACCESS_TOKEN_SECRET":  "access",
		"AUTH_REFRESH_TOKEN_SECRET": "refresh",
		"AUTH_FAKE_HASH":            "fakehash",
		"AUTH_ACCESS_TOKEN_TTL":     "1h",
		"AUTH_REFRESH_TOKEN_TTL":    "168h",
	}

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "valid config",
			env:  baseEnv,
		},
		{
			name: "missing access token secret",
			env: func() map[string]string {
				m := maps.Clone(baseEnv)
				delete(m, "AUTH_ACCESS_TOKEN_SECRET")
				return m
			}(),
			wantErr: true,
		},
		{
			name: "missing refresh token secret",
			env: func() map[string]string {
				m := maps.Clone(baseEnv)
				delete(m, "AUTH_REFRESH_TOKEN_SECRET")
				return m
			}(),
			wantErr: true,
		},
		{
			name: "missing fake hash",
			env: func() map[string]string {
				m := maps.Clone(baseEnv)
				delete(m, "AUTH_FAKE_HASH")
				return m
			}(),
			wantErr: true,
		},
		{
			name: "invalid access token ttl",
			env: func() map[string]string {
				m := maps.Clone(baseEnv)
				m["AUTH_ACCESS_TOKEN_TTL"] = "not-a-duration"
				return m
			}(),
			wantErr: true,
		},
		{
			name: "invalid refresh token ttl",
			env: func() map[string]string {
				m := maps.Clone(baseEnv)
				m["AUTH_REFRESH_TOKEN_TTL"] = "invalid"
				return m
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := getAuthConfig()

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
				assert.Equal(t, time.Hour, cfg.AccessTokenTTL)
			}
		})
	}
}

func TestGetAdminConfig(t *testing.T) {
	baseEnv := map[string]string{
		"ADMIN_EMAIL":    "email",
		"ADMIN_PASSWORD": "password",
	}

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "valid config",
			env:  baseEnv,
		},
		{
			name: "missing email",
			env: func() map[string]string {
				m := maps.Clone(baseEnv)
				delete(m, "ADMIN_EMAIL")
				return m
			}(),
			wantErr: true,
		},
		{
			name: "missing password",
			env: func() map[string]string {
				m := maps.Clone(baseEnv)
				delete(m, "ADMIN_PASSWORD")
				return m
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := getAdminConfig()

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
			}
		})
	}
}

func TestGetCircuitBreakerConfig(t *testing.T) {
	baseEnv := map[string]string{
		"CIRCUIT_BREAKER_MAX_REQUEST":  "7",
		"CIRCUIT_BREAKER_MAX_FAILURES": "2",
	}

	tests := []struct {
		name           string
		env            map[string]string
		expectedConfig *CircuitBreakerConfig
	}{
		{
			name: "override config",
			env:  baseEnv,
			expectedConfig: &CircuitBreakerConfig{
				MaxRequests: 7,
				MaxFailures: 2,
			},
		},
		{
			name: "default configs",
			env: func() map[string]string {
				return make(map[string]string)
			}(),
			expectedConfig: &CircuitBreakerConfig{
				MaxRequests: uint32(5),
				MaxFailures: uint32(5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := getCircuitBreakerConfig()

			require.NotNil(t, cfg)

			assert.Equal(t, tt.expectedConfig.MaxFailures, cfg.MaxFailures)
			assert.Equal(t, tt.expectedConfig.MaxRequests, cfg.MaxRequests)
		})
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		User:     "user",
		Password: "password",
		Database: "mydb",
		Port:     "5432",
	}

	dsn := cfg.DSN()

	expected := "host=localhost user=user password=password dbname=mydb port=5432 sslmode=disable TimeZone=UTC"
	assert.Equal(t, expected, dsn)
}

func TestGetDatabaseConfig(t *testing.T) {
	baseEnv := map[string]string{
		"SQL_DATABASE":       "testdb",
		"SQL_HOST":           "db-host",
		"SQL_USER":           "db-user",
		"SQL_PASSWORD":       "db-pass",
		"SQL_PORT":           "6543",
		"SQL_MAX_OPEN_CONNS": "100",
		"SQL_MAX_IDLE_CONNS": "20",
		"SQL_MIGRATION_PATH": "custom-migrations",
	}

	tests := []struct {
		name           string
		env            map[string]string
		environmentSet bool
		assertions     func(t *testing.T, cfg *DatabaseConfig)
	}{
		{
			name:           "override all values",
			env:            baseEnv,
			environmentSet: true,
			assertions: func(t *testing.T, cfg *DatabaseConfig) {
				assert.Equal(t, "testdb", cfg.Database)
				assert.Equal(t, "db-host", cfg.Host)
				assert.Equal(t, "db-user", cfg.User)
				assert.Equal(t, "db-pass", cfg.Password)
				assert.Equal(t, "6543", cfg.Port)
				assert.Equal(t, 100, cfg.MaxOpenConns)
				assert.Equal(t, 20, cfg.MaxIdleConns)
				assert.Equal(t, "custom-migrations", cfg.MigrationsPath)
			},
		},
		{
			name:           "default values when env is missing",
			env:            map[string]string{},
			environmentSet: true,
			assertions: func(t *testing.T, cfg *DatabaseConfig) {
				assert.Equal(t, "default", cfg.Database)
				assert.Equal(t, "postgresql", cfg.Host)
				assert.Equal(t, "root", cfg.User)
				assert.Equal(t, "", cfg.Password)
				assert.Equal(t, "5432", cfg.Port)
				assert.Equal(t, 400, cfg.MaxOpenConns)
				assert.Equal(t, 10, cfg.MaxIdleConns)
				assert.Equal(t, "migrations", cfg.MigrationsPath)
			},
		},
		{
			name:           "local environment resolves migrations path from project root",
			env:            baseEnv,
			environmentSet: false,
			assertions: func(t *testing.T, cfg *DatabaseConfig) {
				projectRoot := FindProjectRoot()
				require.NotEmpty(t, projectRoot)

				expectedPath := filepath.Join(projectRoot, "custom-migrations")
				assert.Equal(t, expectedPath, cfg.MigrationsPath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			// Treat as production so doesn't load any .env file.
			if tt.environmentSet {
				t.Setenv("ENVIRONMENT", "production")
			}

			cfg := getDatabaseConfig()

			require.NotNil(t, cfg)
			tt.assertions(t, cfg)
		})
	}
}

func TestGetHTTPConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "default port",
			env:      map[string]string{},
			expected: ":8080",
		},
		{
			name: "custom port",
			env: map[string]string{
				"APP_PORT": "9090",
			},
			expected: ":9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := getHTTPConfig()

			require.NotNil(t, cfg)
			assert.Equal(t, tt.expected, cfg.Address)
		})
	}
}

func TestGetLoggingConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "default log level",
			env:      map[string]string{},
			expected: "debug",
		},
		{
			name: "custom log level",
			env: map[string]string{
				"LOG_LEVEL": "info",
			},
			expected: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := getLoggingConfig()

			require.NotNil(t, cfg)
			assert.Equal(t, tt.expected, cfg.Level)
		})
	}
}

func TestGetRedisConfig(t *testing.T) {
	tests := []struct {
		name           string
		env            map[string]string
		expectedAddr   string
		expectedConfig *RedisConfig
	}{
		{
			name: "explicit redis addr",
			env: map[string]string{
				"REDIS_ADDR": "localhost:9999",
			},
			expectedAddr: "localhost:9999",
			expectedConfig: &RedisConfig{
				Database:     0,
				MaxRetries:   3,
				MinIdleConns: 10,
				PoolSize:     100,
				PoolTimeout:  30 * time.Second,
				Password:     "",
			},
		},
		{
			name: "host and port fallback",
			env: map[string]string{
				"REDIS_HOST": "redis-host",
				"REDIS_PORT": "6380",
			},
			expectedAddr: "redis-host:6380",
			expectedConfig: &RedisConfig{
				Database:     0,
				MaxRetries:   3,
				MinIdleConns: 10,
				PoolSize:     100,
				PoolTimeout:  30 * time.Second,
				Password:     "",
			},
		},
		{
			name: "override all numeric values",
			env: map[string]string{
				"REDIS_DATABASE":             "2",
				"REDIS_MAX_RETRIES":          "5",
				"REDIS_MIN_IDLE_CONNS":       "3",
				"REDIS_POOL_SIZE":            "50",
				"REDIS_POOL_TIMEOUT_SECONDS": "5",
			},
			expectedAddr: "redis:6379",
			expectedConfig: &RedisConfig{
				Database:     2,
				MaxRetries:   5,
				MinIdleConns: 3,
				PoolSize:     50,
				PoolTimeout:  5 * time.Second,
				Password:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := getRedisConfig()

			require.NotNil(t, cfg)
			assert.Equal(t, tt.expectedAddr, cfg.Addr)
			assert.Equal(t, tt.expectedConfig.Database, cfg.Database)
			assert.Equal(t, tt.expectedConfig.MaxRetries, cfg.MaxRetries)
			assert.Equal(t, tt.expectedConfig.MinIdleConns, cfg.MinIdleConns)
			assert.Equal(t, tt.expectedConfig.PoolSize, cfg.PoolSize)
			assert.Equal(t, tt.expectedConfig.PoolTimeout, cfg.PoolTimeout)
			assert.Equal(t, tt.expectedConfig.Password, cfg.Password)
		})
	}
}

func TestGetTracingConfig(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		serviceName string
		address     string
	}{
		{
			name:        "default tracing config",
			env:         map[string]string{},
			serviceName: "monitor",
			address:     "tempo:4317",
		},
		{
			name: "custom tracing config",
			env: map[string]string{
				"TRACE_SERVICE_NAME":    "api",
				"TRACE_SERVICE_ADDRESS": "localhost:4317",
			},
			serviceName: "api",
			address:     "localhost:4317",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := getTracingConfig()

			require.NotNil(t, cfg)
			assert.Equal(t, tt.serviceName, cfg.ServiceName)
			assert.Equal(t, tt.address, cfg.Address)
		})
	}
}
