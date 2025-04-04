package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

func AuthMiddleware(cfg *config.Config, userRepo domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Fail(c, 401, "UNAUTHORIZED", "Missing or invalid Authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &model.MyClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			response.Fail(c, 401, "INVALID_TOKEN", "Invalid or corrupted JWT")
			c.Abort()
			return
		}

		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			response.Fail(c, 401, "TOKEN_EXPIRED", "Token has expired")
			c.Abort()
			return
		}

		if claims.Issuer != cfg.JWTIssuer {
			response.Fail(c, 401, "INVALID_ISSUER", "Invalid token issuer")
			c.Abort()
			return
		}

		if len(claims.Audience) == 0 || claims.Audience[0] != cfg.JWTAudience {
			response.Fail(c, 401, "INVALID_AUDIENCE", "Invalid token audience")
			c.Abort()
			return
		}

		user, err := userRepo.GetUserById(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Fail(c, 401, "USER_NOT_FOUND", "User not found")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", user.Email)
		c.Set("role", user.Role)

		c.Next()
	}
}
