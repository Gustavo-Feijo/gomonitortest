package redisinfra

import (
	"bytes"
	"errors"
	"gomonitor/internal/config"
	"log/slog"
	"testing"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	t.Parallel()
	client := New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	// Redis should never be nil, even if out, the client is returned and retries are handled with CB.
	assert.NotNil(t, client)
}

func TestNewRedisOut_ErrorIsLogged(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	logger := slog.New(
		slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}),
	)

	client := New(t.Context(), &config.RedisConfig{Addr: "localhost:0000"}, testCbCfg, logger)

	output := buf.String()

	assert.Contains(t, output, "couldn't connect correctly to Redis")
	assert.NotNil(t, client)
}

// Just to add the sweet 100% coverage.
func TestRedisTracing_ErrorIsLogged(t *testing.T) {
	t.Parallel()
	orig := instrumentRedisTracing
	t.Cleanup(func() { instrumentRedisTracing = orig })

	instrumentRedisTracing = func(
		redis.UniversalClient,
		...redisotel.TracingOption,
	) error {
		return errors.New("boom")
	}

	var buf bytes.Buffer
	logger := slog.New(
		slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}),
	)

	_ = New(t.Context(), testRedisCfg, testCbCfg, logger)

	output := buf.String()
	assert.Contains(t, output, "error while adding redis tracing")
}

func TestRedisCircuitBreaker(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	logger := slog.New(
		slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}),
	)
	client := New(t.Context(), testRedisCfg, testCbCfg, logger)

	client.Close()

	// Execute queries until reaching the CB limit
	for range testCbCfg.MaxFailures + 1 {
		client.Get(t.Context(), "test")
	}
	output := buf.String()

	assert.Contains(t, output, "circuit breaker changing state")
}

func TestRedisSet(t *testing.T) {
	t.Parallel()
	client := New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	// Happy path with Redis working.
	err := client.Set(t.Context(), "test", 5, 0)
	assert.Nil(t, err)

	client.Close()

	// Redis not working.
	err = client.Set(t.Context(), "test", 5, 0)
	assert.NotNil(t, err)
}

func TestRedisGet(t *testing.T) {
	t.Parallel()
	client := New(t.Context(), testRedisCfg, testCbCfg, slog.Default())

	val, err := client.Get(t.Context(), "nonexistent")
	assert.Equal(t, val, "")
	assert.ErrorIs(t, err, redis.Nil)

	_ = client.Set(t.Context(), "existent", 5, 0)

	val, err = client.Get(t.Context(), "existent")
	assert.Equal(t, val, "5")
	assert.Nil(t, err)
}
