package userhandler

import (
	userdto "gomonitor/internal/api/dto/user"
	pkgerrors "gomonitor/internal/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) create(c *gin.Context) {
	var req userdto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(pkgerrors.NewBadRequestError("Invalid JSON payload", err))
		return
	}

	input := req.ToDomainInput()

	user, err := h.service.CreateUser(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := userdto.ToCreateUserResponse(user)

	c.JSON(http.StatusCreated, resp)
}
