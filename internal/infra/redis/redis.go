package redisinfra

import (
	"context"
	"gomonitor/internal/config"
	"log/slog"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/logging"
	"github.com/sony/gobreaker/v2"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
}

// Wrapper on the redis client, handle nil client checks.
type redisClient struct {
	cb     *gobreaker.CircuitBreaker[any]
	client *redis.Client
}

var instrumentRedisTracing = redisotel.InstrumentTracing

// New creates and returns a new redis connection pool.
// Always returns the client, even if can't connect to it.
// Implements circuit breaking and retries in the case Redis is down.
func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) RedisClient {
	redis.SetLogger(&logging.VoidLogger{})
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.Database,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolTimeout:  cfg.Redis.PoolTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("couldn't connect correctly to Redis", slog.Any("err", err))
	}

	if err := instrumentRedisTracing(
		client,
		redisotel.WithCommandFilter(redisotel.DefaultCommandFilter),
		redisotel.WithDialFilter(true),
	); err != nil {
		logger.Warn("error while adding redis tracing", slog.Any("err", err))
	}

	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
		Name:        "redis",
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    time.Minute,
		Timeout:     time.Minute,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.CircuitBreaker.MaxFailures
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Warn("circuit breaker changing state",
				slog.String("breaker", name),
				slog.String("from", from.String()),
				slog.String("to", to.String()),
			)
		},
		IsSuccessful: func(err error) bool {
			// Redis nil is success, just cache miss.
			return err == nil || err == redis.Nil
		},
	})

	return &redisClient{
		client: client,
		cb:     cb,
	}
}

// Get wrapper with Circuit break and result unwrapping.
// If result is redis.Nil, the error is returned normally.
func (rs *redisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := rs.cb.Execute(func() (any, error) {
		return rs.client.Get(ctx, key).Result()
	})

	if err != nil {
		return "", err
	}

	return res.(string), nil
}

// Set wrapper with Circuit break and result unwrapping.
func (rs *redisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	_, err := rs.cb.Execute(func() (any, error) {
		return nil, rs.client.Set(ctx, key, value, ttl).Err()
	})
	return err
}
