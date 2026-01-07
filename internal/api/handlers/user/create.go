package userhandler

import (
	userdto "gomonitor/internal/api/dto/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) create(c *gin.Context) {
	var req userdto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	input := req.ToDomainInput()

	user, err := h.service.CreateUser(c.Request.Context(), input)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	resp := userdto.ToCreateUserResponse(user)

	c.JSON(http.StatusCreated, resp)
}
