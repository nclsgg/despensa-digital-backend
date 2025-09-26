package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	shoppingDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/dto"
	shoppingModel "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/service"
)

type mockShoppingListRepository struct {
	mock.Mock
}

func (m *mockShoppingListRepository) Create(ctx context.Context, shoppingList *shoppingModel.ShoppingList) error {
	args := m.Called(ctx, shoppingList)
	return args.Error(0)
}

func (m *mockShoppingListRepository) GetByID(ctx context.Context, id uuid.UUID) (*shoppingModel.ShoppingList, error) {
	args := m.Called(ctx, id)
	if list, ok := args.Get(0).(*shoppingModel.ShoppingList); ok {
		return list, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockShoppingListRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shoppingModel.ShoppingList, error) {
	args := m.Called(ctx, userID, limit, offset)
	if lists, ok := args.Get(0).([]*shoppingModel.ShoppingList); ok {
		return lists, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockShoppingListRepository) Update(ctx context.Context, shoppingList *shoppingModel.ShoppingList) error {
	args := m.Called(ctx, shoppingList)
	return args.Error(0)
}

func (m *mockShoppingListRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockShoppingListRepository) CreateItem(ctx context.Context, item *shoppingModel.ShoppingListItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *mockShoppingListRepository) UpdateItem(ctx context.Context, item *shoppingModel.ShoppingListItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *mockShoppingListRepository) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *mockShoppingListRepository) GetItemsByShoppingListID(ctx context.Context, shoppingListID uuid.UUID) ([]*shoppingModel.ShoppingListItem, error) {
	args := m.Called(ctx, shoppingListID)
	if items, ok := args.Get(0).([]*shoppingModel.ShoppingListItem); ok {
		return items, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockShoppingListRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

type mockPantryRepository struct {
	mock.Mock
}

func (m *mockPantryRepository) Create(ctx context.Context, pantry *pantryModel.Pantry) (*pantryModel.Pantry, error) {
	args := m.Called(ctx, pantry)
	if p, ok := args.Get(0).(*pantryModel.Pantry); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPantryRepository) Delete(ctx context.Context, pantryID uuid.UUID) error {
	args := m.Called(ctx, pantryID)
	return args.Error(0)
}

func (m *mockPantryRepository) Update(ctx context.Context, pantry *pantryModel.Pantry) error {
	args := m.Called(ctx, pantry)
	return args.Error(0)
}

func (m *mockPantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (*pantryModel.Pantry, error) {
	args := m.Called(ctx, pantryID)
	if pantry, ok := args.Get(0).(*pantryModel.Pantry); ok {
		return pantry, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*pantryModel.Pantry, error) {
	args := m.Called(ctx, userID)
	if pantries, ok := args.Get(0).([]*pantryModel.Pantry); ok {
		return pantries, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPantryRepository) IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, pantryID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPantryRepository) IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, pantryID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPantryRepository) AddUserToPantry(ctx context.Context, pantryUser *pantryModel.PantryUser) error {
	args := m.Called(ctx, pantryUser)
	return args.Error(0)
}

func (m *mockPantryRepository) RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) error {
	args := m.Called(ctx, pantryID, userID)
	return args.Error(0)
}

func (m *mockPantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) ([]*pantryModel.PantryUserInfo, error) {
	args := m.Called(ctx, pantryID)
	if users, ok := args.Get(0).([]*pantryModel.PantryUserInfo); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func newService(repo *mockShoppingListRepository, pantryRepo *mockPantryRepository) shoppingDomain.ShoppingListService {
	return service.NewShoppingListService(repo, pantryRepo, nil, nil, nil)
}

func TestShoppingListService_CreateShoppingList_PantryAccessDenied(t *testing.T) {
	repo := new(mockShoppingListRepository)
	pantryRepo := new(mockPantryRepository)
	service := newService(repo, pantryRepo)

	userID := uuid.New()
	pantryID := uuid.New()
	input := dto.CreateShoppingListDTO{
		Name:        "Weekly",
		PantryID:    &pantryID,
		TotalBudget: 100,
	}

	pantryRepo.On("IsUserInPantry", mock.Anything, pantryID, userID).Return(false, nil).Once()

	result, err := service.CreateShoppingList(context.Background(), userID, input)
	require.ErrorIs(t, err, shoppingDomain.ErrPantryAccessDenied)
	require.Nil(t, result)

	pantryRepo.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestShoppingListService_GetShoppingListByID_NotFound(t *testing.T) {
	repo := new(mockShoppingListRepository)
	pantryRepo := new(mockPantryRepository)
	service := newService(repo, pantryRepo)

	listID := uuid.New()
	repo.On("GetByID", mock.Anything, listID).Return((*shoppingModel.ShoppingList)(nil), gorm.ErrRecordNotFound).Once()

	result, err := service.GetShoppingListByID(context.Background(), uuid.New(), listID)
	require.ErrorIs(t, err, shoppingDomain.ErrShoppingListNotFound)
	require.Nil(t, result)

	repo.AssertExpectations(t)
	pantryRepo.AssertExpectations(t)
}

func TestShoppingListService_GetShoppingListByID_Unauthorized(t *testing.T) {
	repo := new(mockShoppingListRepository)
	pantryRepo := new(mockPantryRepository)
	service := newService(repo, pantryRepo)

	userID := uuid.New()
	otherUser := uuid.New()
	listID := uuid.New()
	shoppingList := &shoppingModel.ShoppingList{ID: listID, UserID: otherUser}

	repo.On("GetByID", mock.Anything, listID).Return(shoppingList, nil).Once()

	result, err := service.GetShoppingListByID(context.Background(), userID, listID)
	require.ErrorIs(t, err, shoppingDomain.ErrUnauthorized)
	require.Nil(t, result)

	repo.AssertExpectations(t)
	pantryRepo.AssertExpectations(t)
}

func TestShoppingListService_UpdateShoppingListItem_ItemNotFound(t *testing.T) {
	repo := new(mockShoppingListRepository)
	pantryRepo := new(mockPantryRepository)
	service := newService(repo, pantryRepo)

	userID := uuid.New()
	listID := uuid.New()
	itemID := uuid.New()

	shoppingList := &shoppingModel.ShoppingList{ID: listID, UserID: userID, Items: []shoppingModel.ShoppingListItem{}}

	repo.On("GetByID", mock.Anything, listID).Return(shoppingList, nil).Once()

	input := dto.UpdateShoppingListItemDTO{Name: ptrString("Item")}

	result, err := service.UpdateShoppingListItem(context.Background(), userID, listID, itemID, input)
	require.ErrorIs(t, err, shoppingDomain.ErrItemNotFound)
	require.Nil(t, result)

	repo.AssertExpectations(t)
	pantryRepo.AssertExpectations(t)
}

func ptrString(value string) *string {
	return &value
}
