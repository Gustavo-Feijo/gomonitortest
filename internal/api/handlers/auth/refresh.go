package authhandler

import (
	authdto "gomonitor/internal/api/dto/auth"
	pkgerrors "gomonitor/internal/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Refresh(c *gin.Context) {
	var req authdto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(pkgerrors.NewBadRequestError("Invalid JSON payload", err))
		return
	}

	input := req.ToDomainInput()

	refresh, err := h.service.Refresh(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := authdto.ToRefreshResponse(refresh)

	c.JSON(http.StatusOK, resp)
}
