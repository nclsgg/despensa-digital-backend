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

type itemService struct {
	repo       domain.ItemRepository
	pantryRepo pantryDomain.PantryRepository
}

func NewItemService(repo domain.ItemRepository, pantryRepo pantryDomain.PantryRepository) domain.ItemService {
	return &itemService{repo, pantryRepo}
}

func (s *itemService) Create(ctx context.Context, input dto.CreateItemDTO, userID uuid.UUID) (*model.Item, error) {
	pantryID := uuid.MustParse(input.PantryID)

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

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
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if input.CategoryID != nil {
		categoryID := uuid.MustParse(*input.CategoryID)
		item.CategoryID = &categoryID
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	item.TotalPrice = item.Quantity * item.PricePerUnit

	return item, nil
}

func (s *itemService) Update(ctx context.Context, id uuid.UUID, input dto.UpdateItemDTO, userID uuid.UUID) (*model.Item, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	item.ApplyUpdate(input)
	item.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	item.TotalPrice = item.Quantity * item.PricePerUnit

	return item, nil
}

func (s *itemService) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Item, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	return item, nil
}

func (s *itemService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, item.PantryID, userID)
	if err != nil || !isMember {
		return errors.New("user not authorized for this operation")
	}

	if item == nil {
		return errors.New("item not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *itemService) ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]*model.Item, error) {
	items, err := s.repo.ListByPantryID(ctx, pantryID)
	if err != nil {
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	return items, nil
}

func (s *itemService) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters dto.ItemFilterDTO, userID uuid.UUID) ([]*model.Item, error) {
	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	items, err := s.repo.FilterByPantryID(ctx, pantryID, filters)
	if err != nil {
		return nil, err
	}

	// Calcular total_price para cada item
	for _, item := range items {
		item.TotalPrice = item.Quantity * item.PricePerUnit
	}

	return items, nil
}
