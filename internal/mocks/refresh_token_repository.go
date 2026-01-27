package mocks

import (
	"context"
	"gomonitor/internal/domain/auth"

	"github.com/google/uuid"
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

func (m *MockRefreshTokenRepository) GetByJTI(ctx context.Context, jti uuid.UUID) (*auth.RefreshToken, error) {
	args := m.Called(ctx, jti)

	var t *auth.RefreshToken
	if args.Get(0) != nil {
		t = args.Get(0).(*auth.RefreshToken)
	}

	return t, args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeByJTI(ctx context.Context, jti uuid.UUID) error {
	args := m.Called(ctx, jti)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByUserID(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) WithTx(tx *gorm.DB) auth.RefreshTokenRepository {
	return m
}
