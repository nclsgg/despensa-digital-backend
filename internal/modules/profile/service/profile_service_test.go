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

	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/service"
)

type mockProfileRepository struct {
	mock.Mock
}

func (m *mockProfileRepository) Create(ctx context.Context, profile *model.Profile) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "profile": profile}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockProfileRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockProfileRepository.Create"), zap.Any("params", __logParams))
	args := m.Called(ctx, profile)
	result0 = args.Error(0)
	return
}

func (m *mockProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (result0 *model.Profile, result1 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockProfileRepository.GetByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockProfileRepository.GetByUserID"), zap.Any("params", __logParams))
	args := m.Called(ctx, userID)
	if profile, ok := args.Get(0).(*model.Profile); ok {
		result0 = profile
		result1 = args.Error(1)
		return
	}
	result0 = nil
	result1 = args.Error(1)
	return
}

func (m *mockProfileRepository) Update(ctx context.Context, profile *model.Profile) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "profile": profile}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockProfileRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockProfileRepository.Update"), zap.Any("params", __logParams))
	args := m.Called(ctx, profile)
	result0 = args.Error(0)
	return
}

func (m *mockProfileRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"m": m, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*mockProfileRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*mockProfileRepository.Delete"), zap.Any("params", __logParams))
	args := m.Called(ctx, id)
	result0 = args.Error(0)
	return
}

func newProfileService(repo *mockProfileRepository) (result0 domain.ProfileService) {
	__logParams := map[string]any{"repo": repo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "newProfileService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "newProfileService"), zap.Any("params", __logParams))
	result0 = service.NewProfileService(repo)
	return
}

func TestCreateProfile_ProfileAlreadyExists(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestCreateProfile_ProfileAlreadyExists"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestCreateProfile_ProfileAlreadyExists"), zap.Any("params", __logParams))
	repo := new(mockProfileRepository)
	svc := newProfileService(repo)
	userID := uuid.New()

	repo.On("GetByUserID", mock.Anything, userID).Return(&model.Profile{ID: uuid.New()}, nil).Once()

	result, err := svc.CreateProfile(context.Background(), userID, dto.CreateProfileDTO{})
	require.ErrorIs(t, err, domain.ErrProfileAlreadyExists)
	require.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestCreateProfile_Success(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestCreateProfile_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestCreateProfile_Success"), zap.Any("params", __logParams))
	repo := new(mockProfileRepository)
	svc := newProfileService(repo)
	userID := uuid.New()

	repo.On("GetByUserID", mock.Anything, userID).Return((*model.Profile)(nil), gorm.ErrRecordNotFound).Once()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Profile")).Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*model.Profile)
		profile.ID = uuid.New()
		now := time.Now().UTC()
		profile.CreatedAt = now
		profile.UpdatedAt = now
	})

	input := dto.CreateProfileDTO{MonthlyIncome: 1000, PreferredBudget: 200, HouseholdSize: 2}

	result, err := svc.CreateProfile(context.Background(), userID, input)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, userID.String(), result.UserID)
	require.Equal(t, input.MonthlyIncome, result.MonthlyIncome)

	repo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGetProfile_NotFound"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGetProfile_NotFound"), zap.Any("params", __logParams))
	repo := new(mockProfileRepository)
	svc := newProfileService(repo)
	userID := uuid.New()

	repo.On("GetByUserID", mock.Anything, userID).Return((*model.Profile)(nil), gorm.ErrRecordNotFound).Once()

	result, err := svc.GetProfileByUserID(context.Background(), userID)
	require.ErrorIs(t, err, domain.ErrProfileNotFound)
	require.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestUpdateProfile_NotFound(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestUpdateProfile_NotFound"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestUpdateProfile_NotFound"), zap.Any("params", __logParams))
	repo := new(mockProfileRepository)
	svc := newProfileService(repo)
	userID := uuid.New()

	repo.On("GetByUserID", mock.Anything, userID).Return((*model.Profile)(nil), gorm.ErrRecordNotFound).Once()

	result, err := svc.UpdateProfile(context.Background(), userID, dto.UpdateProfileDTO{})
	require.ErrorIs(t, err, domain.ErrProfileNotFound)
	require.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestDeleteProfile_NotFound(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestDeleteProfile_NotFound"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestDeleteProfile_NotFound"), zap.Any("params", __logParams))
	repo := new(mockProfileRepository)
	svc := newProfileService(repo)
	userID := uuid.New()

	repo.On("GetByUserID", mock.Anything, userID).Return((*model.Profile)(nil), gorm.ErrRecordNotFound).Once()

	err := svc.DeleteProfile(context.Background(), userID)
	require.ErrorIs(t, err, domain.ErrProfileNotFound)

	repo.AssertExpectations(t)
}
