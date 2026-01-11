package mocks

import (
	"gomonitor/internal/pkg/identity"

	"github.com/stretchr/testify/mock"
)

type MockJwtManager struct {
	mock.Mock
}

func (m *MockJwtManager) GenerateRefreshToken(userID uint, role identity.UserRole) (string, error) {
	args := m.Called(userID, role)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockJwtManager) GenerateAccessToken(userID uint, role identity.UserRole) (string, error) {
	args := m.Called(userID, role)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockJwtManager) ValidateRefreshToken(tokenString string) (*identity.Principal, error) {
	args := m.Called(tokenString)

	var i *identity.Principal
	if args.Get(0) != nil {
		i = args.Get(0).(*identity.Principal)
	}

	return i, args.Error(1)
}

func (m *MockJwtManager) ValidateAccessToken(tokenString string) (*identity.Principal, error) {
	args := m.Called(tokenString)

	var i *identity.Principal
	if args.Get(0) != nil {
		i = args.Get(0).(*identity.Principal)
	}

	return i, args.Error(1)
}
