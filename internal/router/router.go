package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/nclsgg/despensa-digital/backend/config"
	authHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/handler"
	authRepo "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/repository"
	authService "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/service"
	userHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/user/handler"
	userRepo "github.com/nclsgg/despensa-digital/backend/internal/modules/user/repository"
	userService "github.com/nclsgg/despensa-digital/backend/internal/modules/user/service"
	middleware "github.com/nclsgg/despensa-digital/backend/internal/router/middlewares"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config, redis *redis.Client) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRepo := authRepo.NewAuthRepository(db)
	authService := authService.NewAuthService(authRepo, cfg, redis)
	authHandlerInstance := authHandler.NewAuthHandler(authService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandlerInstance.Register)
		authGroup.POST("/login", authHandlerInstance.Login)
		authGroup.POST("/refresh", authHandlerInstance.RefreshToken)
	}

	userRepo := userRepo.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)
	userHandler := userHandler.NewUserHandler(userService)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware(cfg, userRepo))
	{
		userGroup.GET("/:id", middleware.RoleMiddleware([]string{"admin"}), userHandler.GetUser)
		userGroup.GET("/me", userHandler.GetCurrentUser)
		userGroup.GET("/all", middleware.RoleMiddleware([]string{"admin"}), userHandler.GetAllUsers)
	}
}
