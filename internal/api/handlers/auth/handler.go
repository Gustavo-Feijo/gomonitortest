package authhandler

import (
	"gomonitor/internal/config"
	"gomonitor/internal/domain/auth"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/infra/deps"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger  *slog.Logger
	service auth.Service
}

func NewHandler(deps *deps.Deps, authCfg *config.AuthConfig) *Handler {
	userRepo := user.NewRepository(deps.DB)
	svcDeps := &auth.ServiceDeps{
		AuthConfig:   authCfg,
		Hasher:       deps.Hasher,
		UserRepo:     userRepo,
		Logger:       deps.Logger,
		TokenManager: deps.TokenManager,
	}
	svc := auth.NewService(svcDeps)

	return &Handler{
		logger:  deps.Logger,
		service: svc,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("login", h.Login)
		auth.POST("refresh", h.Refresh)
	}
}
