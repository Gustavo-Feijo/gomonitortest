package auth_test

import (
	"gomonitor/internal/config"
	"gomonitor/internal/domain/auth"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/pkg/jwt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	deps := &auth.ServiceDeps{
		UserRepo:     user.NewRepository(&gorm.DB{}),
		Logger:       &slog.Logger{},
		TokenManager: jwt.NewTokenManager(&config.AuthConfig{}),
		AuthConfig:   &config.AuthConfig{},
	}

	service := auth.NewService(deps)
	assert.NotNil(t, service)
}
