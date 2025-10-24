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
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"go.uber.org/zap"
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
	logger := appLogger.FromContext(ctx)

	// Check if profile already exists
	existingProfile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err == nil && existingProfile != nil {
		logger.Warn("Profile already exists for user",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "CreateProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
		)
		return nil, domain.ErrProfileAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to check existing profile",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "CreateProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
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
		logger.Error("Failed to create profile",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "CreateProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("create profile: %w", err)
	}

	logger.Info("Profile created successfully",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "CreateProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("profile_id", profile.ID.String()),
	)

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponseDTO, error) {
	logger := appLogger.FromContext(ctx)

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Profile not found for user",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "GetProfileByUserID"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			return nil, domain.ErrProfileNotFound
		}
		logger.Error("Failed to get profile",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "GetProfileByUserID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get profile: %w", err)
	}

	logger.Info("Profile retrieved successfully",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "GetProfileByUserID"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("profile_id", profile.ID.String()),
	)

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, input dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error) {
	logger := appLogger.FromContext(ctx)

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Profile not found for update",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "UpdateProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			return nil, domain.ErrProfileNotFound
		}
		logger.Error("Failed to get profile for update",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "UpdateProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
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
		logger.Error("Failed to update profile",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "UpdateProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("profile_id", profile.ID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("update profile: %w", err)
	}

	logger.Info("Profile updated successfully",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "UpdateProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("profile_id", profile.ID.String()),
	)

	return s.convertToResponseDTO(profile), nil
}

func (s *profileService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Profile not found for deletion",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "DeleteProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			return domain.ErrProfileNotFound
		}
		logger.Error("Failed to get profile for deletion",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "DeleteProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("get profile: %w", err)
	}

	if err := s.profileRepo.Delete(ctx, profile.ID); err != nil {
		logger.Error("Failed to delete profile",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "DeleteProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("profile_id", profile.ID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("delete profile: %w", err)
	}

	logger.Info("Profile deleted successfully",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "DeleteProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("profile_id", profile.ID.String()),
	)

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
