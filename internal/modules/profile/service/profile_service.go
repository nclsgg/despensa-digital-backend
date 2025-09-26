package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
		return nil, domain.ErrProfileAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check existing profile: %w", err)
	}

	profile := &model.Profile{
		UserID:              userID,
		MonthlyIncome:       input.MonthlyIncome,
		PreferredBudget:     input.PreferredBudget,
		HouseholdSize:       input.HouseholdSize,
		DietaryRestrictions: model.StringArray(normalizeStringSlice(input.DietaryRestrictions)),
		PreferredBrands:     model.StringArray(normalizeStringSlice(input.PreferredBrands)),
		ShoppingFrequency:   input.ShoppingFrequency,
	}

	if err := s.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("create profile: %w", err)
	}

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponseDTO, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, input dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
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
		profile.DietaryRestrictions = model.StringArray(normalizeStringSlice(*input.DietaryRestrictions))
	}
	if input.PreferredBrands != nil {
		profile.PreferredBrands = model.StringArray(normalizeStringSlice(*input.PreferredBrands))
	}
	if input.ShoppingFrequency != nil {
		profile.ShoppingFrequency = *input.ShoppingFrequency
	}

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrProfileNotFound
		}
		return fmt.Errorf("get profile: %w", err)
	}

	if err := s.profileRepo.Delete(ctx, profile.ID); err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	return nil
}

func (s *profileService) convertToResponseDTO(profile *model.Profile) *dto.ProfileResponseDTO {
	return &dto.ProfileResponseDTO{
		ID:                  profile.ID.String(),
		UserID:              profile.UserID.String(),
		MonthlyIncome:       profile.MonthlyIncome,
		PreferredBudget:     profile.PreferredBudget,
		HouseholdSize:       profile.HouseholdSize,
		DietaryRestrictions: toStringSlice(profile.DietaryRestrictions),
		PreferredBrands:     toStringSlice(profile.PreferredBrands),
		ShoppingFrequency:   profile.ShoppingFrequency,
		CreatedAt:           profile.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:           profile.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func normalizeStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func toStringSlice(values model.StringArray) []string {
	if values == nil {
		return []string{}
	}
	return append([]string(nil), values...)
}
