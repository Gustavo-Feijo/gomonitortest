package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// startRedis creates a new redis test container.
// Returns the container, the address, the cleanup function and any error.
func StartRedis(ctx context.Context) (*tcredis.RedisContainer, string, func(ctx context.Context) error, error) {
	container, err := tcredis.Run(ctx, "redis:8.4-alpine")
	if err != nil {
		return nil, "", nil, err
	}

	addr, err := container.Endpoint(ctx, "")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", nil, err
	}

	return container, addr, func(ctx context.Context) error { return container.Terminate(ctx) }, nil
}

// StartTestRedis creates a new redis test container, set up cleanup and environment variables.
func StartTestRedis(t *testing.T) *tcredis.RedisContainer {
	ctx := context.Background()

	container, addr, cleanup, err := StartRedis(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = cleanup(ctx)
	})

	setupTestRedisEnv(t, addr)

	return container
}

func setupTestRedisEnv(t *testing.T, addr string) {
	t.Setenv("REDIS_ADDR", addr)
}
