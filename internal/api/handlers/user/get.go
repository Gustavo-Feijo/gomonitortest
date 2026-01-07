package userhandler

import (
	userdto "gomonitor/internal/api/dto/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getByID(c *gin.Context) {
	var req userdto.GetUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	input := req.ToDomainInput()

	user, err := h.service.GetUser(c.Request.Context(), input)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	resp := userdto.ToGetUserResponse(user)

	c.JSON(http.StatusOK, resp)
}
