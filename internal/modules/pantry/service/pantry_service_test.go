package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPantryRepository struct {
	mock.Mock
}

func (m *mockPantryRepository) Create(ctx context.Context, pantry *model.Pantry) (*model.Pantry, error) {
	args := m.Called(ctx, pantry)
	return args.Get(0).(*model.Pantry), args.Error(1)
}

func (m *mockPantryRepository) AddUserToPantry(ctx context.Context, pantryUser *model.PantryUser) error {
	args := m.Called(ctx, pantryUser)
	return args.Error(0)
}

func (m *mockPantryRepository) IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, pantryID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPantryRepository) IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, pantryID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (*model.Pantry, error) {
	args := m.Called(ctx, pantryID)
	return args.Get(0).(*model.Pantry), args.Error(1)
}

func (m *mockPantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*model.Pantry, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Pantry), args.Error(1)
}

func (m *mockPantryRepository) Update(ctx context.Context, pantry *model.Pantry) error {
	args := m.Called(ctx, pantry)
	return args.Error(0)
}

func (m *mockPantryRepository) Delete(ctx context.Context, pantryID uuid.UUID) error {
	args := m.Called(ctx, pantryID)
	return args.Error(0)
}

func (m *mockPantryRepository) RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) error {
	args := m.Called(ctx, pantryID, userID)
	return args.Error(0)
}

func (m *mockPantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) ([]*model.PantryUserInfo, error) {
	args := m.Called(ctx, pantryID)
	return args.Get(0).([]*model.PantryUserInfo), args.Error(1)
}

func TestCreatePantry(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	ownerID := uuid.New()
	pantry := &model.Pantry{
		ID:      uuid.New(),
		Name:    "Minha despensa",
		OwnerID: ownerID,
	}

	repo.On("Create", ctx, mock.AnythingOfType("*model.Pantry")).Return(pantry, nil)
	repo.On("AddUserToPantry", ctx, mock.AnythingOfType("*model.PantryUser")).Return(nil)

	result, err := svc.CreatePantry(ctx, pantry.Name, ownerID)

	assert.NoError(t, err)
	assert.Equal(t, pantry.Name, result.Name)
	repo.AssertExpectations(t)
}

func TestAddUserToPantry(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	ownerID := uuid.New()
	targetUserID := uuid.New()

	repo.On("IsUserOwner", ctx, pantryID, ownerID).Return(true, nil)
	repo.On("IsUserInPantry", ctx, pantryID, targetUserID).Return(false, nil)
	repo.On("AddUserToPantry", ctx, mock.AnythingOfType("*model.PantryUser")).Return(nil)

	err := svc.AddUserToPantry(ctx, pantryID, ownerID, targetUserID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRemoveUserFromPantry_Success(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	ownerID := uuid.New()
	targetUserID := uuid.New()

	repo.On("IsUserOwner", ctx, pantryID, ownerID).Return(true, nil)
	repo.On("RemoveUserFromPantry", ctx, pantryID, targetUserID).Return(nil)

	err := svc.RemoveUserFromPantry(ctx, pantryID, ownerID, targetUserID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRemoveUserFromPantry_BlockOwner(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	ownerID := uuid.New()

	repo.On("IsUserOwner", ctx, pantryID, ownerID).Return(true, nil)

	err := svc.RemoveUserFromPantry(ctx, pantryID, ownerID, ownerID)
	assert.EqualError(t, err, "owner cannot remove themselves")
	repo.AssertExpectations(t)
}

func TestGetPantry_Success(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	pantry := &model.Pantry{
		ID:      pantryID,
		Name:    "Test",
		OwnerID: userID,
	}

	repo.On("IsUserInPantry", ctx, pantryID, userID).Return(true, nil)
	repo.On("GetByID", ctx, pantryID).Return(pantry, nil)

	result, err := svc.GetPantry(ctx, pantryID, userID)
	assert.NoError(t, err)
	assert.Equal(t, pantry.ID, result.ID)
	repo.AssertExpectations(t)
}

func TestGetPantry_NotMember(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	repo.On("IsUserInPantry", ctx, pantryID, userID).Return(false, nil)

	_, err := svc.GetPantry(ctx, pantryID, userID)
	assert.ErrorIs(t, err, service.ErrUnauthorized)
	repo.AssertExpectations(t)
}

func TestListUsersInPantry_Success(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	expectedUsers := []*model.PantryUserInfo{
		{UserID: uuid.New(), Email: "test1@example.com", Role: "member"},
		{UserID: uuid.New(), Email: "test2@example.com", Role: "owner"},
	}

	repo.On("IsUserInPantry", ctx, pantryID, userID).Return(true, nil)
	repo.On("ListUsersInPantry", ctx, pantryID).Return(expectedUsers, nil)

	result, err := svc.ListUsersInPantry(ctx, pantryID, userID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestListUsersInPantry_Unauthorized(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	repo.On("IsUserInPantry", ctx, pantryID, userID).Return(false, nil)

	_, err := svc.ListUsersInPantry(ctx, pantryID, userID)
	assert.EqualError(t, err, "user is not in the pantry")
	repo.AssertExpectations(t)
}

func TestUpdatePantry_Success(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	pantry := &model.Pantry{
		ID:      pantryID,
		Name:    "Old Name",
		OwnerID: userID,
	}

	repo.On("IsUserOwner", ctx, pantryID, userID).Return(true, nil)
	repo.On("GetByID", ctx, pantryID).Return(pantry, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*model.Pantry")).Return(nil)

	err := svc.UpdatePantry(ctx, pantryID, userID, "New Name")
	assert.NoError(t, err)
	assert.Equal(t, "New Name", pantry.Name)
	repo.AssertExpectations(t)
}

func TestUpdatePantry_Unauthorized(t *testing.T) {
	repo := new(mockPantryRepository)
	svc := service.NewPantryService(repo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	repo.On("IsUserOwner", ctx, pantryID, userID).Return(false, nil)

	err := svc.UpdatePantry(ctx, pantryID, userID, "Name")
	assert.EqualError(t, err, service.ErrUnauthorized.Error())
	repo.AssertExpectations(t)
}
