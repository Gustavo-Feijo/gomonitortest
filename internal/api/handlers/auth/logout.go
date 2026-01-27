package authhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Logout(c *gin.Context) {
	if err := h.service.Logout(c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) LogoutAll(c *gin.Context) {
	if err := h.service.LogoutAll(c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
