package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (result0 *gorm.DB) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "setupTestDB"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "setupTestDB"), zap.Any("params", __logParams))
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.Pantry{}, &model.PantryUser{})
	assert.NoError(t, err)
	result0 = db
	return
}

func TestCreateAndGetByID(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestCreateAndGetByID"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestCreateAndGetByID"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewPantryRepository(db)

	ctx := context.Background()
	ownerID := uuid.New()
	pantry := &model.Pantry{
		ID:      uuid.New(),
		Name:    "Cozinha",
		OwnerID: ownerID,
	}

	_, err := repo.Create(ctx, pantry)
	assert.NoError(t, err)

	found, err := repo.GetByID(ctx, pantry.ID)
	assert.NoError(t, err)
	assert.Equal(t, pantry.Name, found.Name)
}

func TestAddUserToPantryAndIsUserInPantry(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestAddUserToPantryAndIsUserInPantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestAddUserToPantryAndIsUserInPantry"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewPantryRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	pantryID := uuid.New()

	db.Create(&model.Pantry{ID: pantryID, Name: "Despensa", OwnerID: userID})

	pantryUser := &model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantryID,
		UserID:   userID,
		Role:     "owner",
	}

	err := repo.AddUserToPantry(ctx, pantryUser)
	assert.NoError(t, err)

	isMember, err := repo.IsUserInPantry(ctx, pantryID, userID)
	assert.NoError(t, err)
	assert.True(t, isMember)
}

func TestRemoveUserFromPantry(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestRemoveUserFromPantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestRemoveUserFromPantry"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewPantryRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	pantryID := uuid.New()

	db.Create(&model.Pantry{ID: pantryID, Name: "Sala", OwnerID: userID})
	db.Create(&model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantryID,
		UserID:   userID,
		Role:     "owner",
	})

	err := repo.RemoveUserFromPantry(ctx, pantryID, userID)
	assert.NoError(t, err)

	isMember, _ := repo.IsUserInPantry(ctx, pantryID, userID)
	assert.False(t, isMember)
}

func TestGetByUser(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetByUser"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetByUser"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewPantryRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	pantryID := uuid.New()

	db.Create(&model.Pantry{ID: pantryID, Name: "Churrasqueira", OwnerID: userID})
	db.Create(&model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantryID,
		UserID:   userID,
		Role:     "member",
	})

	pantries, err := repo.GetByUser(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, pantries, 1)

	_ = repo.RemoveUserFromPantry(ctx, pantryID, userID)

	pantries, _ = repo.GetByUser(ctx, userID)
	assert.Len(t, pantries, 0)
}

func TestListUsersInPantry(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestListUsersInPantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestListUsersInPantry"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewPantryRepository(db)
	ctx := context.Background()

	user1 := uuid.New()
	user2 := uuid.New()
	pantryID := uuid.New()

	db.Exec(`CREATE TABLE users (id UUID PRIMARY KEY, email TEXT NOT NULL)`)
	db.Exec(`INSERT INTO users (id, email) VALUES (?, ?), (?, ?)`, user1, "user1@email.com", user2, "user2@email.com")

	db.Create(&model.Pantry{
		ID:      pantryID,
		Name:    "Geral",
		OwnerID: user1,
	})

	db.Create(&model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantryID,
		UserID:   user1,
		Role:     "owner",
	})

	db.Create(&model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantryID,
		UserID:   user2,
		Role:     "member",
	})

	users, err := repo.ListUsersInPantry(ctx, pantryID)

	assert.NoError(t, err)
	assert.Len(t, users, 2)

	assert.Equal(t, "user1@email.com", users[0].Email)
	assert.Equal(t, "owner", users[0].Role)
}
