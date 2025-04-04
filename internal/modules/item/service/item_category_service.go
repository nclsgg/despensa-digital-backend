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

type itemCategoryService struct {
	repo       domain.ItemCategoryRepository
	pantryRepo pantryDomain.PantryRepository
}

func NewItemCategoryService(repo domain.ItemCategoryRepository, pantryRepo pantryDomain.PantryRepository) domain.ItemCategoryService {
	return &itemCategoryService{repo, pantryRepo}
}

func (s *itemCategoryService) Create(ctx context.Context, input dto.CreateItemCategoryDTO, userID uuid.UUID) (*model.ItemCategory, error) {
	pantryID := uuid.MustParse(input.PantryID)

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	itemCategory := &model.ItemCategory{
		ID:        uuid.New(),
		PantryID:  pantryID,
		AddedBy:   userID,
		Name:      input.Name,
		Color:     input.Color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, itemCategory); err != nil {
		return nil, err
	}

	return itemCategory, nil
}

func (s *itemCategoryService) CreateDefault(ctx context.Context, input dto.CreateDefaultItemCategoryDTO, userID uuid.UUID) (*model.ItemCategory, error) {
	itemCategory := &model.ItemCategory{
		ID:        uuid.New(),
		AddedBy:   userID,
		Name:      input.Name,
		Color:     input.Color,
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, itemCategory); err != nil {
		return nil, err
	}

	return itemCategory, nil
}

func (s *itemCategoryService) CloneDefaultCategoryToPantry(ctx context.Context, defaultCategoryID, pantryID uuid.UUID, userID uuid.UUID) (*model.ItemCategory, error) {
	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	defaultCat, err := s.repo.FindByID(ctx, defaultCategoryID)
	if err != nil {
		return nil, err
	}

	if !defaultCat.IsDefault {
		return nil, errors.New("category is not default")
	}

	newCategory := &model.ItemCategory{
		ID:        uuid.New(),
		PantryID:  pantryID,
		AddedBy:   userID,
		Name:      defaultCat.Name,
		Color:     defaultCat.Color,
		IsDefault: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, newCategory); err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (s *itemCategoryService) Update(ctx context.Context, id uuid.UUID, input dto.UpdateItemCategoryDTO, userID uuid.UUID) (*model.ItemCategory, error) {
	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, itemCategory.PantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	itemCategory.ApplyUpdate(input)

	if err := s.repo.Update(ctx, itemCategory); err != nil {
		return nil, err
	}

	return itemCategory, nil
}

func (s *itemCategoryService) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.ItemCategory, error) {
	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isMember, err := s.pantryRepo.IsUserInPantry(ctx, itemCategory.PantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	return itemCategory, nil
}

func (s *itemCategoryService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	itemCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	isOwner := itemCategory.AddedBy == userID
	if !isOwner {
		return errors.New("user not authorized for this operation")
	}

	return s.repo.Delete(ctx, id)
}

func (s *itemCategoryService) ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]*model.ItemCategory, error) {
	isMember, err := s.pantryRepo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		return nil, errors.New("user not authorized for this operation")
	}

	itemCategories, err := s.repo.ListByPantryID(ctx, pantryID)
	if err != nil {
		return nil, err
	}

	return itemCategories, nil
}

func (s *itemCategoryService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ItemCategory, error) {
	itemCategories, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return itemCategories, nil
}
