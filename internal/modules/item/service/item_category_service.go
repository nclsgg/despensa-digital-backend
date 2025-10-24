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
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type itemCategoryService struct {
	repo       domain.ItemCategoryRepository
	pantryRepo pantryDomain.PantryRepository
}

func NewItemCategoryService(repo domain.ItemCategoryRepository, pantryRepo pantryDomain.PantryRepository) domain.ItemCategoryService {
	return &itemCategoryService{repo: repo, pantryRepo: pantryRepo}
}

func (s *itemCategoryService) Create(ctx context.Context, input dto.CreateItemCategoryDTO, userID uuid.UUID) (*dto.ItemCategoryResponse, error) {
	logger := appLogger.FromContext(ctx)

	pantryID, err := uuid.Parse(input.PantryID)
	if err != nil {
		logger.Warn("invalid pantry ID",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, domain.ErrInvalidPantry
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, domain.ErrUnauthorized
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
		logger.Error("failed to create item category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("item category created",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "Create"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", itemCategory.ID.String()),
		zap.String("pantry_id", pantryID.String()),
	)
	return toItemCategoryResponse(itemCategory), nil
}

func (s *itemCategoryService) CreateDefault(ctx context.Context, input dto.CreateDefaultItemCategoryDTO, userID uuid.UUID) (*dto.ItemCategoryResponse, error) {
	logger := appLogger.FromContext(ctx)

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
		logger.Error("failed to create default item category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "CreateDefault"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("default item category created",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "CreateDefault"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", itemCategory.ID.String()),
	)
	return toItemCategoryResponse(itemCategory), nil
}

func (s *itemCategoryService) CloneDefaultCategoryToPantry(ctx context.Context, defaultCategoryID, pantryID uuid.UUID, userID uuid.UUID) (*dto.ItemCategoryResponse, error) {
	logger := appLogger.FromContext(ctx)

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}

	defaultCat, err := s.repo.FindByID(ctx, defaultCategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCategoryNotFound
		}
		logger.Error("failed to find default category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("category_id", defaultCategoryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if !defaultCat.IsDefault {
		logger.Warn("category is not default",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("category_id", defaultCategoryID.String()),
		)
		return nil, domain.ErrCategoryNotDefault
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
		logger.Error("failed to clone category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("category cloned to pantry",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", newCategory.ID.String()),
		zap.String("pantry_id", pantryID.String()),
	)
	return toItemCategoryResponse(newCategory), nil
}

func (s *itemCategoryService) Update(ctx context.Context, id uuid.UUID, input dto.UpdateItemCategoryDTO, userID uuid.UUID) (*dto.ItemCategoryResponse, error) {
	logger := appLogger.FromContext(ctx)

	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCategoryNotFound
		}
		logger.Error("failed to find category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("category_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, itemCategory.PantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", itemCategory.PantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", itemCategory.PantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}

	itemCategory.ApplyUpdate(input)
	itemCategory.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, itemCategory); err != nil {
		logger.Error("failed to update category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("category_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("category updated",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "Update"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", id.String()),
	)
	return toItemCategoryResponse(itemCategory), nil
}

func (s *itemCategoryService) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.ItemCategoryResponse, error) {
	logger := appLogger.FromContext(ctx)

	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCategoryNotFound
		}
		logger.Error("failed to find category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "FindByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("category_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, itemCategory.PantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "FindByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", itemCategory.PantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "FindByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", itemCategory.PantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}
	return toItemCategoryResponse(itemCategory), nil
}

func (s *itemCategoryService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrCategoryNotFound
		}
		logger.Error("failed to find category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Delete"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("category_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	isOwner := itemCategory.AddedBy == userID
	if !isOwner {
		logger.Warn("unauthorized category access",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "Delete"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("category_id", id.String()),
		)
		return domain.ErrUnauthorized
	}

	logger.Info("category deleted",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "Delete"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", id.String()),
	)
	return s.repo.Delete(ctx, id)
}

func (s *itemCategoryService) ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]*dto.ItemCategoryResponse, error) {
	logger := appLogger.FromContext(ctx)

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "ListByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "ListByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}

	itemCategories, err := s.repo.ListByPantryID(ctx, pantryID)
	if err != nil {
		logger.Error("failed to list categories",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "ListByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	return toItemCategoryResponseList(itemCategories), nil
}

func (s *itemCategoryService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.ItemCategoryResponse, error) {
	logger := appLogger.FromContext(ctx)

	itemCategories, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		logger.Error("failed to list categories by user",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "ListByUserID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	return toItemCategoryResponseList(itemCategories), nil
}

func toItemCategoryResponse(category *model.ItemCategory) *dto.ItemCategoryResponse {
	if category == nil {
		return nil
	}

	var deletedAt *string
	if category.DeletedAt.Valid {
		formatted := category.DeletedAt.Time.UTC().Format(time.RFC3339)
		deletedAt = &formatted
	}
	return &dto.ItemCategoryResponse{
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
}

func toItemCategoryResponseList(categories []*model.ItemCategory) []*dto.ItemCategoryResponse {
	responses := make([]*dto.ItemCategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, toItemCategoryResponse(category))
	}
	return responses
}
