package mocks

import (
	"context"
	"gomonitor/internal/domain/user"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	args := m.Called(ctx, id)

	var u *user.User
	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}

	return u, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)

	var u *user.User
	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}

	return u, args.Error(1)
}
