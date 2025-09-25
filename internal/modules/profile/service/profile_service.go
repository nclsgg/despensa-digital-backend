package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	"gorm.io/gorm"
)

type profileService struct {
	profileRepo domain.ProfileRepository
}

func NewProfileService(profileRepo domain.ProfileRepository) domain.ProfileService {
	return &profileService{
		profileRepo: profileRepo,
	}
}

func (s *profileService) CreateProfile(ctx context.Context, userID uuid.UUID, input dto.CreateProfileDTO) (*dto.ProfileResponseDTO, error) {
	// Check if profile already exists
	existingProfile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err == nil && existingProfile != nil {
		return nil, fmt.Errorf("profile already exists for user")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error checking existing profile: %w", err)
	}

	profile := &model.Profile{
		UserID:              userID,
		MonthlyIncome:       input.MonthlyIncome,
		PreferredBudget:     input.PreferredBudget,
		HouseholdSize:       input.HouseholdSize,
		DietaryRestrictions: model.StringArray(input.DietaryRestrictions),
		PreferredBrands:     model.StringArray(input.PreferredBrands),
		ShoppingFrequency:   input.ShoppingFrequency,
	}

	if err := s.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("error creating profile: %w", err)
	}

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponseDTO, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("error getting profile: %w", err)
	}

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, input dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("error getting profile: %w", err)
	}

	// Update fields if provided
	if input.MonthlyIncome != nil {
		profile.MonthlyIncome = *input.MonthlyIncome
	}
	if input.PreferredBudget != nil {
		profile.PreferredBudget = *input.PreferredBudget
	}
	if input.HouseholdSize != nil {
		profile.HouseholdSize = *input.HouseholdSize
	}
	if input.DietaryRestrictions != nil {
		profile.DietaryRestrictions = model.StringArray(*input.DietaryRestrictions)
	}
	if input.PreferredBrands != nil {
		profile.PreferredBrands = model.StringArray(*input.PreferredBrands)
	}
	if input.ShoppingFrequency != nil {
		profile.ShoppingFrequency = *input.ShoppingFrequency
	}

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("error updating profile: %w", err)
	}

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("profile not found")
		}
		return fmt.Errorf("error getting profile: %w", err)
	}

	if err := s.profileRepo.Delete(ctx, profile.ID); err != nil {
		return fmt.Errorf("error deleting profile: %w", err)
	}

	return nil
}

func (s *profileService) convertToResponseDTO(profile *model.Profile) *dto.ProfileResponseDTO {
	return &dto.ProfileResponseDTO{
		ID:                  profile.ID,
		UserID:              profile.UserID,
		MonthlyIncome:       profile.MonthlyIncome,
		PreferredBudget:     profile.PreferredBudget,
		HouseholdSize:       profile.HouseholdSize,
		DietaryRestrictions: []string(profile.DietaryRestrictions),
		PreferredBrands:     []string(profile.PreferredBrands),
		ShoppingFrequency:   profile.ShoppingFrequency,
		CreatedAt:           profile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:           profile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
