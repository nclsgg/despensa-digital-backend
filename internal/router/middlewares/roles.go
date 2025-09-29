package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
	"time"
)

func RoleMiddleware(roles []string) (result0 gin.HandlerFunc) {
	__logParams := map[string]any{"roles": roles}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "RoleMiddleware"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "RoleMiddleware"), zap.Any("params", __logParams))
	result0 = func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Fail(c, 401, "UNAUTHORIZED", "Missing user role in context")
			c.Abort()
			return
		}

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		response.Fail(c, 403, "FORBIDDEN_ROLE", "You do not have permission to access this resource")
		c.Abort()
	}
	return
}
