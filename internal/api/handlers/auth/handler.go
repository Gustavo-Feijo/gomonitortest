package authhandler

import (
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/domain/auth"
	"gomonitor/internal/pkg/jwt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger       *slog.Logger
	service      auth.Service
	tokenManager jwt.TokenManager
}

func NewHandler(logger *slog.Logger, svc auth.Service, tokenManager jwt.TokenManager) *Handler {
	return &Handler{
		logger:       logger,
		service:      svc,
		tokenManager: tokenManager,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("login", h.Login)
		auth.POST("refresh", h.Refresh)

		logout := auth.Group("logout", middlewares.AuthMiddleware(h.tokenManager))
		{
			logout.POST("", h.Logout)
			logout.POST("all", h.LogoutAll)
		}
	}
}
