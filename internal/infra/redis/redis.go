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
	Raw() *redis.Client
	Get(ctx context.Context, key string) (string, error)
}

// Wrapper on the redis client, handle nil client checks.
type redisClient struct {
	cb     *gobreaker.CircuitBreaker[any]
	client *redis.Client
}

// New creates and returns a new redis connection pool.
func New(ctx context.Context, cfg *config.RedisConfig, logger *slog.Logger) (RedisClient, error) {
	redis.SetLogger(&logging.VoidLogger{})
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.Database,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		PoolTimeout:  cfg.PoolTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("couldn't connect correctly to Redis", slog.Any("err", err))
	}

	if err := redisotel.InstrumentTracing(
		client,
		redisotel.WithCommandFilter(redisotel.DefaultCommandFilter),
		redisotel.WithDialFilter(true),
	); err != nil {
		logger.Warn("error while adding redis tracing", slog.Any("err", err))
	}

	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
		Name:        "redis",
		MaxRequests: 5,
		Interval:    time.Minute,
		Timeout:     time.Minute,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
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
	}, nil
}

// Raw is a method to retrieve the raw client directly for one off operations.
// It's preferable to implement the wrapper to the necessary method.
func (rs *redisClient) Raw() *redis.Client {
	return rs.client
}

// Get wrapper with nil client check and result extracted.
func (rs *redisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := rs.cb.Execute(func() (any, error) {
		return rs.client.Get(ctx, key).Result()
	})

	if err != nil {
		return "", err
	}

	return res.(string), err
}
