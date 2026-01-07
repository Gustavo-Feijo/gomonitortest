package authhandler

import (
	authdto "gomonitor/internal/api/dto/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Refresh(c *gin.Context) {
	var req authdto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	input := req.ToDomainInput()

	refresh, err := h.service.Refresh(c.Request.Context(), input)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	resp := authdto.ToRefreshResponse(refresh)

	c.JSON(http.StatusOK, resp)
}
