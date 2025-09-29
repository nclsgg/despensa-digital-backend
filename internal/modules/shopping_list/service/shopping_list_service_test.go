package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func (m *mockShoppingListRepository) Create(ctx context.Context, shoppingList *shoppingModel.ShoppingList) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "shoppingList": shoppingList}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.Create"), zap.Any("params", __logParams))
	args := m.Called(ctx, shoppingList)
	result0 = args.Error(0)
	return
}

func (m *mockShoppingListRepository) GetByID(ctx context.Context, id uuid.UUID) (result0 *shoppingModel.ShoppingList, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.GetByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.GetByID"), zap.Any("params", __logParams))
	args := m.Called(ctx, id)
	if list, ok := args.Get(0).(*shoppingModel.ShoppingList); ok {
		result0 = list
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockShoppingListRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) (result0 []*shoppingModel.ShoppingList, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "userID": userID, "limit": limit, "offset": offset}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.GetByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.GetByUserID"), zap.Any("params", __logParams))
	args := m.Called(ctx, userID, limit, offset)
	if lists, ok := args.Get(0).([]*shoppingModel.ShoppingList); ok {
		result0 = lists
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockShoppingListRepository) Update(ctx context.Context, shoppingList *shoppingModel.ShoppingList) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "shoppingList": shoppingList}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.Update"), zap.Any("params", __logParams))
	args := m.Called(ctx, shoppingList)
	result0 = args.Error(0)
	return
}

func (m *mockShoppingListRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.Delete"), zap.Any("params", __logParams))
	args := m.Called(ctx, id)
	result0 = args.Error(0)
	return
}

func (m *mockShoppingListRepository) CreateItem(ctx context.Context, item *shoppingModel.ShoppingListItem) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.CreateItem"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.CreateItem"), zap.Any("params", __logParams))
	args := m.Called(ctx, item)
	result0 = args.Error(0)
	return
}

func (m *mockShoppingListRepository) UpdateItem(ctx context.Context, item *shoppingModel.ShoppingListItem) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.UpdateItem"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.UpdateItem"), zap.Any("params", __logParams))
	args := m.Called(ctx, item)
	result0 = args.Error(0)
	return
}

func (m *mockShoppingListRepository) DeleteItem(ctx context.Context, itemID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "itemID": itemID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.DeleteItem"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.DeleteItem"), zap.Any("params", __logParams))
	args := m.Called(ctx, itemID)
	result0 = args.Error(0)
	return
}

func (m *mockShoppingListRepository) GetItemsByShoppingListID(ctx context.Context, shoppingListID uuid.UUID) (result0 []*shoppingModel.ShoppingListItem, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "shoppingListID": shoppingListID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.GetItemsByShoppingListID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.GetItemsByShoppingListID"), zap.Any("params", __logParams))
	args := m.Called(ctx, shoppingListID)
	if items, ok := args.Get(0).([]*shoppingModel.ShoppingListItem); ok {
		result0 = items
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockShoppingListRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (result0 int64, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockShoppingListRepository.CountByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockShoppingListRepository.CountByUserID"), zap.Any("params", __logParams))
	args := m.Called(ctx, userID)
	result0 = args.Get(0).(int64)
	result1 = args.Error(1)
	return
}

type mockPantryRepository struct {
	mock.Mock
}

func (m *mockPantryRepository) Create(ctx context.Context, pantry *pantryModel.Pantry) (result0 *pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantry": pantry}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.Create"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.Create"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantry)
	if p, ok := args.Get(0).(*pantryModel.Pantry); ok {
		result0 = p
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) Delete(ctx context.Context, pantryID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.Delete"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID)
	result0 = args.Error(0)
	return
}

func (m *mockPantryRepository) Update(ctx context.Context, pantry *pantryModel.Pantry) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantry": pantry}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.Update"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantry)
	result0 = args.Error(0)
	return
}

func (m *mockPantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (result0 *pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.GetByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.GetByID"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID)
	if pantry, ok := args.Get(0).(*pantryModel.Pantry); ok {
		result0 = pantry
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) (result0 []*pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.GetByUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.GetByUser"), zap.Any("params", __logParams))
	args := m.Called(ctx, userID)
	if pantries, ok := args.Get(0).([]*pantryModel.Pantry); ok {
		result0 = pantries
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 bool, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.IsUserInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.IsUserInPantry"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID, userID)
	result0 = args.Bool(0)
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (result0 bool, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.IsUserOwner"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.IsUserOwner"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID, userID)
	result0 = args.Bool(0)
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) AddUserToPantry(ctx context.Context, pantryUser *pantryModel.PantryUser) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryUser": pantryUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.AddUserToPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.AddUserToPantry"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryUser)
	result0 = args.Error(0)
	return
}

func (m *mockPantryRepository) RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.RemoveUserFromPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.RemoveUserFromPantry"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID, userID)
	result0 = args.Error(0)
	return
}

func (m *mockPantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) (result0 []*pantryModel.PantryUserInfo, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.ListUsersInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.ListUsersInPantry"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID)
	if users, ok := args.Get(0).([]*pantryModel.PantryUserInfo); ok {
		result0 = users
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func newService(repo *mockShoppingListRepository, pantryRepo *mockPantryRepository) (result0 shoppingDomain.ShoppingListService) {
	__logParams := map[string]any{"repo": repo, "pantryRepo": pantryRepo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "newService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "newService"), zap.Any("params", __logParams))
	result0 = service.NewShoppingListService(repo, pantryRepo, nil, nil, nil)
	return
}

func TestShoppingListService_CreateShoppingList_PantryAccessDenied(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestShoppingListService_CreateShoppingList_PantryAccessDenied"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestShoppingListService_CreateShoppingList_PantryAccessDenied"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestShoppingListService_GetShoppingListByID_NotFound"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestShoppingListService_GetShoppingListByID_NotFound"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestShoppingListService_GetShoppingListByID_Unauthorized"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestShoppingListService_GetShoppingListByID_Unauthorized"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestShoppingListService_UpdateShoppingListItem_ItemNotFound"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestShoppingListService_UpdateShoppingListItem_ItemNotFound"), zap.Any("params", __logParams))
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

func ptrString(value string) (result0 *string) {
	__logParams := map[string]any{"value": value}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "ptrString"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "ptrString"), zap.Any("params", __logParams))
	result0 = &value
	return
}
