package container_test

import (
	"gomonitor/internal/config"
	"gomonitor/internal/container"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/mocks"
	"log/slog"
	"testing"

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
	container := container.New(deps, &config.Config{})
	require.NotNil(t, container)
}
