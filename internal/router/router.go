package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	authHandler "github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/handler"
	authRepo "github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/repository"
	authService "github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/service"
	userHandler "github.com/nclsgg/dispensa-digital/backend/internal/modules/user/handler"
	userRepo "github.com/nclsgg/dispensa-digital/backend/internal/modules/user/repository"
	userService "github.com/nclsgg/dispensa-digital/backend/internal/modules/user/service"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRepo := authRepo.NewAuthRepository(db)
	authService := authService.NewAuthService(authRepo)
	authHandler := authHandler.NewAuthHandler(authService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
	}

	userRepo := userRepo.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)
	userHandler := userHandler.NewUserHandler(userService)

	userGroup := r.Group("/user")
	{
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.GET("/all", userHandler.GetAllUsers)
	}
}
