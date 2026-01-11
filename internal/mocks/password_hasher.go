package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockPasswordHasher) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}
