package middlewares

import (
	"errors"
	"fmt"
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/pkg/ratelimit"

	"github.com/gin-gonic/gin"
)

func IPRateLimiterMiddleware(ratelimiter ratelimit.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ip:%s", ip)
		allowed, err := ratelimiter.Allow(c.Request.Context(), key)
		if err != nil {
			_ = c.Error(pkgerrors.NewInternalError(err))
			c.Abort()
			return
		}

		if !allowed {
			_ = c.Error(pkgerrors.NewTooManyRequestsError("Too many requests from this IP", err))
			c.Abort()
			return
		}

		c.Next()
	}
}

func UserRateLimiterMiddleware(ratelimiter ratelimit.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal, ok := identity.PrincipalFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(pkgerrors.NewInternalError(
				errors.New("UserRateLimiter used without AuthMiddleware"),
			))
			c.Abort()
			return
		}

		key := fmt.Sprintf("userId:%d", principal.UserID)
		allowed, err := ratelimiter.Allow(c.Request.Context(), key)
		if err != nil {
			_ = c.Error(pkgerrors.NewInternalError(err))
			c.Abort()
			return
		}

		if !allowed {
			_ = c.Error(pkgerrors.NewTooManyRequestsError("Too many requests from this user", err))
			c.Abort()
			return
		}

		c.Next()
	}
}
