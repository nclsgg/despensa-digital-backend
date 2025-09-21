package repository

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
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

func (r *itemRepository) CountByPantryID(ctx context.Context, pantryID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Item{}).Where("pantry_id = ?", pantryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *itemRepository) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters dto.ItemFilterDTO) ([]*model.Item, error) {
	query := r.db.WithContext(ctx).Where("pantry_id = ?", pantryID)

	// Filtro por preço mínimo
	if filters.MinPrice != nil {
		query = query.Where("(quantity * price_per_unit) >= ?", *filters.MinPrice)
	}

	// Filtro por preço máximo
	if filters.MaxPrice != nil {
		query = query.Where("(quantity * price_per_unit) <= ?", *filters.MaxPrice)
	}

	// Filtro por data de vencimento até
	if filters.ExpiresUntil != "" {
		layout := "2006-01-02"
		expiresUntil, err := time.Parse(layout, filters.ExpiresUntil)
		if err == nil {
			query = query.Where("expires_at IS NOT NULL AND expires_at <= ?", expiresUntil)
		}
	}

	// Filtro por nome (busca parcial, case insensitive)
	if filters.Name != nil && strings.TrimSpace(*filters.Name) != "" {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(*filters.Name)) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchTerm)
	}

	// Filtro por categoria
	if filters.CategoryID != nil {
		categoryUUID, err := uuid.Parse(*filters.CategoryID)
		if err == nil {
			query = query.Where("category_id = ?", categoryUUID)
		}
	}

	// Ordenação
	if filters.SortBy != nil {
		sortDirection := "asc"
		if filters.SortDirection != nil && strings.ToLower(*filters.SortDirection) == "desc" {
			sortDirection = "desc"
		}

		switch strings.ToLower(*filters.SortBy) {
		case "price":
			query = query.Order("(quantity * price_per_unit) " + sortDirection)
		case "expires_at":
			// Colocar nulls por último quando ordenar por data
			if sortDirection == "asc" {
				query = query.Order("expires_at ASC NULLS LAST")
			} else {
				query = query.Order("expires_at DESC NULLS LAST")
			}
		case "category":
			query = query.Order("category_id " + sortDirection + " NULLS LAST")
		case "name":
			query = query.Order("name " + sortDirection)
		default:
			// Ordenação padrão por nome
			query = query.Order("name ASC")
		}
	} else {
		// Ordenação padrão
		query = query.Order("name ASC")
	}

	var items []*model.Item
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
