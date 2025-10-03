package service_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"

	llmDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	profileModel "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
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

func (m *mockPantryRepository) UpdatePantryUserRole(ctx context.Context, pantryID, userID uuid.UUID, newRole string) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID, "userID": userID, "newRole": newRole}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.UpdatePantryUserRole"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.UpdatePantryUserRole"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID, userID, newRole)
	result0 = args.Error(0)
	return
}

func (m *mockPantryRepository) GetPantryUser(ctx context.Context, pantryID, userID uuid.UUID) (result0 *pantryModel.PantryUser, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.GetPantryUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.GetPantryUser"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID, userID)
	if pu, ok := args.Get(0).(*pantryModel.PantryUser); ok {
		result0 = pu
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
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

type mockProfileRepository struct {
	mock.Mock
}

func (m *mockProfileRepository) Create(ctx context.Context, profile *profileModel.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *mockProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*profileModel.Profile, error) {
	args := m.Called(ctx, userID)
	if profile, ok := args.Get(0).(*profileModel.Profile); ok {
		return profile, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProfileRepository) Update(ctx context.Context, profile *profileModel.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *mockProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type fakeLLMService struct {
	response *llmDTO.LLMResponseDTO
}

func (f *fakeLLMService) ProcessRequest(ctx context.Context, request *llmDTO.LLMRequestDTO) (*llmDTO.LLMResponseDTO, error) {
	return nil, nil
}

func (f *fakeLLMService) GenerateText(ctx context.Context, prompt string, options map[string]interface{}) (*llmDTO.LLMResponseDTO, error) {
	return f.response, nil
}

func (f *fakeLLMService) BuildPrompt(ctx context.Context, templateID string, variables map[string]string) (string, error) {
	return "", nil
}

func (f *fakeLLMService) GetAvailableProviders() []string {
	return nil
}

func (f *fakeLLMService) SetProvider(providerName string) error {
	return nil
}

func (f *fakeLLMService) GetCurrentProvider() string {
	return ""
}

func newService(repo *mockShoppingListRepository, pantryRepo *mockPantryRepository) (result0 shoppingDomain.ShoppingListService) {
	__logParams := map[string]any{"repo": repo, "pantryRepo": pantryRepo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "newService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "newService"), zap.Any("params", __logParams))
	profileRepo := new(mockProfileRepository)
	profileRepo.On("GetByUserID", mock.Anything, mock.Anything).Return((*profileModel.Profile)(nil), gorm.ErrRecordNotFound).Maybe()
	result0 = service.NewShoppingListService(repo, pantryRepo, nil, profileRepo, nil)
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

func TestShoppingListService_UpdateShoppingListItem_RecalculateTotals(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestShoppingListService_UpdateShoppingListItem_RecalculateTotals"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestShoppingListService_UpdateShoppingListItem_RecalculateTotals"), zap.Any("params", __logParams))
	repo := new(mockShoppingListRepository)
	pantryRepo := new(mockPantryRepository)
	service := newService(repo, pantryRepo)

	userID := uuid.New()
	listID := uuid.New()
	itemID := uuid.New()
	secondItemID := uuid.New()

	shoppingList := &shoppingModel.ShoppingList{
		ID:     listID,
		UserID: userID,
		Items: []shoppingModel.ShoppingListItem{
			{
				ID:             itemID,
				ShoppingListID: listID,
				Name:           "Arroz",
				Quantity:       1,
				Unit:           "kg",
				EstimatedPrice: 10,
			},
			{
				ID:             secondItemID,
				ShoppingListID: listID,
				Name:           "Feijao",
				Quantity:       2,
				Unit:           "kg",
				EstimatedPrice: 5,
				ActualPrice:    4,
				Purchased:      true,
			},
		},
	}

	repo.On("GetByID", mock.Anything, listID).Return(shoppingList, nil).Once()
	repo.On("UpdateItem", mock.Anything, mock.MatchedBy(func(item *shoppingModel.ShoppingListItem) bool {
		return item.ID == itemID && math.Abs(item.Quantity-3) < 1e-6 && math.Abs(item.EstimatedPrice-12) < 1e-6
	})).Return(nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(list *shoppingModel.ShoppingList) bool {
		return list.ID == listID && math.Abs(list.EstimatedCost-46) < 1e-6 && math.Abs(list.ActualCost-8) < 1e-6
	})).Return(nil).Once()
	repo.On("GetItemsByShoppingListID", mock.Anything, listID).Return([]*shoppingModel.ShoppingListItem{
		{
			ID:             itemID,
			ShoppingListID: listID,
			Name:           "Arroz",
			Quantity:       3,
			Unit:           "kg",
			EstimatedPrice: 12,
		},
		{
			ID:             secondItemID,
			ShoppingListID: listID,
			Name:           "Feijao",
			Quantity:       2,
			Unit:           "kg",
			EstimatedPrice: 5,
			ActualPrice:    4,
			Purchased:      true,
		},
	}, nil).Once()

	input := dto.UpdateShoppingListItemDTO{
		Quantity:       ptrFloat64(3),
		EstimatedPrice: ptrFloat64(12),
	}

	result, err := service.UpdateShoppingListItem(context.Background(), userID, listID, itemID, input)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.InEpsilon(t, 3, result.Quantity, 1e-6)
	require.InEpsilon(t, 12, result.EstimatedPrice, 1e-6)

	repo.AssertExpectations(t)
	pantryRepo.AssertExpectations(t)
}

func TestShoppingListService_UpdateShoppingList_FinalizePurchasedOnly(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestShoppingListService_UpdateShoppingList_FinalizePurchasedOnly"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestShoppingListService_UpdateShoppingList_FinalizePurchasedOnly"), zap.Any("params", __logParams))
	repo := new(mockShoppingListRepository)
	pantryRepo := new(mockPantryRepository)
	service := newService(repo, pantryRepo)

	userID := uuid.New()
	listID := uuid.New()
	purchasedID := uuid.New()
	pendingID := uuid.New()

	shoppingList := &shoppingModel.ShoppingList{
		ID:     listID,
		UserID: userID,
		Items: []shoppingModel.ShoppingListItem{
			{
				ID:             purchasedID,
				ShoppingListID: listID,
				Name:           "Cafe",
				Quantity:       1,
				Unit:           "unidade",
				EstimatedPrice: 15,
				ActualPrice:    12,
				Purchased:      true,
			},
			{
				ID:             pendingID,
				ShoppingListID: listID,
				Name:           "Açucar",
				Quantity:       2,
				Unit:           "kg",
				EstimatedPrice: 6,
				Purchased:      false,
			},
		},
	}

	updatedList := &shoppingModel.ShoppingList{
		ID:            listID,
		UserID:        userID,
		Status:        "completed",
		EstimatedCost: shoppingList.EstimatedCost,
		ActualCost:    12,
		Items: []shoppingModel.ShoppingListItem{
			shoppingList.Items[0],
			shoppingList.Items[1],
		},
	}

	repo.On("GetByID", mock.Anything, listID).Return(shoppingList, nil).Once()
	repo.On("UpdateItem", mock.Anything, mock.MatchedBy(func(item *shoppingModel.ShoppingListItem) bool {
		return item.ID == purchasedID && item.Purchased
	})).Return(nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(list *shoppingModel.ShoppingList) bool {
		return list.Status == "completed" && math.Abs(list.ActualCost-12) < 1e-6
	})).Return(nil).Once()
	repo.On("GetByID", mock.Anything, listID).Return(updatedList, nil).Once()

	statusCompleted := "completed"
	result, err := service.UpdateShoppingList(context.Background(), userID, listID, dto.UpdateShoppingListDTO{Status: &statusCompleted})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.InEpsilon(t, 12, result.ActualCost, 1e-6)
	require.Equal(t, "completed", result.Status)
	require.Len(t, result.Items, 2)
	require.False(t, result.Items[1].Purchased)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.MatchedBy(func(item *shoppingModel.ShoppingListItem) bool {
		return item.ID == pendingID
	}))
	pantryRepo.AssertExpectations(t)
}

