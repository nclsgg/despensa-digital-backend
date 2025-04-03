package repository_test

import (
	"context"
	"testing"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB configura um banco de dados SQLite in-memory para os testes.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	// Realiza a migração para a tabela de User
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}

	return db
}

func TestCreateUserAndGetUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuthRepository(db)
	ctx := context.Background()

	// Cria um novo usuário
	user := &model.User{
		Email:    "teste@exemplo.com",
		Password: "senha123",
	}
	err := repo.CreateUser(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Recupera o usuário pelo email
	foundUser, err := repo.GetUser(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestGetUserById(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuthRepository(db)
	ctx := context.Background()

	// Cria um usuário para teste
	user := &model.User{
		Email:    "teste2@exemplo.com",
		Password: "senha456",
	}
	err := repo.CreateUser(ctx, user)
	assert.NoError(t, err)

	// Busca o usuário pelo ID
	foundUser, err := repo.GetUserById(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)
}
