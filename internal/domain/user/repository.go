package user

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	WithTx(tx *gorm.DB) UserRepository
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{db: tx}
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var usr User
	err := r.db.Model(&User{}).
		Where("email = ?", email).
		First(&usr).Error

	if err != nil {
		return nil, err
	}

	return &usr, nil
}
