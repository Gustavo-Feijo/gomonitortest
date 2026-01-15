package mocks

import (
	"context"
	"gomonitor/internal/domain/user"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, input user.CreateUserInput) (*user.User, error) {
	args := m.Called(ctx, input)
	var u *user.User
	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}
	return u, args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, input user.GetUserInput) (*user.User, error) {
	args := m.Called(ctx, input)
	var u *user.User
	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}
	return u, args.Error(1)
}
