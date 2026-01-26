package auth

import (
	"context"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *RefreshToken) error
	WithTx(tx *gorm.DB) RefreshTokenRepository
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (r *refreshTokenRepository) WithTx(tx *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: tx}
}

func (r *refreshTokenRepository) Create(ctx context.Context, refreshToken *RefreshToken) error {
	return r.db.WithContext(ctx).Create(refreshToken).Error
}
