package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/service"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockUserRepository.GetUserById"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockUserRepository.GetUserById"), zap.Any("params", __logParams))
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*model.User); ok {
		result0 = user
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "email": email}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockUserRepository.GetUserByEmail"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockUserRepository.GetUserByEmail"), zap.Any("params", __logParams))
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*model.User); ok {
		result0 = user
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) (result0 []model.User, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockUserRepository.GetAllUsers"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockUserRepository.GetAllUsers"), zap.Any("params", __logParams))
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]model.User); ok {
		result0 = users
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *model.User) (result0 error) {
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

func newUserService(repo *mockUserRepository) (result0 domain.UserService) {
	__logParams := map[string]any{"repo": repo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "newUserService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "newUserService"), zap.Any("params", __logParams))
	result0 = service.NewUserService(repo)
	return
}

func TestGetUserByID_NotFound(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetUserByID_NotFound"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetUserByID_NotFound"), zap.Any("params", __logParams))
	repo := new(mockUserRepository)
	svc := newUserService(repo)
	userID := uuid.New()

	repo.On("GetUserById", mock.Anything, userID).Return((*model.User)(nil), gorm.ErrRecordNotFound).Once()

	result, err := svc.GetUserById(context.Background(), userID)
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	require.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetUserByID_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetUserByID_Success"), zap.Any("params", __logParams))
	repo := new(mockUserRepository)
	svc := newUserService(repo)
	userID := uuid.New()
	user := &model.User{ID: userID}

	repo.On("GetUserById", mock.Anything, userID).Return(user, nil).Once()

	result, err := svc.GetUserById(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, user, result)

	repo.AssertExpectations(t)
}

func TestGetAllUsers_Error(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetAllUsers_Error"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetAllUsers_Error"), zap.Any("params", __logParams))
	repo := new(mockUserRepository)
	svc := newUserService(repo)

	repo.On("GetAllUsers", mock.Anything).Return(nil, errors.New("repository failure")).Once()

	result, err := svc.GetAllUsers(context.Background())
	require.Error(t, err)
	require.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestCompleteProfile_UserNotFound(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestCompleteProfile_UserNotFound"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestCompleteProfile_UserNotFound"), zap.Any("params", __logParams))
	repo := new(mockUserRepository)
	svc := newUserService(repo)
	userID := uuid.New()

	repo.On("GetUserById", mock.Anything, userID).Return((*model.User)(nil), gorm.ErrRecordNotFound).Once()

	err := svc.CompleteProfile(context.Background(), userID, "Ana", "Silva")
	require.ErrorIs(t, err, domain.ErrUserNotFound)

	repo.AssertExpectations(t)
}

func TestCompleteProfile_Success(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestCompleteProfile_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestCompleteProfile_Success"), zap.Any("params", __logParams))
	repo := new(mockUserRepository)
	svc := newUserService(repo)
	userID := uuid.New()
	user := &model.User{ID: userID}

	repo.On("GetUserById", mock.Anything, userID).Return(user, nil).Once()
	repo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
		updated := args.Get(1).(*model.User)
		require.True(t, updated.ProfileCompleted)
		require.Equal(t, "Ana", updated.FirstName)
		require.Equal(t, "Silva", updated.LastName)
	})

	require.NoError(t, svc.CompleteProfile(context.Background(), userID, "Ana", "Silva"))

	repo.AssertExpectations(t)
}
