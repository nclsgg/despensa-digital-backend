package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/service"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*model.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*model.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]model.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func newUserService(repo *mockUserRepository) domain.UserService {
	return service.NewUserService(repo)
}

func TestGetUserByID_NotFound(t *testing.T) {
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
	repo := new(mockUserRepository)
	svc := newUserService(repo)

	repo.On("GetAllUsers", mock.Anything).Return(nil, errors.New("repository failure")).Once()

	result, err := svc.GetAllUsers(context.Background())
	require.Error(t, err)
	require.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestCompleteProfile_UserNotFound(t *testing.T) {
	repo := new(mockUserRepository)
	svc := newUserService(repo)
	userID := uuid.New()

	repo.On("GetUserById", mock.Anything, userID).Return((*model.User)(nil), gorm.ErrRecordNotFound).Once()

	err := svc.CompleteProfile(context.Background(), userID, "Ana", "Silva")
	require.ErrorIs(t, err, domain.ErrUserNotFound)

	repo.AssertExpectations(t)
}

func TestCompleteProfile_Success(t *testing.T) {
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
