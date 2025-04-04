package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(id)
	if usr, ok := args.Get(0).(*model.User); ok {
		return usr, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(email)
	if usr, ok := args.Get(0).(*model.User); ok {
		return usr, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	args := m.Called()
	if usrs, ok := args.Get(0).([]model.User); ok {
		return usrs, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetUser(t *testing.T) {
	repo := new(mockUserRepository)
	service := service.NewUserService(repo)

	userID := uuid.New()

	userMock := &model.User{
		ID:    userID,
		Name:  "Test",
		Email: "test@example.com",
	}

	repo.On("GetUserById", userID).Return(userMock, nil)

	user, err := service.GetUserById(context.Background(), userID)
	if err != nil {
		t.Fatalf("Erro ao buscar usuário: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, userMock, user)

	repo.AssertExpectations(t)
}

func TestGetAllUsers(t *testing.T) {
	repo := new(mockUserRepository)
	service := service.NewUserService(repo)

	usersMock := []model.User{
		{
			ID:    uuid.New(),
			Name:  "Test",
			Email: "test1@example.com",
		},
		{
			ID:    uuid.New(),
			Name:  "Test 2",
			Email: "text2@example.com",
		},
	}

	repo.On("GetAllUsers").Return(usersMock, nil)

	users, err := service.GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("Erro ao buscar usuários: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, usersMock, users)

	repo.AssertExpectations(t)
}