func TestShoppingListService_GenerateAIShoppingList_PopulatesItems(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestShoppingListService_GenerateAIShoppingList_PopulatesItems"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestShoppingListService_GenerateAIShoppingList_PopulatesItems"), zap.Any("params", __logParams))
	repo := new(mockShoppingListRepository)
	pantryRepo := new(mockPantryRepository)
	profileRepo := new(mockProfileRepository)

	userID := uuid.New()
	pantryID := uuid.New()
	pantry := &pantryModel.Pantry{ID: pantryID, Name: "Casa"}

	profileRepo.On("GetByUserID", mock.Anything, userID).Return((*profileModel.Profile)(nil), gorm.ErrRecordNotFound).Maybe()
	pantryRepo.On("GetByID", mock.Anything, pantryID).Return(pantry, nil).Maybe()
	pantryRepo.On("IsUserInPantry", mock.Anything, pantryID, userID).Return(true, nil).Once()

	aiResponse := `{"items":[{"name":"Arroz","quantity":2,"unit":"kg","estimated_price":30,"category":"Grãos","priority":1,"reason":"Reposição"},{"name":"Feijao","quantity":3,"unit":"un","estimated_price":15,"category":"Grãos","priority":2,"reason":"Consumo semanal"}],"reasoning":"Lista gerada para teste","estimated_total":45}`
	llmStub := &fakeLLMService{
		response: &llmDTO.LLMResponseDTO{Response: aiResponse},
	}
	service := service.NewShoppingListService(repo, pantryRepo, nil, profileRepo, llmStub)

	var capturedList *shoppingModel.ShoppingList
	repo.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		list := args.Get(1).(*shoppingModel.ShoppingList)
		require.Equal(t, userID, list.UserID)
		require.Equal(t, "ai", list.GeneratedBy)
		require.Len(t, list.Items, 2)

		itemsByName := make(map[string]shoppingModel.ShoppingListItem)
		for _, item := range list.Items {
			itemsByName[item.Name] = item
		}

		arroz, ok := itemsByName["Arroz"]
		require.True(t, ok)
		require.InEpsilon(t, 2, arroz.Quantity, 1e-6)
		require.InEpsilon(t, 15, arroz.EstimatedPrice, 1e-6)

		feijao, ok := itemsByName["Feijao"]
		require.True(t, ok)
		require.InEpsilon(t, 3, feijao.Quantity, 1e-6)
		require.InEpsilon(t, 5, feijao.EstimatedPrice, 1e-6)

		require.InEpsilon(t, 45, list.EstimatedCost, 1e-6)

		if list.ID == uuid.Nil {
			list.ID = uuid.New()
		}

		persisted := &shoppingModel.ShoppingList{
			ID:                  list.ID,
			UserID:              list.UserID,
			PantryID:            list.PantryID,
			Name:                list.Name,
			Status:              list.Status,
			TotalBudget:         list.TotalBudget,
			EstimatedCost:       list.EstimatedCost,
			ActualCost:          list.ActualCost,
			GeneratedBy:         list.GeneratedBy,
			HouseholdSize:       list.HouseholdSize,
			MonthlyIncome:       list.MonthlyIncome,
			DietaryRestrictions: list.DietaryRestrictions,
			Items:               make([]shoppingModel.ShoppingListItem, len(list.Items)),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		for idx := range list.Items {
			item := list.Items[idx]
			if item.ID == uuid.Nil {
				item.ID = uuid.New()
			}
			item.ShoppingListID = list.ID
			item.CreatedAt = time.Now()
			item.UpdatedAt = time.Now()
			persisted.Items[idx] = item
		}

		capturedList = persisted
		repo.On("GetByID", mock.Anything, list.ID).Return(capturedList, nil).Once()
	}).Once()

	input := dto.GenerateAIShoppingListDTO{
		Name:     "Lista Inteligente",
		PantryID: pantryID,
	}
	result, err := service.GenerateAIShoppingList(context.Background(), userID, input)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "ai", result.GeneratedBy)
	require.InEpsilon(t, 45, result.EstimatedCost, 1e-6)
	require.Len(t, result.Items, 2)

	var arrozDTO, feijaoDTO *dto.ShoppingListItemResponseDTO
	for idx := range result.Items {
		item := &result.Items[idx]
		switch item.Name {
		case "Arroz":
			arrozDTO = item
		case "Feijao":
			feijaoDTO = item
		}
	}
	require.NotNil(t, arrozDTO)
	require.NotNil(t, feijaoDTO)
	require.InEpsilon(t, 2, arrozDTO.Quantity, 1e-6)
	require.InEpsilon(t, 15, arrozDTO.EstimatedPrice, 1e-6)
	require.InEpsilon(t, 3, feijaoDTO.Quantity, 1e-6)
	require.InEpsilon(t, 5, feijaoDTO.EstimatedPrice, 1e-6)

	repo.AssertExpectations(t)
	pantryRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
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

func ptrFloat64(value float64) (result0 *float64) {
	__logParams := map[string]any{"value": value}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "ptrFloat64"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "ptrFloat64"), zap.Any("params", __logParams))
	result0 = &value
	return
}
