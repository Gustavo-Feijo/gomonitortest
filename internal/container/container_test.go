package container_test

import (
	"gomonitor/internal/config"
	"gomonitor/internal/container"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/mocks"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewContainer(t *testing.T) {
	deps := &deps.Deps{
		DB:           &gorm.DB{},
		Hasher:       &mocks.MockPasswordHasher{},
		Logger:       slog.Default(),
		Redis:        &mocks.MockRedisClient{},
		TokenManager: &mocks.MockJwtManager{},
	}
	container := container.New(deps, &config.Config{
		RateLimit: &config.RateLimitConfig{
			IPLimit:    10,
			IPWindow:   time.Minute,
			UserLimit:  5,
			UserWindow: time.Minute,
		},
	})
	require.NotNil(t, container)
}
