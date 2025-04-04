package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	"gorm.io/gorm"
)

type itemCategoryRepository struct {
	db *gorm.DB
}

func NewItemCategoryRepository(db *gorm.DB) domain.ItemCategoryRepository {
	return &itemCategoryRepository{db}
}

func (r *itemCategoryRepository) Create(ctx context.Context, itemCategory *model.ItemCategory) error {
	return r.db.WithContext(ctx).Create(itemCategory).Error
}

func (r *itemCategoryRepository) Update(ctx context.Context, itemCategory *model.ItemCategory) error {
	return r.db.WithContext(ctx).Save(itemCategory).Error
}

func (r *itemCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.ItemCategory, error) {
	var itemCategory model.ItemCategory
	if err := r.db.WithContext(ctx).First(&itemCategory, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &itemCategory, nil
}

func (r *itemCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ItemCategory{}, "id = ?", id).Error
}

func (r *itemCategoryRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) ([]*model.ItemCategory, error) {
	var itemCategories []*model.ItemCategory
	if err := r.db.WithContext(ctx).Where("pantry_id = ?", pantryID).Find(&itemCategories).Error; err != nil {
		return nil, err
	}
	return itemCategories, nil
}

func (r *itemCategoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ItemCategory, error) {
	var itemCategories []*model.ItemCategory
	if err := r.db.WithContext(ctx).Where("added_by = ?", userID).Find(&itemCategories).Error; err != nil {
		return nil, err
	}
	return itemCategories, nil
}
