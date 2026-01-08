package userhandler

import (
	userdto "gomonitor/internal/api/dto/user"
	pkgerrors "gomonitor/internal/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getByID(c *gin.Context) {
	var req userdto.GetUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(pkgerrors.NewBadRequestError("Invalid ID parameter", err))
		return
	}

	input := req.ToDomainInput()

	user, err := h.service.GetUser(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := userdto.ToGetUserResponse(user)

	c.JSON(http.StatusOK, resp)
}
