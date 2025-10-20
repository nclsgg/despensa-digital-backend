package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	zap.ReplaceGlobals(zap.NewNop())
}

func setupItemTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&model.Item{}))

	return db
}

func TestItemRepositoryCountByPantryIDCountsEntries(t *testing.T) {
	db := setupItemTestDB(t)
	repo := NewItemRepository(db)

	pantryID := uuid.New()
	addedBy := uuid.New()
	now := time.Now().UTC()

	cheeseID := uuid.New()
	cheese := &model.Item{
		ID:           cheeseID,
		PantryID:     pantryID,
		AddedBy:      addedBy,
		Name:         "Queijo",
		Quantity:     400, // grams
		PricePerUnit: 30,  // price per kilogram
		Unit:         "g",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, db.WithContext(context.Background()).Create(cheese).Error)

	count, err := repo.CountByPantryID(context.Background(), pantryID)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	butter := &model.Item{
		ID:           uuid.New(),
		PantryID:     pantryID,
		AddedBy:      addedBy,
		Name:         "Manteiga",
		Quantity:     2,
		PricePerUnit: 10,
		Unit:         "un",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, db.WithContext(context.Background()).Create(butter).Error)

	count, err = repo.CountByPantryID(context.Background(), pantryID)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}
