package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type pantryRepository struct {
	db *gorm.DB
}

func NewPantryRepository(db *gorm.DB) (result0 domain.PantryRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewPantryRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewPantryRepository"), zap.Any("params", __logParams))
	result0 = &pantryRepository{db}
	return
}

func (r *pantryRepository) Create(ctx context.Context, pantry *model.Pantry) (result0 *model.Pantry, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantry": pantry}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.Create"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.Create"), zap.Any("params", __logParams))
	err := r.db.WithContext(ctx).Create(pantry).Error
	result0 = pantry
	result1 = err
	return
}

func (r *pantryRepository) Delete(ctx context.Context, pantryID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.Delete"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Delete(&model.Pantry{}, "id = ?", pantryID).Error
	return
}

func (r *pantryRepository) Update(ctx context.Context, pantry *model.Pantry) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantry": pantry}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.Update"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(pantry).Error
	return
}

func (r *pantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (result0 *model.Pantry, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.GetByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.GetByID"), zap.Any("params", __logParams))
	var pantry model.Pantry
	err := r.db.WithContext(ctx).First(&pantry, "id = ?", pantryID).Error
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryRepository.GetByID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &pantry
	result1 = nil
	return
}

func (r *pantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) (result0 []*model.Pantry, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.GetByUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.GetByUser"), zap.Any("params", __logParams))
	var pantries []*model.Pantry
	err := r.db.WithContext(ctx).
		Joins("JOIN pantry_users ON pantries.id = pantry_users.pantry_id").
		Where("pantry_users.user_id = ? AND pantry_users.deleted_at IS NULL", userID).
		Find(&pantries).Error
	result0 = pantries
	result1 = err
	return
}

func (r *pantryRepository) IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 bool, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.IsUserInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.IsUserInPantry"), zap.Any("params", __logParams))
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PantryUser{}).
		Where("pantry_id = ? AND user_id = ? AND deleted_at IS NULL", pantryID, userID).
		Count(&count).Error
	result0 = count > 0
	result1 = err
	return
}

func (r *pantryRepository) IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (result0 bool, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.IsUserOwner"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.IsUserOwner"), zap.Any("params", __logParams))
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PantryUser{}).
		Where("pantry_id = ? AND user_id = ? AND role = ? AND deleted_at IS NULL", pantryID, userID, "owner").
		Count(&count).Error
	result0 = count > 0
	result1 = err
	return
}

func (r *pantryRepository) AddUserToPantry(ctx context.Context, pantryUser *model.PantryUser) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryUser": pantryUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.AddUserToPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.AddUserToPantry"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(pantryUser).Error
	return
}

func (r *pantryRepository) RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.RemoveUserFromPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.RemoveUserFromPantry"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).
		Where("pantry_id = ? AND user_id = ?", pantryID, userID).
		Delete(&model.PantryUser{}).Error
	return
}

func (r *pantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) (result0 []*model.PantryUserInfo, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryRepository.ListUsersInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryRepository.ListUsersInPantry"), zap.Any("params", __logParams))
	var users []*model.PantryUserInfo
	err := r.db.WithContext(ctx).
		Table("pantry_users").
		Select("pantry_users.id as id, pantry_users.pantry_id as pantry_id, pantry_users.user_id as user_id, users.email, pantry_users.role").
		Joins("JOIN users ON users.id = pantry_users.user_id").
		Where("pantry_users.pantry_id = ? AND pantry_users.deleted_at IS NULL", pantryID).
		Scan(&users).Error
	result0 = users
	result1 = err
	return
}
