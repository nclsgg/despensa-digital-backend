package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	_ "github.com/nclsgg/despensa-digital/backend/cmd/server/docs"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/router"
	"github.com/nclsgg/despensa-digital/backend/pkg/database"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
)

// @title Despensa Digital API
// @version 1.0
// @description API da aplicação Despensa Digital
// @termsOfService http://swagger.io/terms/

// @contact.name Nicolas Guadagno
// @contact.url http://github.com/nclsgg
// @contact.email nicolasguadagno@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:5310
// @BasePath /
func main() {
	// Initialize logger based on environment
	ginMode := os.Getenv("GIN_MODE")

	var logger *zap.Logger

	if ginMode == "release" {
		// Production configuration
		logger = appLogger.NewProduction()
		appLogger.WithAppInfo(logger, "despensa-digital", "1.0.0")
	} else {
		// Development configuration
		logger = appLogger.NewDevelopment()
		appLogger.WithAppInfo(logger, "despensa-digital", "1.0.0")
	}

	defer logger.Sync()

	logger.Info("Starting Despensa Digital API",
		zap.String(appLogger.FieldEnvironment, ginMode),
	)

	cfg := config.LoadConfig()

	logger.Info("Configuration loaded",
		zap.String("port", cfg.Port),
		zap.String("cors_origin", cfg.CorsOrigin),
	)

	db := database.ConnectPostgres(cfg)

	sqlDB, _ := db.DB()
	database.MigrateItems(db)
	defer sqlDB.Close()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CorsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.SetupRoutes(r, db, cfg, logger)

	logger.Info("Server starting",
		zap.String("port", cfg.Port),
	)

	r.Run(":" + cfg.Port)
}
