package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nclsgg/dispensa-digital/backend/config"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/dispensa-digital/backend/internal/router"
	"github.com/nclsgg/dispensa-digital/backend/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	db := database.Connect(cfg)

	sqlDB, _ := db.DB()
	db.AutoMigrate(&model.User{})
	defer sqlDB.Close()

	r := gin.Default()

	router.SetupRoutes(r, db)

	r.Run(":" + cfg.Port)
}
