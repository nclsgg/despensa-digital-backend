package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) domain.ItemRepository {
	return &itemRepository{db}
}

func (r *itemRepository) Create(ctx context.Context, item *model.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *itemRepository) Update(ctx context.Context, item *model.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *itemRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	var item model.Item
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Item{}, "id = ?", id).Error
}

func (r *itemRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) ([]*model.Item, error) {
	var items []*model.Item
	if err := r.db.WithContext(ctx).Where("pantry_id = ?", pantryID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
