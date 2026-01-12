package authhandler

import (
	"gomonitor/internal/domain/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger  *slog.Logger
	service auth.Service
}

func NewHandler(logger *slog.Logger, svc auth.Service) *Handler {
	return &Handler{
		logger:  logger,
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
