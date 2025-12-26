package redisinfra

import (
	"gomonitor/internal/config"
	"log"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// New creates and returns a new redis connection pool.
func New(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.Database,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		PoolTimeout:  cfg.PoolTimeout,
	})

	if err := redisotel.InstrumentTracing(client); err != nil {
		log.Printf("Error while adding redis tracing: %v", err)
	}

	if err := redisotel.InstrumentMetrics(client); err != nil {
		log.Printf("Error while adding redis metrics: %v", err)
	}

	return client
}
