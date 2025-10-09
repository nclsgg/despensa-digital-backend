package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) (result0 *recipeRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewRecipeRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewRecipeRepository"), zap.Any("params", __logParams))
	result0 = &recipeRepository{db: db}
	return
}

func (r *recipeRepository) Create(ctx context.Context, recipe *model.Recipe) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "recipe": recipe}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeRepository.Create"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(recipe).Error
	return
}

func (r *recipeRepository) CreateMany(ctx context.Context, recipes []*model.Recipe) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "recipes": recipes}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeRepository.CreateMany"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeRepository.CreateMany"), zap.Any("params", __logParams))

	// Use a transaction for atomic creation
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, recipe := range recipes {
			if err := tx.Create(recipe).Error; err != nil {
				return err
			}
		}
		return nil
	})

	result0 = err
	return
}

func (r *recipeRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (result0 *model.Recipe, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeRepository.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeRepository.FindByID"), zap.Any("params", __logParams))

	var recipe model.Recipe
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&recipe).Error

	if err != nil {
		result0 = nil
		result1 = err
		return
	}

	result0 = &recipe
	result1 = nil
	return
}

func (r *recipeRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (result0 []*model.Recipe, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeRepository.FindByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeRepository.FindByUserID"), zap.Any("params", __logParams))

	var recipes []*model.Recipe
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&recipes).Error

	if err != nil {
		result0 = nil
		result1 = err
		return
	}

	result0 = recipes
	result1 = nil
	return
}

func (r *recipeRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeRepository.Delete"), zap.Any("params", __logParams))

	result0 = r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Recipe{}).Error
	return
}
