package main

import (
	"time"

	"github.com/gin-contrib/cors"
	_ "github.com/nclsgg/despensa-digital/backend/cmd/server/docs"

	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/config"
	authModel "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	"github.com/nclsgg/despensa-digital/backend/internal/router"
	"github.com/nclsgg/despensa-digital/backend/pkg/database"
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
	cfg := config.LoadConfig()
	db := database.ConnectPostgres(cfg)
	redis := database.ConnectRedis(cfg)

	sqlDB, _ := db.DB()
	db.AutoMigrate(&authModel.User{})
	db.AutoMigrate(&pantryModel.Pantry{})
	db.AutoMigrate(&pantryModel.PantryUser{})
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

	router.SetupRoutes(r, db, cfg, redis)

	r.Run(":" + cfg.Port)
}
