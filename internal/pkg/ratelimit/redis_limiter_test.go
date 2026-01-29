package ratelimit_test

import (
	redisinfra "gomonitor/internal/infra/redis"
	"gomonitor/internal/pkg/ratelimit"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedisLimiter(t *testing.T) {
	t.Parallel()
	client := redisinfra.New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	rl := ratelimit.NewRedisLimiter(client,
		ratelimit.WithLimit(5),
		ratelimit.WithPrefix("rltest"),
		ratelimit.WithWindow(5*time.Second),
	)

	require.NotNil(t, rl)
}

func TestRedisLimiterAllow(t *testing.T) {
	t.Parallel()
	client := redisinfra.New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	testLimit := 5
	window := time.Second * 2
	key := t.Name()
	rl := ratelimit.NewRedisLimiter(client,
		ratelimit.WithLimit(testLimit),
		ratelimit.WithPrefix("rltest"),
		ratelimit.WithWindow(window),
	)

	require.NotNil(t, rl)

	for range testLimit {
		allowed, err := rl.Allow(t.Context(), key)
		require.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rl.Allow(t.Context(), key)
	require.NoError(t, err)
	assert.False(t, allowed)

	// Wait for window reset.
	time.Sleep(window + 100*time.Millisecond)

	allowed, err = rl.Allow(t.Context(), key)
	require.NoError(t, err)
	assert.True(t, allowed)
}
