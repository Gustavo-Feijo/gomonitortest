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

func TestNewManager(t *testing.T) {
	t.Parallel()
	client := redisinfra.New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	rl := ratelimit.New(
		ratelimit.WithLimiter(
			ratelimit.NewRedisLimiter(
				client,
				ratelimit.WithLimit(10),
				ratelimit.WithPrefix("rate_limit"),
				ratelimit.WithWindow(time.Minute),
			),
		),
		ratelimit.WithFallback(
			ratelimit.NewMemoryLimiter(
				ratelimit.WithLimit(5),
				ratelimit.WithPrefix("rate_limit"),
				ratelimit.WithWindow(time.Minute),
			),
		),
	)

	require.NotNil(t, rl)
}

func TestNewManagerAllowWithFallback(t *testing.T) {
	client := redisinfra.New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	redisTestLimit := 5
	inMemoryTestLimit := 2
	window := time.Second * 2
	key := t.Name()

	rl := ratelimit.New(
		ratelimit.WithLimiter(
			ratelimit.NewRedisLimiter(
				client,
				ratelimit.WithLimit(redisTestLimit),
				ratelimit.WithPrefix("rate_limit"),
				ratelimit.WithWindow(window),
			),
		),
		ratelimit.WithFallback(
			ratelimit.NewMemoryLimiter(
				ratelimit.WithLimit(inMemoryTestLimit),
				ratelimit.WithPrefix("rate_limit"),
				ratelimit.WithWindow(window),
			),
		),
	)

	require.NotNil(t, rl)

	// First trying the main limiter behavior
	for range redisTestLimit {
		allowed, err := rl.Allow(t.Context(), key)
		require.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rl.Allow(t.Context(), key)
	require.NoError(t, err)
	assert.False(t, allowed)

	_ = client.Close()

	for range inMemoryTestLimit {
		allowed, err := rl.Allow(t.Context(), key)
		require.NoError(t, err)
		assert.True(t, allowed)
	}

	memAllowed, err := rl.Allow(t.Context(), key)
	require.NoError(t, err)
	assert.False(t, memAllowed)
}

func TestNewManagerAllowWithoutFallback(t *testing.T) {
	client := redisinfra.New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	redisTestLimit := 5
	window := time.Second * 2
	key := t.Name()

	rl := ratelimit.New(
		ratelimit.WithLimiter(
			ratelimit.NewRedisLimiter(
				client,
				ratelimit.WithLimit(redisTestLimit),
				ratelimit.WithPrefix("rate_limit"),
				ratelimit.WithWindow(window),
			),
		),
	)

	require.NotNil(t, rl)

	// First trying the main limiter behavior
	for range redisTestLimit {
		allowed, err := rl.Allow(t.Context(), key)
		require.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rl.Allow(t.Context(), key)
	require.NoError(t, err)
	assert.False(t, allowed)

	_ = client.Close()

	noFallbackAllow, err := rl.Allow(t.Context(), key)
	require.Error(t, err)
	assert.False(t, noFallbackAllow)
}
