package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/nclsgg/despensa-digital/backend/cmd/server/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

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

	// LLM module imports
	llmHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/handler"
	llmService "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"

	// Recipe module imports
	recipeHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/handler"
	recipeService "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/service"

	// Profile module imports
	profileHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/handler"
	profileRepo "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/repository"
	profileService "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/service"

	// Shopping list module imports
	shoppingListHandler "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/handler"
	shoppingListRepo "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/repository"
	shoppingListService "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/service"

	middleware "github.com/nclsgg/despensa-digital/backend/internal/router/middlewares"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	__logParams := map[string]any{"r": r, "db": db, "cfg": cfg}
	__logStart := time.Now()
	defer func() {
		zap.

			// User repository (needed for auth middleware)
			L().Info("function.exit", zap.String("func", "SetupRoutes"), zap.Any("result", nil), zap.Duration("duration",

			// OAuth-only auth routes
			time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "SetupRoutes"), zap.Any("params", __logParams))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	userRepoInstance := userRepo.NewUserRepository(db)

	authRepoInstance := authRepo.NewAuthRepository(db)
	authServiceInstance := authService.NewAuthService(authRepoInstance, cfg)

	// OAuth handler
	oauthHandlerInstance := authHandler.NewOAuthHandler(authServiceInstance, cfg)
	oauthHandlerInstance.InitOAuth()

	authGroup := r.Group("/api/v1/auth")
	{
		// OAuth routes only
		authGroup.GET("/oauth/:provider", oauthHandlerInstance.OAuthLogin)
		authGroup.GET("/oauth/:provider/callback", oauthHandlerInstance.OAuthCallback)

		// Protected profile completion route
		authGroup.PATCH("/complete-profile", middleware.AuthMiddleware(cfg, userRepoInstance), oauthHandlerInstance.CompleteProfile)
	}

	// User routes
	userServiceInstance := userService.NewUserService(userRepoInstance)
	userHandlerInstance := userHandler.NewUserHandler(userServiceInstance)

	userGroup := r.Group("/api/v1/user")
	userGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	userGroup.Use(middleware.ProfileCompleteMiddleware())
	{
		userGroup.GET("/:id", middleware.RoleMiddleware([]string{"admin"}), userHandlerInstance.GetUser)
		userGroup.GET("/me", userHandlerInstance.GetCurrentUser)
		userGroup.GET("/all", middleware.RoleMiddleware([]string{"admin"}), userHandlerInstance.GetAllUsers)
	}

	// LLM routes (needed for recipe handlers)
	llmServiceInstance := llmService.NewLLMService()
	llmHandlerInstance := llmHandler.NewLLMHandler(llmServiceInstance)

	// Recipe routes setup (needed for pantry ingredients endpoint)
	pantryRepoInstance := pantryRepo.NewPantryRepository(db)
	itemRepoInstance := itemRepo.NewItemRepository(db)
	pantryServiceInstance := pantryService.NewPantryService(pantryRepoInstance, userRepoInstance, itemRepoInstance)
	itemServiceInstance := itemService.NewItemService(itemRepoInstance, pantryRepoInstance)

	// Profile module setup
	profileRepoInstance := profileRepo.NewProfileRepository(db)
	profileServiceInstance := profileService.NewProfileService(profileRepoInstance)
	profileHandlerInstance := profileHandler.NewProfileHandler(profileServiceInstance)

	// Shopping list module setup
	shoppingListRepoInstance := shoppingListRepo.NewShoppingListRepository(db)
	shoppingListServiceInstance := shoppingListService.NewShoppingListService(
		shoppingListRepoInstance,
		pantryRepoInstance,
		itemRepoInstance,
		profileRepoInstance,
		llmServiceInstance,
	)
	shoppingListHandlerInstance := shoppingListHandler.NewShoppingListHandler(shoppingListServiceInstance)

	recipeServiceInstance := recipeService.NewRecipeService(
		llmServiceInstance,
		itemRepoInstance,
		pantryServiceInstance,
	)
	recipeHandlerInstance := recipeHandler.NewRecipeHandler(recipeServiceInstance, llmServiceInstance)

	// Pantry routes
	pantryHandlerInstance := pantryHandler.NewPantryHandler(pantryServiceInstance, itemServiceInstance)

	pantryGroup := r.Group("/api/v1/pantries")
	pantryGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	pantryGroup.Use(middleware.ProfileCompleteMiddleware())
	{
		pantryGroup.POST("", pantryHandlerInstance.CreatePantry)
		pantryGroup.GET("", pantryHandlerInstance.ListPantries)
		pantryGroup.GET("/:id", pantryHandlerInstance.GetPantry)
		pantryGroup.DELETE("/:id", pantryHandlerInstance.DeletePantry)
		pantryGroup.PUT("/:id", pantryHandlerInstance.UpdatePantry)
		pantryGroup.POST("/:id/users", pantryHandlerInstance.AddUserToPantry)
		pantryGroup.DELETE("/:id/users", pantryHandlerInstance.RemoveUserFromPantry)
		pantryGroup.GET("/:id/users", pantryHandlerInstance.ListUsersInPantry)
		pantryGroup.GET("/:id/ingredients", recipeHandlerInstance.GetAvailableIngredients)
	}

	// Item routes - reuse the itemRepoInstance
	itemHandlerInstance := itemHandler.NewItemHandler(itemServiceInstance)

	itemGroup := r.Group("/api/v1/items")
	itemGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	itemGroup.Use(middleware.ProfileCompleteMiddleware())
	{
		itemGroup.POST("", itemHandlerInstance.CreateItem)
		itemGroup.GET("/pantry/:id", itemHandlerInstance.ListItems)
		itemGroup.POST("/pantry/:id/filter", itemHandlerInstance.FilterItems)
		itemGroup.GET("/:id", itemHandlerInstance.GetItem)
		itemGroup.PUT("/:id", itemHandlerInstance.UpdateItem)
		itemGroup.DELETE("/:id", itemHandlerInstance.DeleteItem)
	}

	// Item Category routes
	itemCategoryRepoInstance := itemRepo.NewItemCategoryRepository(db)
	itemCategoryServiceInstance := itemService.NewItemCategoryService(itemCategoryRepoInstance, pantryRepoInstance)
	itemCategoryHandlerInstance := itemHandler.NewItemCategoryHandler(itemCategoryServiceInstance)

	itemCategoryGroup := r.Group("/api/v1/item-categories")
	itemCategoryGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	itemCategoryGroup.Use(middleware.ProfileCompleteMiddleware())
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

	llmGroup := r.Group("/api/v1/llm")
	llmGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	llmGroup.Use(middleware.ProfileCompleteMiddleware())
	{
		llmGroup.POST("/chat", llmHandlerInstance.ProcessChatRequest)
		llmGroup.POST("/process", llmHandlerInstance.ProcessLLMRequest)
		llmGroup.POST("/prompt/build", llmHandlerInstance.BuildPrompt)
		llmGroup.GET("/providers/status", llmHandlerInstance.GetProviderStatus)
		llmGroup.POST("/providers/config", llmHandlerInstance.ConfigureProvider)
		llmGroup.GET("/providers/available", llmHandlerInstance.GetAvailableProviders)
		llmGroup.POST("/providers/switch", llmHandlerInstance.SwitchProvider)
		llmGroup.POST("/providers/test", llmHandlerInstance.TestProvider)
	}

	// Profile routes
	profileGroup := r.Group("/api/v1/profile")
	profileGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	{
		profileGroup.POST("", profileHandlerInstance.CreateProfile)
		profileGroup.GET("", profileHandlerInstance.GetProfile)
		profileGroup.PUT("", profileHandlerInstance.UpdateProfile)
		profileGroup.DELETE("", profileHandlerInstance.DeleteProfile)
	}

	// Shopping list routes
	shoppingListGroup := r.Group("/api/v1/shopping-lists")
	shoppingListGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	shoppingListGroup.Use(middleware.ProfileCompleteMiddleware())
	{
		shoppingListGroup.POST("", shoppingListHandlerInstance.CreateShoppingList)
		shoppingListGroup.GET("", shoppingListHandlerInstance.GetShoppingLists)
		shoppingListGroup.GET("/:id", shoppingListHandlerInstance.GetShoppingList)
		shoppingListGroup.PUT("/:id", shoppingListHandlerInstance.UpdateShoppingList)
		shoppingListGroup.DELETE("/:id", shoppingListHandlerInstance.DeleteShoppingList)
		shoppingListGroup.PUT("/:id/items/:itemId", shoppingListHandlerInstance.UpdateShoppingListItem)
		shoppingListGroup.DELETE("/:id/items/:itemId", shoppingListHandlerInstance.DeleteShoppingListItem)
		shoppingListGroup.POST("/generate", shoppingListHandlerInstance.GenerateAIShoppingList)
	}

	recipeGroup := r.Group("/api/v1/recipes")
	recipeGroup.Use(middleware.AuthMiddleware(cfg, userRepoInstance))
	recipeGroup.Use(middleware.ProfileCompleteMiddleware())
	{
		recipeGroup.POST("/generate", recipeHandlerInstance.GenerateRecipe)
		recipeGroup.GET("/ingredients", recipeHandlerInstance.GetAvailableIngredients)
		recipeGroup.GET("/pantries/:pantry_id/ingredients", recipeHandlerInstance.GetAvailableIngredients)
		recipeGroup.POST("/chat", recipeHandlerInstance.ChatWithLLM)
		recipeGroup.GET("/providers", recipeHandlerInstance.GetLLMProviders)
		recipeGroup.POST("/providers/set", recipeHandlerInstance.SetLLMProvider)
		recipeGroup.POST("/tokens/estimate", recipeHandlerInstance.EstimateTokens)
	}

	// Swagger routes
	r.GET(
		"/swagger/*any",
		middleware.AuthMiddleware(cfg, userRepoInstance),
		middleware.RoleMiddleware([]string{"admin"}),
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}
