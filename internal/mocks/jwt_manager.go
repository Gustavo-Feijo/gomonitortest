package mocks

import (
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/pkg/jwt"

	"github.com/stretchr/testify/mock"
)

type MockJwtManager struct {
	mock.Mock
}

func (m *MockJwtManager) GenerateRefreshToken(userID uint, role identity.UserRole) (*jwt.RefreshTokenResult, error) {
	args := m.Called(userID, role)
	var j *jwt.RefreshTokenResult
	if args.Get(0) != nil {
		j = args.Get(0).(*jwt.RefreshTokenResult)
	}
	return j, args.Error(1)
}

func (m *MockJwtManager) GenerateAccessToken(userID uint, role identity.UserRole) (*jwt.AccessTokenResult, error) {
	args := m.Called(userID, role)
	var j *jwt.AccessTokenResult
	if args.Get(0) != nil {
		j = args.Get(0).(*jwt.AccessTokenResult)
	}
	return j, args.Error(1)
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
