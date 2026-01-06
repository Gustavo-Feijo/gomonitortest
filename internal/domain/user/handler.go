package user

import (
	"gomonitor/internal/infra/deps"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	logger  *slog.Logger
	service *service
}

func NewHandler(deps *deps.Deps) *handler {
	svcDeps := &ServiceDeps{
		Logger: deps.Logger,
		DB:     deps.DB,
	}
	svc := NewService(svcDeps)

	return &handler{
		logger:  deps.Logger,
		service: svc,
	}
}

func (h *handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.create)
		users.GET("/:id", h.getByID)
	}
}

func (h *handler) create(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *handler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, user)
}
