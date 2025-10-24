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

func parseTimePointer(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil
	}
	return &parsedTime
}

func formatTimePointer(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.UTC().Format(time.RFC3339)
	return &formatted
}

func toItemResponse(item *model.Item) *dto.ItemResponse {
	if item == nil {
		return nil
	}

	totalPrice := item.TotalPrice
	if totalPrice == 0 {
		totalPrice = item.Quantity * item.PricePerUnit
	}

	var categoryID *string
	if item.CategoryID != nil {
		id := item.CategoryID.String()
		categoryID = &id
	}
	return &dto.ItemResponse{
		ID:           item.ID.String(),
		PantryID:     item.PantryID.String(),
		AddedBy:      item.AddedBy.String(),
		Name:         item.Name,
		Quantity:     item.Quantity,
		Unit:         item.Unit,
		PricePerUnit: item.PricePerUnit,
		TotalPrice:   totalPrice,
		CategoryID:   categoryID,
		ExpiresAt:    formatTimePointer(item.ExpiresAt),
		CreatedAt:    item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    item.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toItemResponseList(items []*model.Item) []*dto.ItemResponse {
	responses := make([]*dto.ItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toItemResponse(item))
	}
	return responses
}

type itemService struct {
	repo       domain.ItemRepository
	pantryRepo pantryDomain.PantryRepository
}

func NewItemService(repo domain.ItemRepository, pantryRepo pantryDomain.PantryRepository) domain.ItemService {
	return &itemService{repo, pantryRepo}
}

func (s *itemService) Create(ctx context.Context, input dto.CreateItemDTO, userID uuid.UUID) (*dto.ItemResponse, error) {
	logger := appLogger.FromContext(ctx)

	pantryID, err := uuid.Parse(input.PantryID)
	if err != nil {
		logger.Warn("invalid pantry ID",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, domain.ErrInvalidPantry
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, domain.ErrUnauthorized
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
		logger.Error("failed to create item",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Create"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("item created",
		zap.String(appLogger.FieldModule, "item"),
		zap.String(appLogger.FieldFunction, "Create"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("item_id", item.ID.String()),
		zap.String("pantry_id", pantryID.String()),
	)
	return toItemResponse(item), nil
}

func (s *itemService) Update(ctx context.Context, id uuid.UUID, input dto.UpdateItemDTO, userID uuid.UUID) (*dto.ItemResponse, error) {
	logger := appLogger.FromContext(ctx)

	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find item",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("item_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", item.PantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", item.PantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}

	item.ApplyUpdate(input)
	item.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, item); err != nil {
		logger.Error("failed to update item",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Update"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("item_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("item updated",
		zap.String(appLogger.FieldModule, "item"),
		zap.String(appLogger.FieldFunction, "Update"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("item_id", id.String()),
	)
	return toItemResponse(item), nil
}

func (s *itemService) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.ItemResponse, error) {
	logger := appLogger.FromContext(ctx)

	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrItemNotFound
		}
		logger.Error("failed to find item",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "FindByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("item_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "FindByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", item.PantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "FindByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", item.PantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}
	return toItemResponse(item), nil
}

func (s *itemService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrItemNotFound
		}
		logger.Error("failed to find item",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Delete"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("item_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Delete"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", item.PantryID.String()),
			zap.Error(err),
		)
		return err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "Delete"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", item.PantryID.String()),
		)
		return domain.ErrUnauthorized
	}

	logger.Info("item deleted",
		zap.String(appLogger.FieldModule, "item"),
		zap.String(appLogger.FieldFunction, "Delete"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("item_id", id.String()),
	)
	return s.repo.Delete(ctx, id)
}

func (s *itemService) ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]*dto.ItemResponse, error) {
	logger := appLogger.FromContext(ctx)

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "ListByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "ListByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}

	items, err := s.repo.ListByPantryID(ctx, pantryID)
	if err != nil {
		logger.Error("failed to list items",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "ListByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	return toItemResponseList(items), nil
}

func (s *itemService) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters dto.ItemFilterDTO, userID uuid.UUID) ([]*dto.ItemResponse, error) {
	logger := appLogger.FromContext(ctx)

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("failed to check pantry membership",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "FilterByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("unauthorized pantry access",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "FilterByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, domain.ErrUnauthorized
	}

	items, err := s.repo.FilterByPantryID(ctx, pantryID, filters)
	if err != nil {
		logger.Error("failed to filter items",
			zap.String(appLogger.FieldModule, "item"),
			zap.String(appLogger.FieldFunction, "FilterByPantryID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	return toItemResponseList(items), nil
}
