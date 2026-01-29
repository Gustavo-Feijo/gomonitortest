package ratelimit

import (
	"context"
	"errors"
	redisinfra "gomonitor/internal/infra/redis"
	"time"

	"github.com/google/uuid"
)

type redisLimiter struct {
	redisClient redisinfra.RedisClient
	BaseLimiterConfig
}

func NewRedisLimiter(redisClient redisinfra.RedisClient, opts ...BaseOption) RateLimiter {
	rl := &redisLimiter{
		redisClient: redisClient,
		BaseLimiterConfig: BaseLimiterConfig{
			limit:     100,
			keyPrefix: "rate_limit",
			window:    time.Minute,
		},
	}

	for _, opt := range opts {
		opt(&rl.BaseLimiterConfig)
	}

	return rl
}

func (r *redisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := r.keyPrefix + ":" + key
	now := time.Now().UnixMilli()
	windowMillis := r.window.Milliseconds()
	limit := r.limit

	script := `
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local limit = tonumber(ARGV[3])
	local uuid = ARGV[4]
	
	redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
	
	local current = redis.call('ZCARD', key)
	
	if current >= limit then
		return 0
	end
	
	redis.call('ZADD', key, now, uuid)
	redis.call('EXPIRE', key, math.ceil(window / 1000) + 1)
	
	return 1
	`

	result, err := r.redisClient.Eval(ctx, script, []string{redisKey}, now, windowMillis, limit, uuid.New().String())
	if err != nil {
		return false, err
	}

	val, ok := result.(int64)
	if !ok {
		return false, errors.New("couldn't cast eval to int64")
	}

	return val == 1, nil
}
