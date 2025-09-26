package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/service"
)

type mockProfileRepository struct {
	mock.Mock
}

func (m *mockProfileRepository) Create(ctx context.Context, profile *model.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *mockProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	args := m.Called(ctx, userID)
	if profile, ok := args.Get(0).(*model.Profile); ok {
		return profile, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *mockProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func newProfileService(repo *mockProfileRepository) domain.ProfileService {
	return service.NewProfileService(repo)
}

func TestCreateProfile_ProfileAlreadyExists(t *testing.T) {
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
	repo := new(mockProfileRepository)
	svc := newProfileService(repo)
	userID := uuid.New()

	repo.On("GetByUserID", mock.Anything, userID).Return((*model.Profile)(nil), gorm.ErrRecordNotFound).Once()

	err := svc.DeleteProfile(context.Background(), userID)
	require.ErrorIs(t, err, domain.ErrProfileNotFound)

	repo.AssertExpectations(t)
}
