package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nclsgg/despensa-digital/backend/config"
	authModel "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	userModel "github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

// ProfileCompleteMiddleware verifies if user has completed their profile
func ProfileCompleteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for complete-profile endpoint itself
		if strings.Contains(c.Request.URL.Path, "complete-profile") {
			c.Next()
			return
		}

		userInterface, exists := c.Get("user")
		if !exists {
			response.Fail(c, 401, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		user := userInterface.(*userModel.User)
		if !user.ProfileCompleted {
			response.Fail(c, 403, "PROFILE_INCOMPLETE", "Profile must be completed to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}

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
	claims := &authModel.MyClaims{}

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

	c.Set("userID", claims.UserID) // Mudança de "user_id" para "userID" para consistência
	c.Set("user", user)
	c.Set("email", user.Email)
	c.Set("role", user.Role)

	c.Next()
	}
}
