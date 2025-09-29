package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	itemDto "github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	itemModel "github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/service"
	userModel "github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockPantryRepository struct {
	mock.Mock
}

func (m *mockPantryRepository) Create(ctx context.Context, pantry *model.Pantry) (result0 *model.Pantry, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantry": pantry}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.Create"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.Create"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantry)
	result0 = args.Get(0).(*model.Pantry)
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) AddUserToPantry(ctx context.Context, pantryUser *model.PantryUser) (result0 error) {
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

func (m *mockPantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (result0 *model.Pantry, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.GetByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.GetByID"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID)
	result0 = args.Get(0).(*model.Pantry)
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) (result0 []*model.Pantry, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.GetByUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.GetByUser"), zap.Any("params", __logParams))
	args := m.Called(ctx, userID)
	result0 = args.Get(0).([]*model.Pantry)
	result1 = args.Error(1)
	return
}

func (m *mockPantryRepository) Update(ctx context.Context, pantry *model.Pantry) (result0 error) {
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

func (m *mockPantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) (result0 []*model.PantryUserInfo, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockPantryRepository.ListUsersInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockPantryRepository.ListUsersInPantry"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID)
	result0 = args.Get(0).([]*model.PantryUserInfo)
	result1 = args.Error(1)
	return
}

type mockUserRepository struct {
	mock.Mock
}

type mockItemRepository struct {
	mock.Mock
}

func (m *mockItemRepository) Create(ctx context.Context, item *itemModel.Item) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockItemRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockItemRepository.Create"), zap.Any("params", __logParams))
	args := m.Called(ctx, item)
	result0 = args.Error(0)
	return
}

func (m *mockItemRepository) Update(ctx context.Context, item *itemModel.Item) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockItemRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockItemRepository.Update"), zap.Any("params", __logParams))
	args := m.Called(ctx, item)
	result0 = args.Error(0)
	return
}

func (m *mockItemRepository) FindByID(ctx context.Context, id uuid.UUID) (result0 *itemModel.Item, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockItemRepository.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockItemRepository.FindByID"), zap.Any("params", __logParams))
	args := m.Called(ctx, id)
	result0 = args.Get(0).(*itemModel.Item)
	result1 = args.Error(1)
	return
}

func (m *mockItemRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockItemRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockItemRepository.Delete"), zap.Any("params", __logParams))
	args := m.Called(ctx, id)
	result0 = args.Error(0)
	return
}

func (m *mockItemRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 []*itemModel.Item, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockItemRepository.ListByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockItemRepository.ListByPantryID"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID)
	result0 = args.Get(0).([]*itemModel.Item)
	result1 = args.Error(1)
	return
}

func (m *mockItemRepository) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters itemDto.ItemFilterDTO) (result0 []*itemModel.Item, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID, "filters": filters}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockItemRepository.FilterByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockItemRepository.FilterByPantryID"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID, filters)
	result0 = args.Get(0).([]*itemModel.Item)
	result1 = args.Error(1)
	return
}

func (m *mockItemRepository) CountByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 int, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockItemRepository.CountByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockItemRepository.CountByPantryID"), zap.Any("params", __logParams))
	args := m.Called(ctx, pantryID)
	result0 = args.Int(0)
	result1 = args.Error(1)
	return
}

func (m *mockUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (result0 *userModel.User, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockUserRepository.GetUserById"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockUserRepository.GetUserById"), zap.Any("params", __logParams))
	args := m.Called(ctx, id)
	if usr, ok := args.Get(0).(*userModel.User); ok {
		result0 = usr
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (result0 *userModel.User, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "email": email}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockUserRepository.GetUserByEmail"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockUserRepository.GetUserByEmail"), zap.Any("params", __logParams))
	args := m.Called(ctx, email)
	if usr, ok := args.Get(0).(*userModel.User); ok {
		result0 = usr
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) (result0 []userModel.User, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockUserRepository.GetAllUsers"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockUserRepository.GetAllUsers"), zap.Any("params", __logParams))
	args := m.Called(ctx)
	if usrs, ok := args.Get(0).([]userModel.User); ok {
		result0 = usrs
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *userModel.User) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "user": user}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockUserRepository.UpdateUser"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockUserRepository.UpdateUser"), zap.Any("params", __logParams))
	args := m.Called(ctx, user)
	result0 = args.Error(0)
	return
}

