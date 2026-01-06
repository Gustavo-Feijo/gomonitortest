package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, user *User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}
