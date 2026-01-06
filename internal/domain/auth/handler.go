package auth

import (
	"gomonitor/internal/config"
	"gomonitor/internal/infra/deps"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	logger  *slog.Logger
	service *service
}

func NewHandler(deps *deps.Deps, authCfg *config.AuthConfig) *handler {
	svcDeps := &ServiceDeps{
		AuthConfig: authCfg,
		DB:         deps.DB,
		Logger:     deps.Logger,
	}
	svc := NewService(svcDeps)

	return &handler{
		logger:  deps.Logger,
		service: svc,
	}
}

func (h *handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("login", h.Login)
		auth.POST("refresh", h.Refresh)
	}
}

func (h *handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	login, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, login)
}

func (h *handler) Refresh(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	login, err := h.service.Refresh(c.Request.Context(), req)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, login)
}
