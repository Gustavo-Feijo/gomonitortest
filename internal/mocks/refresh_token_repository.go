package mocks

import (
	"context"
	"gomonitor/internal/domain/auth"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, refreshToken *auth.RefreshToken) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) WithTx(tx *gorm.DB) auth.RefreshTokenRepository {
	return m
}
