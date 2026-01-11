package userhandler

import (
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      user.Service
	tokenManager jwt.TokenManager
}

func NewHandler(deps *deps.Deps) *Handler {
	userRepo := user.NewRepository(deps.DB)
	svcDeps := &user.ServiceDeps{
		Hasher:   deps.Hasher,
		Logger:   deps.Logger,
		UserRepo: userRepo,
	}
	svc := user.NewService(svcDeps)

	return &Handler{
		service:      svc,
		tokenManager: deps.TokenManager,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users", middlewares.AuthMiddleware(h.tokenManager))
	{
		users.POST("", h.create)
		users.GET("/:id", h.getByID)
	}
}
