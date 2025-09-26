package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
	"gorm.io/gorm"
)

type shoppingListRepository struct {
	db *gorm.DB
}

func NewShoppingListRepository(db *gorm.DB) domain.ShoppingListRepository {
	return &shoppingListRepository{db: db}
}

func (r *shoppingListRepository) Create(ctx context.Context, shoppingList *model.ShoppingList) error {
	return r.db.WithContext(ctx).Create(shoppingList).Error
}

func (r *shoppingListRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ShoppingList, error) {
	var shoppingList model.ShoppingList
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&shoppingList).Error
	if err != nil {
		return nil, err
	}
	return &shoppingList, nil
}

func (r *shoppingListRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.ShoppingList, error) {
	var shoppingLists []*model.ShoppingList
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Preload("Items").Find(&shoppingLists).Error
	return shoppingLists, err
}

func (r *shoppingListRepository) Update(ctx context.Context, shoppingList *model.ShoppingList) error {
	return r.db.WithContext(ctx).Save(shoppingList).Error
}

func (r *shoppingListRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Delete items first
	if err := r.db.WithContext(ctx).Where("shopping_list_id = ?", id).Delete(&model.ShoppingListItem{}).Error; err != nil {
		return err
	}
	// Then delete the shopping list
	return r.db.WithContext(ctx).Delete(&model.ShoppingList{}, id).Error
}

func (r *shoppingListRepository) CreateItem(ctx context.Context, item *model.ShoppingListItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *shoppingListRepository) UpdateItem(ctx context.Context, item *model.ShoppingListItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *shoppingListRepository) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ShoppingListItem{}, itemID).Error
}

func (r *shoppingListRepository) GetItemsByShoppingListID(ctx context.Context, shoppingListID uuid.UUID) ([]*model.ShoppingListItem, error) {
	var items []*model.ShoppingListItem
	err := r.db.WithContext(ctx).Where("shopping_list_id = ?", shoppingListID).Find(&items).Error
	return items, err
}

func (r *shoppingListRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ShoppingList{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
