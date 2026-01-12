package userhandler

import (
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/pkg/jwt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger       *slog.Logger
	service      user.Service
	tokenManager jwt.TokenManager
}

func NewHandler(logger *slog.Logger, svc user.Service, tokenManager jwt.TokenManager) *Handler {
	return &Handler{
		logger:       logger,
		service:      svc,
		tokenManager: tokenManager,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users", middlewares.AuthMiddleware(h.tokenManager))
	{
		users.POST("", h.create)
		users.GET("/:id", h.getByID)
	}
}
