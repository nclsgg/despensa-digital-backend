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

type itemCategoryService struct {
	repo       domain.ItemCategoryRepository
	pantryRepo pantryDomain.PantryRepository
}

func NewItemCategoryService(repo domain.ItemCategoryRepository, pantryRepo pantryDomain.PantryRepository) (result0 domain.ItemCategoryService) {
	__logParams := map[string]any{"repo": repo, "pantryRepo": pantryRepo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewItemCategoryService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewItemCategoryService"), zap.Any("params", __logParams))
	result0 = &itemCategoryService{repo: repo, pantryRepo: pantryRepo}
	return
}

func (s *itemCategoryService) Create(ctx context.Context, input dto.CreateItemCategoryDTO, userID uuid.UUID) (result0 *dto.ItemCategoryResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "input": input, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.Create"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.Create"), zap.Any("params", __logParams))
	pantryID, err := uuid.Parse(input.PantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.Create"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = domain.ErrInvalidPantry
		return
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.Create"), zap.Error(err), zap.Any("params", __logParams))
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
	itemCategory := &model.ItemCategory{
		ID:        uuid.New(),
		PantryID:  pantryID,
		AddedBy:   userID,
		Name:      input.Name,
		Color:     input.Color,
		IsDefault: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, itemCategory); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.Create"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemCategoryResponse(itemCategory)
	result1 = nil
	return
}

func (s *itemCategoryService) CreateDefault(ctx context.Context, input dto.CreateDefaultItemCategoryDTO, userID uuid.UUID) (result0 *dto.ItemCategoryResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "input": input, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.CreateDefault"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.CreateDefault"), zap.Any("params", __logParams))
	now := time.Now().UTC()
	itemCategory := &model.ItemCategory{
		ID:        uuid.New(),
		AddedBy:   userID,
		Name:      input.Name,
		Color:     input.Color,
		IsDefault: true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, itemCategory); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.CreateDefault"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemCategoryResponse(itemCategory)
	result1 = nil
	return
}

func (s *itemCategoryService) CloneDefaultCategoryToPantry(ctx context.Context, defaultCategoryID, pantryID uuid.UUID, userID uuid.UUID) (result0 *dto.ItemCategoryResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "defaultCategoryID": defaultCategoryID, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.CloneDefaultCategoryToPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.CloneDefaultCategoryToPantry"), zap.Any("params", __logParams))
	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.CloneDefaultCategoryToPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	defaultCat, err := s.repo.FindByID(ctx, defaultCategoryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.CloneDefaultCategoryToPantry"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrCategoryNotFound
			return
		}
		result0 = nil
		result1 = err
		return
	}

	if !defaultCat.IsDefault {
		result0 = nil
		result1 = domain.ErrCategoryNotDefault
		return
	}

	now := time.Now().UTC()
	newCategory := &model.ItemCategory{
		ID:        uuid.New(),
		PantryID:  pantryID,
		AddedBy:   userID,
		Name:      defaultCat.Name,
		Color:     defaultCat.Color,
		IsDefault: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, newCategory); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.CloneDefaultCategoryToPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemCategoryResponse(newCategory)
	result1 = nil
	return
}

func (s *itemCategoryService) Update(ctx context.Context, id uuid.UUID, input dto.UpdateItemCategoryDTO, userID uuid.UUID) (result0 *dto.ItemCategoryResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id, "input": input, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.Update"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.Update"), zap.Any("params", __logParams))
	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.Update"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrCategoryNotFound
			return
		}
		result0 = nil
		result1 = err
		return
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, itemCategory.PantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.Update"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	itemCategory.ApplyUpdate(input)
	itemCategory.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, itemCategory); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.Update"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemCategoryResponse(itemCategory)
	result1 = nil
	return
}

func (s *itemCategoryService) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (result0 *dto.ItemCategoryResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.FindByID"), zap.Any("params", __logParams))
	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.FindByID"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrCategoryNotFound
			return
		}
		result0 = nil
		result1 = err
		return
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, itemCategory.PantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.FindByID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}
	result0 = toItemCategoryResponse(itemCategory)
	result1 = nil
	return
}

func (s *itemCategoryService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.Delete"), zap.Any("params", __logParams))
	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.Delete"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = domain.ErrCategoryNotFound
			return
		}
		result0 = err
		return
	}

	isOwner := itemCategory.AddedBy == userID
	if !isOwner {
		result0 = domain.ErrUnauthorized
		return
	}
	result0 = s.repo.Delete(ctx, id)
	return
}

func (s *itemCategoryService) ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (result0 []*dto.ItemCategoryResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.ListByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.ListByPantryID"), zap.Any("params", __logParams))
	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.ListByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	itemCategories, err := s.repo.ListByPantryID(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.ListByPantryID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemCategoryResponseList(itemCategories)
	result1 = nil
	return
}

func (s *itemCategoryService) ListByUserID(ctx context.Context, userID uuid.UUID) (result0 []*dto.ItemCategoryResponse, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryService.ListByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryService.ListByUserID"), zap.Any("params", __logParams))
	itemCategories, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryService.ListByUserID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = toItemCategoryResponseList(itemCategories)
	result1 = nil
	return
}

func toItemCategoryResponse(category *model.ItemCategory) (result0 *dto.ItemCategoryResponse) {
	__logParams := map[string]any{"category": category}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "toItemCategoryResponse"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "toItemCategoryResponse"), zap.Any("params", __logParams))
	if category == nil {
		result0 = nil
		return
	}

	var deletedAt *string
	if category.DeletedAt.Valid {
		formatted := category.DeletedAt.Time.UTC().Format(time.RFC3339)
		deletedAt = &formatted
	}
	result0 = &dto.ItemCategoryResponse{
		ID:        category.ID.String(),
		PantryID:  category.PantryID.String(),
		AddedBy:   category.AddedBy.String(),
		Name:      category.Name,
		Color:     category.Color,
		IsDefault: category.IsDefault,
		CreatedAt: category.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.UTC().Format(time.RFC3339),
		DeletedAt: deletedAt,
	}
	return
}

func toItemCategoryResponseList(categories []*model.ItemCategory) (result0 []*dto.ItemCategoryResponse) {
	__logParams := map[string]any{"categories": categories}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "toItemCategoryResponseList"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "toItemCategoryResponseList"), zap.Any("params", __logParams))
	responses := make([]*dto.ItemCategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, toItemCategoryResponse(category))
	}
	result0 = responses
	return
}
