package authhandler

import (
	authdto "gomonitor/internal/api/dto/auth"
	"gomonitor/internal/observability/logging"
	pkgerrors "gomonitor/internal/pkg/errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var req authdto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(pkgerrors.NewBadRequestError("Invalid JSON payload", err))
		return
	}

	input := req.ToDomainInput()

	login, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := authdto.ToLoginResponse(login)

	logging.FromContext(c.Request.Context()).Info("successfull login attempt", slog.String("user", input.Email))

	c.JSON(http.StatusOK, resp)
}