func TestCreatePantry(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestCreatePantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestCreatePantry"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestAddUserToPantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestAddUserToPantry"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

	ctx := context.Background()
	pantryID := uuid.New()
	ownerID := uuid.New()
	targetUserID := uuid.New()
	targetUser := "teste@email.com"

	repo.On("IsUserOwner", ctx, pantryID, ownerID).Return(true, nil)
	repo.On("IsUserInPantry", ctx, pantryID, targetUserID).Return(false, nil)
	userRepo.On("GetUserByEmail", ctx, targetUser).Return(&userModel.User{ID: targetUserID, Email: targetUser}, nil)
	repo.On("AddUserToPantry", ctx, mock.AnythingOfType("*model.PantryUser")).Return(nil)

	err := svc.AddUserToPantry(ctx, pantryID, ownerID, targetUser)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRemoveUserFromPantry_Success(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestRemoveUserFromPantry_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestRemoveUserFromPantry_Success"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

	ctx := context.Background()
	pantryID := uuid.New()
	ownerID := uuid.New()
	targetUserID := uuid.New()
	targetUser := "teste@email.com"

	repo.On("IsUserOwner", ctx, pantryID, ownerID).Return(true, nil)
	repo.On("RemoveUserFromPantry", ctx, pantryID, targetUserID).Return(nil)
	userRepo.On("GetUserByEmail", ctx, targetUser).Return(&userModel.User{ID: targetUserID, Email: targetUser}, nil)

	err := svc.RemoveUserFromPantry(ctx, pantryID, ownerID, targetUser)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRemoveUserFromPantry_BlockOwner(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestRemoveUserFromPantry_BlockOwner"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestRemoveUserFromPantry_BlockOwner"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

	ctx := context.Background()
	pantryID := uuid.New()
	ownerID := uuid.New()
	ownerEmail := "teste@email.com"

	repo.On("IsUserOwner", ctx, pantryID, ownerID).Return(true, nil)
	userRepo.On("GetUserByEmail", ctx, ownerEmail).Return(&userModel.User{ID: ownerID, Email: ownerEmail}, nil)

	err := svc.RemoveUserFromPantry(ctx, pantryID, ownerID, ownerEmail)
	assert.EqualError(t, err, "owner cannot remove themselves")
	repo.AssertExpectations(t)
}

func TestGetPantry_Success(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetPantry_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetPantry_Success"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetPantry_NotMember"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetPantry_NotMember"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	repo.On("IsUserInPantry", ctx, pantryID, userID).Return(false, nil)

	_, err := svc.GetPantry(ctx, pantryID, userID)
	assert.ErrorIs(t, err, service.ErrUnauthorized)
	repo.AssertExpectations(t)
}

func TestListUsersInPantry_Success(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestListUsersInPantry_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestListUsersInPantry_Success"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	userInfo1 := &model.PantryUserInfo{ID: uuid.New(), PantryID: pantryID, UserID: uuid.New(), Email: "test1@example.com", Role: "member"}
	userInfo2 := &model.PantryUserInfo{ID: uuid.New(), PantryID: pantryID, UserID: uuid.New(), Email: "test2@example.com", Role: "owner"}
	expectedUsers := []*model.PantryUserInfo{userInfo1, userInfo2}

	repo.On("IsUserInPantry", ctx, pantryID, userID).Return(true, nil)
	repo.On("ListUsersInPantry", ctx, pantryID).Return(expectedUsers, nil)
	userRepo.On("GetUserById", ctx, userInfo1.UserID).Return(&userModel.User{ID: userInfo1.UserID, FirstName: "Test", LastName: "One"}, nil)
	userRepo.On("GetUserById", ctx, userInfo2.UserID).Return(&userModel.User{ID: userInfo2.UserID, FirstName: "Admin", LastName: "User"}, nil)

	result, err := svc.ListUsersInPantry(ctx, pantryID, userID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestListUsersInPantry_Unauthorized(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestListUsersInPantry_Unauthorized"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestListUsersInPantry_Unauthorized"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	repo.On("IsUserInPantry", ctx, pantryID, userID).Return(false, nil)

	_, err := svc.ListUsersInPantry(ctx, pantryID, userID)
	assert.EqualError(t, err, "user is not in the pantry")
	repo.AssertExpectations(t)
}

func TestUpdatePantry_Success(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestUpdatePantry_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestUpdatePantry_Success"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestUpdatePantry_Unauthorized"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestUpdatePantry_Unauthorized"), zap.Any("params", __logParams))
	repo := new(mockPantryRepository)
	userRepo := new(mockUserRepository)
	itemRepo := new(mockItemRepository)
	svc := service.NewPantryService(repo, userRepo, itemRepo)

	ctx := context.Background()
	pantryID := uuid.New()
	userID := uuid.New()

	repo.On("IsUserOwner", ctx, pantryID, userID).Return(false, nil)

	err := svc.UpdatePantry(ctx, pantryID, userID, "Name")
	assert.EqualError(t, err, service.ErrUnauthorized.Error())
	repo.AssertExpectations(t)
}
