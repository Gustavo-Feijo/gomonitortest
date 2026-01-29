package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}
