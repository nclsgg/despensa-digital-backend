package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type itemCategoryRepository struct {
	db *gorm.DB
}

func NewItemCategoryRepository(db *gorm.DB) (result0 domain.ItemCategoryRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewItemCategoryRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewItemCategoryRepository"), zap.Any("params", __logParams))
	result0 = &itemCategoryRepository{db}
	return
}

func (r *itemCategoryRepository) Create(ctx context.Context, itemCategory *model.ItemCategory) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "itemCategory": itemCategory}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryRepository.Create"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(itemCategory).Error
	return
}

func (r *itemCategoryRepository) Update(ctx context.Context, itemCategory *model.ItemCategory) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "itemCategory": itemCategory}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryRepository.Update"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(itemCategory).Error
	return
}

func (r *itemCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (result0 *model.ItemCategory, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryRepository.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryRepository.FindByID"), zap.Any("params", __logParams))
	var itemCategory model.ItemCategory
	if err := r.db.WithContext(ctx).First(&itemCategory, "id = ?", id).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryRepository.FindByID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &itemCategory
	result1 = nil
	return
}

func (r *itemCategoryRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryRepository.Delete"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Delete(&model.ItemCategory{}, "id = ?", id).Error
	return
}

func (r *itemCategoryRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 []*model.ItemCategory, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryRepository.ListByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryRepository.ListByPantryID"), zap.Any("params", __logParams))
	var itemCategories []*model.ItemCategory
	if err := r.db.WithContext(ctx).Where("pantry_id = ?", pantryID).Find(&itemCategories).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryRepository.ListByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = itemCategories
	result1 = nil
	return
}

func (r *itemCategoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) (result0 []*model.ItemCategory, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryRepository.ListByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryRepository.ListByUserID"), zap.Any("params", __logParams))
	var itemCategories []*model.ItemCategory
	if err := r.db.WithContext(ctx).Where("added_by = ?", userID).Find(&itemCategories).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryRepository.ListByUserID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = itemCategories
	result1 = nil
	return
}
