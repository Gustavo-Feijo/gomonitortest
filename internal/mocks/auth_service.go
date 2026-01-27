package mocks

import (
	"context"
	"gomonitor/internal/domain/auth"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, input auth.LoginInput) (*auth.LoginOutput, error) {
	args := m.Called(ctx, input)
	var lo *auth.LoginOutput
	if args.Get(0) != nil {
		lo = args.Get(0).(*auth.LoginOutput)
	}
	return lo, args.Error(1)
}

func (m *MockAuthService) Refresh(ctx context.Context, input auth.RefreshInput) (*auth.RefreshOutput, error) {
	args := m.Called(ctx, input)
	var ro *auth.RefreshOutput
	if args.Get(0) != nil {
		ro = args.Get(0).(*auth.RefreshOutput)
	}
	return ro, args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
