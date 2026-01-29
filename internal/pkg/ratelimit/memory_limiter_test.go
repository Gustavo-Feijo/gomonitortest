package ratelimit_test

import (
	"gomonitor/internal/pkg/ratelimit"
	"gomonitor/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryLimiter(t *testing.T) {
	t.Parallel()

	rl := ratelimit.NewMemoryLimiter(
		ratelimit.WithLimit(5),
		ratelimit.WithPrefix("rltest"),
		ratelimit.WithWindow(5*time.Second),
	)

	require.NotNil(t, rl)
}

func TestMemoryLimiterAllow(t *testing.T) {
	t.Parallel()

	testLimit := 150
	window := time.Second * 2
	key := t.Name()
	rl := ratelimit.NewMemoryLimiter(
		ratelimit.WithLimit(testLimit),
		ratelimit.WithPrefix("rltest"),
		ratelimit.WithWindow(window),
	)

	require.NotNil(t, rl)

	for range testLimit {
		allowed, err := rl.Allow(t.Context(), key)
		require.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rl.Allow(t.Context(), key)
	require.NoError(t, err)
	assert.False(t, allowed)

	// Wait for window reset.
	time.Sleep(window + 100*time.Millisecond)

	for range testLimit {
		allowed, err := rl.Allow(t.Context(), key)
		require.NoError(t, err)
		assert.True(t, allowed)
	}
}

func TestMemoryLimiterAllowDoneCtx(t *testing.T) {
	t.Parallel()

	testLimit := 5
	window := time.Second * 2
	key := t.Name()
	rl := ratelimit.NewMemoryLimiter(
		ratelimit.WithLimit(testLimit),
		ratelimit.WithPrefix("rltest"),
		ratelimit.WithWindow(window),
	)
	cancelledCtx := testutil.GetCancelledCtx(t.Context())

	allowed, err := rl.Allow(cancelledCtx, key)
	require.Error(t, err)
	assert.False(t, allowed)
}
