package middlewares

import (
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenManager *jwt.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			_ = c.Error(pkgerrors.NewUnauthorizedError("Authorization header required"))
			c.Abort()
			return
		}

		principal, err := tokenManager.ValidateAccessToken(token)
		if err != nil {
			_ = c.Error(pkgerrors.NewUnauthorizedError("Invalid or expired token", err))
			c.Abort()
			return
		}

		authenticatedContext := identity.WithPrincipal(c.Request.Context(), principal)
		c.Request = c.Request.WithContext(authenticatedContext)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return c.Query("token")
}
