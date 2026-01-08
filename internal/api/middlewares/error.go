package middlewares

import (
	"gomonitor/internal/observability/logging"
	pkgerrors "gomonitor/internal/pkg/errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		appErr, ok := err.(*pkgerrors.AppError)
		if !ok {
			logger := logging.FromContext(c.Request.Context())
			logger.Error("unexpected error occurred", slog.Any("err", err))

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "An unexpected error occurred",
			})
			return
		}

		if appErr.StatusCode >= 500 {
			logger := logging.FromContext(c.Request.Context())
			logger.Error("unexpected error occurred", slog.Any("err", err))
		}

		c.JSON(appErr.StatusCode, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
	}
}
