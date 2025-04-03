package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/internal/router"
	"github.com/nclsgg/despensa-digital/backend/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectPostgres(cfg)
	redis := database.ConnectRedis(cfg)

	sqlDB, _ := db.DB()
	db.AutoMigrate(&model.User{})
	defer sqlDB.Close()

	r := gin.Default()

	router.SetupRoutes(r, db, cfg, redis)

	r.Run(":" + cfg.Port)
}
