package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

func TestGetUserById(t *testing.T) {
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

func TestGetAllUsers(t *testing.T) {
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
