package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}
