package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("ENVIRONMENT", "test")

	// Required Admin config
	t.Setenv("ADMIN_EMAIL", "admin@test.com")
	t.Setenv("ADMIN_PASSWORD", "password")

	// Required Auth config
	t.Setenv("AUTH_ACCESS_TOKEN_SECRET", "access")
	t.Setenv("AUTH_REFRESH_TOKEN_SECRET", "refresh")
	t.Setenv("AUTH_FAKE_HASH", "fakehash")
	t.Setenv("AUTH_ACCESS_TOKEN_TTL", "1h")
	t.Setenv("AUTH_REFRESH_TOKEN_TTL", "168h")

	cfg, err := Load()

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.NotNil(t, cfg.Admin)
	assert.NotNil(t, cfg.Auth)
	assert.NotNil(t, cfg.CircuitBreaker)
	assert.NotNil(t, cfg.Database)
	assert.NotNil(t, cfg.HTTP)
	assert.NotNil(t, cfg.Logging)
	assert.NotNil(t, cfg.Redis)
	assert.NotNil(t, cfg.Tracing)
}

func TestLoad_Error(t *testing.T) {
	// Set env as production to simulate the behaviour of missing envs.
	// If set to anything different, will get all necessary values from the .env files, breaking the test idea.
	t.Setenv("ENVIRONMENT", "production")

	// Missing ADMIN_EMAIL
	t.Setenv("ADMIN_PASSWORD", "password")

	t.Setenv("AUTH_ACCESS_TOKEN_SECRET", "access")
	t.Setenv("AUTH_REFRESH_TOKEN_SECRET", "refresh")
	t.Setenv("AUTH_FAKE_HASH", "fakehash")
	t.Setenv("AUTH_ACCESS_TOKEN_TTL", "1h")
	t.Setenv("AUTH_REFRESH_TOKEN_TTL", "168h")

	cfg, err := Load()

	require.Error(t, err)
	assert.Nil(t, cfg)
}
