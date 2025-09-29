package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) (result0 domain.ProfileRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewProfileRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewProfileRepository"), zap.Any("params", __logParams))
	result0 = &profileRepository{db: db}
	return
}

func (r *profileRepository) Create(ctx context.Context, profile *model.Profile) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "profile": profile}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileRepository.Create"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(profile).Error
	return
}

func (r *profileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (result0 *model.Profile, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileRepository.GetByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileRepository.GetByUserID"), zap.Any("params", __logParams))
	var profile model.Profile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*profileRepository.GetByUserID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &profile
	result1 = nil
	return
}

func (r *profileRepository) Update(ctx context.Context, profile *model.Profile) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "profile": profile}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileRepository.Update"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(profile).Error
	return
}

func (r *profileRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*profileRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*profileRepository.Delete"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Delete(&model.Profile{}, id).Error
	return
}
