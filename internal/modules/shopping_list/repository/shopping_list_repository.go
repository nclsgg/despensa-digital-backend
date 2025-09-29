package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type shoppingListRepository struct {
	db *gorm.DB
}

func NewShoppingListRepository(db *gorm.DB) (result0 domain.ShoppingListRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewShoppingListRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewShoppingListRepository"), zap.Any("params", __logParams))
	result0 = &shoppingListRepository{db: db}
	return
}

func (r *shoppingListRepository) Create(ctx context.Context, shoppingList *model.ShoppingList) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "shoppingList": shoppingList}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.Create"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(shoppingList).Error
	return
}

func (r *shoppingListRepository) GetByID(ctx context.Context, id uuid.UUID) (result0 *model.ShoppingList, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.GetByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.GetByID"), zap.Any("params", __logParams))
	var shoppingList model.ShoppingList
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&shoppingList).Error
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListRepository.GetByID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &shoppingList
	result1 = nil
	return
}

func (r *shoppingListRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) (result0 []*model.ShoppingList, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "userID": userID, "limit": limit, "offset": offset}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.GetByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.GetByUserID"), zap.Any("params", __logParams))
	var shoppingLists []*model.ShoppingList
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Preload("Items").Find(&shoppingLists).Error
	result0 = shoppingLists
	result1 = err
	return
}

func (r *shoppingListRepository) Update(ctx context.Context, shoppingList *model.ShoppingList) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "shoppingList": shoppingList}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func",

			// Delete items first
			"*shoppingListRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.Update"),

		// Then delete the shopping list
		zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(shoppingList).Error
	return
}

func (r *shoppingListRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.Delete"), zap.Any("params", __logParams))

	if err := r.db.WithContext(ctx).Where("shopping_list_id = ?", id).Delete(&model.ShoppingListItem{}).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListRepository.Delete"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}
	result0 = r.db.WithContext(ctx).Delete(&model.ShoppingList{}, id).Error
	return
}

func (r *shoppingListRepository) CreateItem(ctx context.Context, item *model.ShoppingListItem) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.CreateItem"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.CreateItem"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(item).Error
	return
}

func (r *shoppingListRepository) UpdateItem(ctx context.Context, item *model.ShoppingListItem) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.UpdateItem"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.UpdateItem"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(item).Error
	return
}

func (r *shoppingListRepository) DeleteItem(ctx context.Context, itemID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "itemID": itemID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.DeleteItem"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.DeleteItem"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Delete(&model.ShoppingListItem{}, itemID).Error
	return
}

func (r *shoppingListRepository) GetItemsByShoppingListID(ctx context.Context, shoppingListID uuid.UUID) (result0 []*model.ShoppingListItem, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "shoppingListID": shoppingListID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.GetItemsByShoppingListID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.GetItemsByShoppingListID"), zap.Any("params", __logParams))
	var items []*model.ShoppingListItem
	err := r.db.WithContext(ctx).Where("shopping_list_id = ?", shoppingListID).Find(&items).Error
	result0 = items
	result1 = err
	return
}

func (r *shoppingListRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (result0 int64, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListRepository.CountByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListRepository.CountByUserID"), zap.Any("params", __logParams))
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ShoppingList{}).Where("user_id = ?", userID).Count(&count).Error
	result0 = count
	result1 = err
	return
}
