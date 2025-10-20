package repository

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) (result0 domain.ItemRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewItemRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewItemRepository"), zap.Any("params", __logParams))
	result0 = &itemRepository{db}
	return
}

func (r *itemRepository) Create(ctx context.Context, item *model.Item) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemRepository.Create"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(item).Error
	return
}

func (r *itemRepository) Update(ctx context.Context, item *model.Item) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemRepository.Update"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(item).Error
	return
}

func (r *itemRepository) FindByID(ctx context.Context, id uuid.UUID) (result0 *model.Item, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemRepository.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemRepository.FindByID"), zap.Any("params", __logParams))
	var item model.Item
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemRepository.FindByID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &item
	result1 = nil
	return
}

func (r *itemRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemRepository.Delete"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Delete(&model.Item{}, "id = ?", id).Error
	return
}

func (r *itemRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 []*model.Item, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemRepository.ListByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemRepository.ListByPantryID"), zap.Any("params", __logParams))
	var items []*model.Item
	if err := r.db.WithContext(ctx).Where("pantry_id = ?", pantryID).Find(&items).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemRepository.ListByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = items
	result1 = nil
	return
}

func (r *itemRepository) CountByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 int, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemRepository.CountByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemRepository.CountByPantryID"), zap.Any("params",

		// Filtro por preço mínimo
		__logParams))
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Item{}).
		Where("pantry_id = ?", pantryID).
		Distinct("id").
		Count(&count).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemRepository.CountByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = 0
		result1 = err
		return
	}
	result0 = int(count)
	result1 = nil
	return
}

func (r *itemRepository) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters dto.ItemFilterDTO) (result0 []*model.Item, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "pantryID": pantryID, "filters": filters}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemRepository.FilterByPantryID"), zap.Any("result", map[

		// Filtro por preço máximo
		string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func",

		// Filtro por data de vencimento até
		"*itemRepository.FilterByPantryID"), zap.Any("params", __logParams))
	query := r.db.WithContext(ctx).Where("pantry_id = ?", pantryID)

	if filters.MinPrice != nil {
		query = query.Where("(quantity * price_per_unit) >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("(quantity * price_per_unit) <= ?", *filters.MaxPrice)
	}

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
		zap.L().Error("function.error", zap.String("func", "*itemRepository.FilterByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = items
	result1 = nil
	return
}
