package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/nclsgg/despensa-digital/backend/cmd/server/docs"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/nclsgg/despensa-digital/backend/config"
	authHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/handler"
	authRepo "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/repository"
	authService "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/service"
	itemHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/item/handler"
	itemRepo "github.com/nclsgg/despensa-digital/backend/internal/modules/item/repository"
	itemService "github.com/nclsgg/despensa-digital/backend/internal/modules/item/service"
	pantryHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/handler"
	pantryRepo "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/repository"
	pantryService "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/service"
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

	// Auth routes
	authRepoInstance := authRepo.NewAuthRepository(db)
	authServiceInstance := authService.NewAuthService(authRepoInstance, cfg, redis)
	authHandlerInstance := authHandler.NewAuthHandler(authServiceInstance)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandlerInstance.Register)
		authGroup.POST("/login", authHandlerInstance.Login)
		authGroup.POST("/logout", authHandlerInstance.Logout)
		authGroup.POST("/refresh", authHandlerInstance.RefreshToken)
	}

	// User routes
	userRepoInstance := userRepo.NewUserRepository(db)
	userServiceInstance := userService.NewUserService(userRepoInstance)
	userHandlerInstance := userHandler.NewUserHandler(userServiceInstance)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	{
		userGroup.GET("/:id", middleware.RoleMiddleware([]string{"admin"}), userHandlerInstance.GetUser)
		userGroup.GET("/me", userHandlerInstance.GetCurrentUser)
		userGroup.GET("/all", middleware.RoleMiddleware([]string{"admin"}), userHandlerInstance.GetAllUsers)
	}

	// Pantry routes
	pantryRepoInstance := pantryRepo.NewPantryRepository(db)
	pantryServiceInstance := pantryService.NewPantryService(pantryRepoInstance, userRepoInstance)
	pantryHandlerInstance := pantryHandler.NewPantryHandler(pantryServiceInstance)

	pantryGroup := r.Group("/pantries")
	pantryGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	{
		pantryGroup.POST("", pantryHandlerInstance.CreatePantry)
		pantryGroup.GET("", pantryHandlerInstance.ListPantries)
		pantryGroup.GET("/:id", pantryHandlerInstance.GetPantry)
		pantryGroup.DELETE("/:id", pantryHandlerInstance.DeletePantry)
		pantryGroup.PUT("/:id", pantryHandlerInstance.UpdatePantry)
		pantryGroup.POST("/:id/users", pantryHandlerInstance.AddUserToPantry)
		pantryGroup.DELETE("/:id/users", pantryHandlerInstance.RemoveUserFromPantry)
		pantryGroup.GET("/:id/users", pantryHandlerInstance.ListUsersInPantry)
	}

	// Item routes
	itemRepoInstance := itemRepo.NewItemRepository(db)
	itemServiceInstance := itemService.NewItemService(itemRepoInstance, pantryRepoInstance)
	itemHandlerInstance := itemHandler.NewItemHandler(itemServiceInstance)

	itemGroup := r.Group("/items")
	itemGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	{
		itemGroup.POST("", itemHandlerInstance.CreateItem)
		itemGroup.GET("/pantry/:id", itemHandlerInstance.ListItems)
		itemGroup.GET("/:id", itemHandlerInstance.GetItem)
		itemGroup.PUT("/:id", itemHandlerInstance.UpdateItem)
		itemGroup.DELETE("/:id", itemHandlerInstance.DeleteItem)
	}

	// Item Category routes
	itemCategoryRepoInstance := itemRepo.NewItemCategoryRepository(db)
	itemCategoryServiceInstance := itemService.NewItemCategoryService(itemCategoryRepoInstance, pantryRepoInstance)
	itemCategoryHandlerInstance := itemHandler.NewItemCategoryHandler(itemCategoryServiceInstance)

	itemCategoryGroup := r.Group("/item-categories")
	itemCategoryGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	{
		itemCategoryGroup.POST("", itemCategoryHandlerInstance.CreateItemCategory)
		itemCategoryGroup.POST("/default", middleware.RoleMiddleware([]string{"admin"}), itemCategoryHandlerInstance.CreateDefaultItemCategory)
		itemCategoryGroup.POST("/from-default/:default_id/pantry/:pantry_id", itemCategoryHandlerInstance.CloneDefaultCategoryToPantry)
		itemCategoryGroup.GET("/pantry/:id", itemCategoryHandlerInstance.ListItemCategoriesByPantry)
		itemCategoryGroup.GET("/:id", itemCategoryHandlerInstance.GetItemCategory)
		itemCategoryGroup.PUT("/:id", itemCategoryHandlerInstance.UpdateItemCategory)
		itemCategoryGroup.DELETE("/:id", itemCategoryHandlerInstance.DeleteItemCategory)
		itemCategoryGroup.GET("/user", itemCategoryHandlerInstance.ListItemCategoriesByUser)
	}

	// Swagger routes
	r.GET(
		"/swagger/*any",
		middleware.AuthMiddleware(cfg, userRepoInstance),
		middleware.RoleMiddleware([]string{"admin"}),
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}
