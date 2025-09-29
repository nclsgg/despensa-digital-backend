package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB configura um banco de dados SQLite in-memory para os testes.
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

func TestCreateUserAndGetUser(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String(

			// Cria um novo usuário
			"func", "TestCreateUserAndGetUser"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestCreateUserAndGetUser"), zap.Any("params",

		// Recupera o usuário pelo email
		__logParams))
	db := setupTestDB(t)
	repo := repository.NewAuthRepository(db)
	ctx := context.Background()

	user := &model.User{
		Email:    "teste@exemplo.com",
		Password: "senha123",
	}
	err := repo.CreateUser(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	foundUser, err := repo.GetUser(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestGetUserById(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String(

			// Cria um usuário para teste
			"func", "TestGetUserById"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetUserById"), zap.Any("params",

		// Busca o usuário pelo ID
		__logParams))
	db := setupTestDB(t)
	repo := repository.NewAuthRepository(db)
	ctx := context.Background()

	user := &model.User{
		Email:    "teste2@exemplo.com",
		Password: "senha456",
	}
	err := repo.CreateUser(ctx, user)
	assert.NoError(t, err)

	foundUser, err := repo.GetUserById(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)
}
