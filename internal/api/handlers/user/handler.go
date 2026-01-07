package userhandler

import (
	"gomonitor/internal/domain/user"
	"gomonitor/internal/infra/deps"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger  *slog.Logger
	service *user.Service
}

func NewHandler(deps *deps.Deps) *Handler {
	svcDeps := &user.ServiceDeps{
		Logger: deps.Logger,
		DB:     deps.DB,
	}
	svc := user.NewService(svcDeps)

	return &Handler{
		logger:  deps.Logger,
		service: svc,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.create)
		users.GET("/:id", h.getByID)
	}
}
