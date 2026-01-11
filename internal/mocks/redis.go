package mocks

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Raw() *redis.Client {
	args := m.Called()
	return args.Get(0).(*redis.Client)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(string), args.Error(1)
}
