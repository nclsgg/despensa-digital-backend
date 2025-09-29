package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (result0 *gorm.DB) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "setupTestDB"), zap.Any("result",

			// Realiza a migração para a tabela de User
			result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "setupTestDB"), zap.Any("params", __logParams))
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "setupTestDB"), zap.Error(err), zap.Any("params", __logParams))
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		zap.L().Error("function.error", zap.String("func", "setupTestDB"), zap.Error(err), zap.Any("params", __logParams))
		t.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}
	result0 = db
	return
}

func TestGetUserById(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetUserById"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetUserById"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	userMock := &model.User{
		ID:    uuid.New(),
		Email: "teste@exemplo.com",
	}

	err := db.Create(userMock).Error
	assert.NoError(t, err)
	assert.NotZero(t, userMock.ID)

	user, err := repo.GetUserById(ctx, userMock.ID)

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, userMock, user)
}

func TestGetUserByEmail(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetUserByEmail"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetUserByEmail"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	userMock := &model.User{
		ID:    uuid.New(),
		Email: "teste@exemplo.com",
	}

	err := db.Create(userMock).Error
	assert.NoError(t, err)
	assert.NotZero(t, userMock.ID)

	user, err := repo.GetUserByEmail(ctx, userMock.Email)

	assert.NoError(t, err)
	assert.NotZero(t, user.Email)
	assert.Equal(t, userMock, user)
}

func TestGetAllUsers(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetAllUsers"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetAllUsers"), zap.Any("params", __logParams))
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	usersMock := []model.User{
		{
			ID:    uuid.New(),
			Email: "teste@exemplo.com",
		},
		{
			ID:    uuid.New(),
			Email: "teste2@exemplo.com",
		},
	}

	for i := range usersMock {
		err := db.Create(&usersMock[i]).Error
		assert.NoError(t, err)
		assert.NotZero(t, &usersMock[i].ID)
	}

	user, err := repo.GetAllUsers(ctx)

	assert.NoError(t, err)
	assert.Equal(t, len(usersMock), len(user))
}
