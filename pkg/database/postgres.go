package database

import (
	"log"
	"time"

	"github.com/nclsgg/despensa-digital/backend/config"
	authModel "github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	creditsModel "github.com/nclsgg/despensa-digital/backend/internal/modules/credits/model"
	itemModel "github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	profileModel "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	shoppingListModel "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres(cfg *config.Config) (result0 *gorm.DB) {
	__logParams := map[string]any{"cfg": cfg}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "ConnectPostgres"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "ConnectPostgres"), zap.Any("params", __logParams))
	log.Println("Connecting to database...")

	maxAttempts := 3
	delay := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "ConnectPostgres"), zap.Error(err), zap.Any("params", __logParams))
			log.Printf("Error connecting to database: %v\n", err)
			log.Printf("Retrying in %v...\n", delay)
			time.Sleep(delay)
			continue
		}

		log.Println("Connected to database")
		DB = db
		result0 = db
		return
	}

	log.Fatalf("Failed to connect to database after %d attempts", maxAttempts)
	result0 = nil
	return
}

func MigrateItems(db *gorm.DB) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "MigrateItems"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "MigrateItems"), zap.Any("params", __logParams))
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
		&creditsModel.CreditWallet{},
		&creditsModel.CreditTransaction{},
	)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "MigrateItems"), zap.Error(err), zap.Any("params", __logParams))
		log.Fatalf("Error migrating database: %v", err)
	}

	db.Exec(`
		DO $$
		BEGIN
			-- Remove a coluna total_price antiga se existir
			IF EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name='items' AND column_name='total_price'
			) THEN
				EXECUTE 'ALTER TABLE items DROP COLUMN total_price';
			END IF;

			-- Remove coluna price_quantity se existir (não é mais utilizada)
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name='items' AND column_name='price_quantity'
			) THEN
				EXECUTE 'ALTER TABLE items DROP COLUMN price_quantity';
			END IF;

			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name='shopping_list_items' AND column_name='price_quantity'
			) THEN
				EXECUTE 'ALTER TABLE shopping_list_items DROP COLUMN price_quantity';
			END IF;
			
			-- Cria a coluna com a fórmula atualizada: quantidade * preço por unidade
			EXECUTE 'ALTER TABLE items ADD COLUMN total_price numeric GENERATED ALWAYS AS ((quantity * price_per_unit)) STORED';
		END
		$$;
	`)

	log.Println("Database migrated")
}
