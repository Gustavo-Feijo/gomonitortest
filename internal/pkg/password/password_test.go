package password_test

import (
	"gomonitor/internal/pkg/password"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "invalid password",
			password:    strings.Repeat("a", 80),
			expectedErr: bcrypt.ErrPasswordTooLong,
		},
		{
			name:     "success",
			password: "valid-password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := password.NewPasswordHasher(bcrypt.DefaultCost)

			hash, err := pm.HashPassword(tt.password)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, hash)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		hashGen     func() string
		expectedErr error
	}{
		{
			name:     "mismatched hash and password",
			password: "test-password",
			hashGen: func() string {
				pm := password.NewPasswordHasher(bcrypt.DefaultCost)
				hash, _ := pm.HashPassword("different-test-password")
				return hash
			},
			expectedErr: bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:     "success",
			password: "test-password",
			hashGen: func() string {
				pm := password.NewPasswordHasher(bcrypt.DefaultCost)
				hash, _ := pm.HashPassword("test-password")
				return hash
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := password.NewPasswordHasher(bcrypt.DefaultCost)

			err := pm.VerifyPassword(tt.hashGen(), tt.password)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
