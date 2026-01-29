package ratelimit

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

type BaseLimiterConfig struct {
	limit     int
	keyPrefix string
	window    time.Duration
}

type BaseOption func(*BaseLimiterConfig)

func WithLimit(limit int) BaseOption {
	return func(lc *BaseLimiterConfig) {
		lc.limit = limit
	}
}

func WithPrefix(prefix string) BaseOption {
	return func(lc *BaseLimiterConfig) {
		lc.keyPrefix = prefix
	}
}

func WithWindow(window time.Duration) BaseOption {
	return func(lc *BaseLimiterConfig) {
		lc.window = window
	}
}
