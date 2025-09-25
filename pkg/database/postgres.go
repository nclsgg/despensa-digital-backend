package database

import (
	"log"
	"time"

	"github.com/nclsgg/despensa-digital/backend/config"
	authModel "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	itemModel "github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	profileModel "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	shoppingListModel "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres(cfg *config.Config) *gorm.DB {
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

func MigrateItems(db *gorm.DB) {
	log.Println("Migrating database...")

	err := db.AutoMigrate(
		&authModel.User{},
		&itemModel.Item{},
		&pantryModel.Pantry{},
		&pantryModel.PantryUser{},
		&itemModel.Item{},
		&itemModel.ItemCategory{},
		&profileModel.Profile{},
		&shoppingListModel.ShoppingList{},
		&shoppingListModel.ShoppingListItem{},
	)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name='items' AND column_name='total_price'
			) THEN
				EXECUTE 'ALTER TABLE items ADD COLUMN total_price numeric GENERATED ALWAYS AS (quantity * price_per_unit) STORED';
			END IF;
		END
		$$;
	`)

	log.Println("Database migrated")
}
