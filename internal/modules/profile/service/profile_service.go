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
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type profileService struct {
	profileRepo domain.ProfileRepository
}

func NewProfileService(profileRepo domain.ProfileRepository) (result0 domain.ProfileService) {
	__logParams := map[string]any{"profileRepo": profileRepo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewProfileService"), zap.Any("result", result0), zap.Duration("duration", time.

			// Check if profile already exists
			Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewProfileService"), zap.Any("params", __logParams))
	result0 = &profileService{
		profileRepo: profileRepo,
	}
	return
}

func (s *profileService) CreateProfile(ctx context.Context, userID uuid.UUID, input dto.CreateProfileDTO) (result0 *dto.ProfileResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileService.CreateProfile"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileService.CreateProfile"), zap.Any("params", __logParams))

	existingProfile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err == nil && existingProfile != nil {
		result0 = nil
		result1 = domain.ErrProfileAlreadyExists
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		result0 = nil
		result1 = fmt.Errorf("check existing profile: %w", err)
		return
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
		zap.L().Error("function.error", zap.String("func", "*profileService.CreateProfile"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("create profile: %w", err)
		return
	}
	result0 = s.convertToResponseDTO(profile)
	result1 = nil
	return
}

func (s *profileService) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (result0 *dto.ProfileResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileService.GetProfileByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileService.GetProfileByUserID"), zap.Any("params", __logParams))
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*profileService.GetProfileByUserID"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrProfileNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get profile: %w", err)
		return
	}
	result0 = s.convertToResponseDTO(profile)
	result1 = nil
	return
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, input dto.UpdateProfileDTO) (result0 *dto.ProfileResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileService.UpdateProfile"), zap.Any("result", map[string]any{"result0":

		// Update fields if provided
		result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileService.UpdateProfile"), zap.Any("params", __logParams))
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*profileService.UpdateProfile"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrProfileNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get profile: %w", err)
		return
	}

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
		zap.L().Error("function.error", zap.String("func", "*profileService.UpdateProfile"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("update profile: %w", err)
		return
	}
	result0 = s.convertToResponseDTO(profile)
	result1 = nil
	return
}

func (s *profileService) DeleteProfile(ctx context.Context, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileService.DeleteProfile"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileService.DeleteProfile"), zap.Any("params", __logParams))
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*profileService.DeleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = domain.ErrProfileNotFound
			return
		}
		result0 = fmt.Errorf("get profile: %w", err)
		return
	}

	if err := s.profileRepo.Delete(ctx, profile.ID); err != nil {
		zap.L().Error("function.error", zap.String("func", "*profileService.DeleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		result0 = fmt.Errorf("delete profile: %w", err)
		return
	}
	result0 = nil
	return
}

func (s *profileService) convertToResponseDTO(profile *model.Profile) (result0 *dto.ProfileResponseDTO) {
	__logParams := map[string]any{"s": s, "profile": profile}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileService.convertToResponseDTO"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileService.convertToResponseDTO"), zap.Any("params", __logParams))
	result0 = &dto.ProfileResponseDTO{
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
	return
}

func normalizeStringSlice(values []string) (result0 []string) {
	__logParams := map[string]any{"values": values}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "normalizeStringSlice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "normalizeStringSlice"), zap.Any("params", __logParams))
	if values == nil {
		result0 = []string{}
		return
	}
	result0 = values
	return
}

func toStringSlice(values model.StringArray) (result0 []string) {
	__logParams := map[string]any{"values": values}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "toStringSlice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "toStringSlice"), zap.Any("params", __logParams))
	if values == nil {
		result0 = []string{}
		return
	}
	result0 = append([]string(nil), values...)
	return
}
