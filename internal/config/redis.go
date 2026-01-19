package config

import (
	"fmt"
	"time"
)

// Redis configuration.
type RedisConfig struct {
	Addr         string
	Database     int
	MaxRetries   int
	MinIdleConns int
	Password     string
	PoolSize     int
	PoolTimeout  time.Duration
}

// getRedisConfig loads the redis environments and return the full config.
func getRedisConfig() *RedisConfig {
	addr := getEnv("REDIS_ADDR", "")

	if addr == "" {
		host := getEnv("REDIS_HOST", "redis")
		port := getEnv("REDIS_PORT", "6379")
		addr = fmt.Sprintf("%s:%s", host, port)
	}

	return &RedisConfig{
		Addr:         addr,
		Password:     getEnv("REDIS_PASSWORD", ""),
		Database:     getIntEnv("REDIS_DATABASE", 0),
		MaxRetries:   getIntEnv("REDIS_MAX_RETRIES", 3),
		MinIdleConns: getIntEnv("REDIS_MIN_IDLE_CONNS", 10),
		PoolSize:     getIntEnv("REDIS_POOL_SIZE", 100),
		PoolTimeout:  time.Duration(getIntEnv("REDIS_POOL_TIMEOUT_SECONDS", 30)) * time.Second,
	}
}
