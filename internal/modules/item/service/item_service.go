package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	pantryDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func parseTimePointer(dateStr string) (result0 *time.Time) {
	__logParams := map[string]any{"dateStr": dateStr}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "parseTimePointer"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "parseTimePointer"), zap.Any("params", __logParams))
	if dateStr == "" {
		result0 = nil
		return
	}
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "parseTimePointer"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		return
	}
	result0 = &parsedTime
	return
}

func formatTimePointer(t *time.Time) (result0 *string) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "formatTimePointer"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "formatTimePointer"), zap.Any("params", __logParams))
	if t == nil {
		result0 = nil
		return
	}
	formatted := t.UTC().Format(time.RFC3339)
	result0 = &formatted
	return
}

func toItemResponse(item *model.Item) (result0 *dto.ItemResponse) {
	__logParams := map[string]any{"item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "toItemResponse"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "toItemResponse"), zap.Any("params", __logParams))
	if item == nil {
		result0 = nil
		return
	}

	item.TotalPrice = item.Quantity * item.PricePerUnit

	var categoryID *string
	if item.CategoryID != nil {
		id := item.CategoryID.String()
		categoryID = &id
	}
	result0 = &dto.ItemResponse{
		ID:           item.ID.String(),
		PantryID:     item.PantryID.String(),
		AddedBy:      item.AddedBy.String(),
		Name:         item.Name,
		Quantity:     item.Quantity,
		Unit:         item.Unit,
		PricePerUnit: item.PricePerUnit,
		TotalPrice:   item.TotalPrice,
		CategoryID:   categoryID,
		ExpiresAt:    formatTimePointer(item.ExpiresAt),
		CreatedAt:    item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    item.UpdatedAt.UTC().Format(time.RFC3339),
	}
	return
}

func toItemResponseList(items []*model.Item) (result0 []*dto.ItemResponse) {
	__logParams := map[string]any{"items": items}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "toItemResponseList"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "toItemResponseList"), zap.Any("params", __logParams))
	responses := make([]*dto.ItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toItemResponse(item))
	}
	result0 = responses
	return
}

type itemService struct {
	repo       domain.ItemRepository
	pantryRepo pantryDomain.PantryRepository
}

func NewItemService(repo domain.ItemRepository, pantryRepo pantryDomain.PantryRepository) (result0 domain.ItemService) {
	__logParams := map[string]any{"repo": repo, "pantryRepo": pantryRepo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewItemService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewItemService"), zap.Any("params", __logParams))
	result0 = &itemService{repo, pantryRepo}
	return
}

func (s *itemService) Create(ctx context.Context, input dto.CreateItemDTO, userID uuid.UUID) (result0 *dto.ItemResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "input": input, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemService.Create"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemService.Create"), zap.Any("params", __logParams))
	pantryID, err := uuid.Parse(input.PantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Create"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = domain.ErrInvalidPantry
		return
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Create"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	now := time.Now().UTC()
	item := &model.Item{
		ID:           uuid.New(),
		PantryID:     pantryID,
		AddedBy:      userID,
		Name:         input.Name,
		Quantity:     input.Quantity,
		PricePerUnit: input.PricePerUnit,
		Unit:         input.Unit,
		CategoryID:   nil,
		ExpiresAt:    parseTimePointer(input.ExpiresAt),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if input.CategoryID != nil {
		categoryID, parseErr := uuid.Parse(*input.CategoryID)
		if parseErr == nil {
			item.CategoryID = &categoryID
		}
	}

	if err := s.repo.Create(ctx, item); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Create"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemResponse(item)
	result1 = nil
	return
}

func (s *itemService) Update(ctx context.Context, id uuid.UUID, input dto.UpdateItemDTO, userID uuid.UUID) (result0 *dto.ItemResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id, "input": input, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemService.Update"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemService.Update"), zap.Any("params", __logParams))
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Update"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Update"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	item.ApplyUpdate(input)
	item.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, item); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Update"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemResponse(item)
	result1 = nil
	return
}

func (s *itemService) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (result0 *dto.ItemResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemService.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemService.FindByID"), zap.Any("params", __logParams))
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.FindByID"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrItemNotFound
			return
		}
		result0 = nil
		result1 = err
		return
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.FindByID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}
	result0 = toItemResponse(item)
	result1 = nil
	return
}

func (s *itemService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemService.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemService.Delete"), zap.Any("params", __logParams))
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Delete"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = domain.ErrItemNotFound
			return
		}
		result0 = err
		return
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.Delete"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}
	if !isMember {
		result0 = domain.ErrUnauthorized
		return
	}
	result0 = s.repo.Delete(ctx, id)
	return
}

func (s *itemService) ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (result0 []*dto.ItemResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemService.ListByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemService.ListByPantryID"), zap.Any("params", __logParams))
	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.ListByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	items, err := s.repo.ListByPantryID(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.ListByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemResponseList(items)
	result1 = nil
	return
}

func (s *itemService) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters dto.ItemFilterDTO, userID uuid.UUID) (result0 []*dto.ItemResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "filters": filters, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemService.FilterByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemService.FilterByPantryID"), zap.Any("params", __logParams))
	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.FilterByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	items, err := s.repo.FilterByPantryID(ctx, pantryID, filters)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemService.FilterByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemResponseList(items)
	result1 = nil
	return
}
