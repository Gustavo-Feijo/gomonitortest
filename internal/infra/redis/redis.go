package redisinfra

import (
	"context"
	"fmt"
	"gomonitor/internal/config"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// Wrapper on the redis client, handle nil client checks.
type RedisClient struct {
	client *redis.Client
}

// New creates and returns a new redis connection pool.
func New(ctx context.Context, cfg *config.RedisConfig) (*RedisClient, error) {
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
		return &RedisClient{nil}, err
	}

	if err := redisotel.InstrumentTracing(
		client,
		redisotel.WithCommandFilter(redisotel.DefaultCommandFilter),
		redisotel.WithDialFilter(true),
	); err != nil {
		return &RedisClient{nil}, fmt.Errorf("error while adding redis tracing: %v", err)
	}

	return &RedisClient{client}, nil
}

// Raw is a method to retrieve the raw client directly for one off operations.
// It's preferable to implement the wrapper to the necessary method.
func (rs *RedisClient) Raw() *redis.Client {
	return rs.client
}

// Get wrapper with nil client check and result extracted.
func (rs *RedisClient) Get(ctx context.Context, key string) (string, error) {
	if rs.client == nil {
		return "", redis.Nil
	}
	return rs.client.Get(ctx, key).Result()
}
