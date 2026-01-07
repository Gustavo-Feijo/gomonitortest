package authhandler

import (
	authdto "gomonitor/internal/api/dto/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var req authdto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	input := req.ToDomainInput()

	login, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	resp := authdto.ToLoginResponse(login)

	c.JSON(http.StatusOK, resp)
}
