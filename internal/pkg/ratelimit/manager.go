package ratelimit

import (
	"context"
	"errors"
)

type limitManager struct {
	limiter  RateLimiter
	fallback RateLimiter
}

type LimitManagerOptions func(*limitManager)

func New(opts ...LimitManagerOptions) RateLimiter {
	lm := &limitManager{
		limiter: NewMemoryLimiter(),
	}

	for _, opt := range opts {
		opt(lm)
	}

	return lm
}

func (lm *limitManager) Allow(ctx context.Context, key string) (bool, error) {
	res, err := lm.limiter.Allow(ctx, key)
	if err == nil {
		return res, nil
	}

	if lm.fallback == nil {
		return false, errors.New("main limiter out, no fallback set")
	}

	return lm.fallback.Allow(ctx, key)
}

func WithLimiter(limiter RateLimiter) LimitManagerOptions {
	return func(lm *limitManager) {
		lm.limiter = limiter
	}
}

func WithFallback(fallback RateLimiter) LimitManagerOptions {
	return func(lm *limitManager) {
		lm.fallback = fallback
	}
}
