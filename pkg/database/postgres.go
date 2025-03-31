package database

import (
	"log"
	"time"

	"github.com/nclsgg/dispensa-digital/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) *gorm.DB {
	log.Println("Connecting to database...")

	maxAttempts := 3
	delay := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Printf("Error connecting to database: %v\n", err)
			log.Printf("Retrying in %v...\n", delay)
			time.Sleep(delay)
			continue
		}

		log.Println("Connected to database")
		DB = db
		return db
	}

	log.Fatalf("Failed to connect to database after %d attempts", maxAttempts)
	return nil
}
