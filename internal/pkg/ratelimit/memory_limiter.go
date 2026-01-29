package ratelimit

import (
	"context"
	"sync"
	"time"
)

type memoryLimiter struct {
	mu             sync.Mutex
	cleanupCounter int
	BaseLimiterConfig
	entries map[string]*memoryEntry
}

type memoryEntry struct {
	count   int
	resetAt time.Time
}

func NewMemoryLimiter(opts ...BaseOption) RateLimiter {
	ml := &memoryLimiter{
		BaseLimiterConfig: BaseLimiterConfig{
			limit:     100,
			keyPrefix: "rate_limit",
			window:    time.Minute,
		},
		entries: make(map[string]*memoryEntry),
	}

	for _, opt := range opts {
		opt(&ml.BaseLimiterConfig)
	}

	return ml
}

func (m *memoryLimiter) Allow(ctx context.Context, key string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		now := time.Now()
		memKey := m.keyPrefix + ":" + key

		m.mu.Lock()
		defer m.mu.Unlock()

		m.cleanupCounter++
		if m.cleanupCounter%100 == 0 {
			m.cleanupExpired(now)
		}

		entry, ok := m.entries[memKey]
		if !ok || now.After(entry.resetAt) {
			m.entries[memKey] = &memoryEntry{
				count:   1,
				resetAt: now.Add(m.window),
			}
			return true, nil
		}

		if entry.count >= m.limit {
			return false, nil
		}

		entry.count++
		return true, nil
	}
}

func (m *memoryLimiter) cleanupExpired(now time.Time) {
	for k, v := range m.entries {
		if now.After(v.resetAt) {
			delete(m.entries, k)
		}
	}
}
