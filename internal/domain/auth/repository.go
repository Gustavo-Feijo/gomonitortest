package auth

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *RefreshToken) error
	GetByJTI(ctx context.Context, jti uuid.UUID) (*RefreshToken, error)
	RevokeByJTI(ctx context.Context, jti uuid.UUID) error
	RevokeByUserID(ctx context.Context, id uint) error
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

func (r *refreshTokenRepository) GetByJTI(ctx context.Context, jti uuid.UUID) (*RefreshToken, error) {
	var refreshToken RefreshToken
	if err := r.db.WithContext(ctx).First(&refreshToken, jti).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *refreshTokenRepository) RevokeByJTI(ctx context.Context, jti uuid.UUID) error {
	return r.db.
		WithContext(ctx).
		Model(&RefreshToken{}).
		Where("jti = ? AND revoked_at IS NULL", jti).
		Update("revoked_at", gorm.Expr("NOW()")).
		Error
}

func (r *refreshTokenRepository) RevokeByUserID(ctx context.Context, id uint) error {
	return r.db.
		WithContext(ctx).
		Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", gorm.Expr("NOW()")).
		Error
}
