package auth

import (
	"context"
	"gomonitor/internal/domain/user"

	"gorm.io/gorm"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	GetUserByID(ctx context.Context, id uint) (*user.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var usr user.User
	err := r.db.Model(&user.User{}).
		Where("email = ?", email).
		First(&usr).Error

	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *repository) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	var usr user.User
	err := r.db.Model(&user.User{}).
		Where("id = ?", id).
		First(&usr).Error

	if err != nil {
		return nil, err
	}

	return &usr, nil
}
