package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	itemDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	"gorm.io/gorm"
)

type fakeItemCategoryRepository struct {
	store     map[uuid.UUID]*model.ItemCategory
	createErr error
}

func newFakeItemCategoryRepository() *fakeItemCategoryRepository {
	return &fakeItemCategoryRepository{store: make(map[uuid.UUID]*model.ItemCategory)}
}

func (f *fakeItemCategoryRepository) Create(ctx context.Context, category *model.ItemCategory) error {
	if f.createErr != nil {
		return f.createErr
	}
	clone := *category
	f.store[category.ID] = &clone
	return nil
}

func (f *fakeItemCategoryRepository) Update(ctx context.Context, category *model.ItemCategory) error {
	if _, ok := f.store[category.ID]; !ok {
		return errors.New("not found")
	}
	clone := *category
	f.store[category.ID] = &clone
	return nil
}

func (f *fakeItemCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.ItemCategory, error) {
	if cat, ok := f.store[id]; ok {
		clone := *cat
		return &clone, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeItemCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := f.store[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(f.store, id)
	return nil
}

func (f *fakeItemCategoryRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) ([]*model.ItemCategory, error) {
	var result []*model.ItemCategory
	for _, cat := range f.store {
		if cat.PantryID == pantryID {
			clone := *cat
			result = append(result, &clone)
		}
	}
	return result, nil
}

func (f *fakeItemCategoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ItemCategory, error) {
	var result []*model.ItemCategory
	for _, cat := range f.store {
		if cat.AddedBy == userID {
			clone := *cat
			result = append(result, &clone)
		}
	}
	return result, nil
}

type fakePantryRepository struct {
	memberships       map[uuid.UUID]map[uuid.UUID]bool
	isUserInPantryErr error
}

func newFakePantryRepository() *fakePantryRepository {
	return &fakePantryRepository{memberships: make(map[uuid.UUID]map[uuid.UUID]bool)}
}

func (f *fakePantryRepository) setMembership(pantryID, userID uuid.UUID, isMember bool) {
	if _, ok := f.memberships[pantryID]; !ok {
		f.memberships[pantryID] = make(map[uuid.UUID]bool)
	}
	f.memberships[pantryID][userID] = isMember
}

func (f *fakePantryRepository) Create(ctx context.Context, pantry *pantryModel.Pantry) (*pantryModel.Pantry, error) {
	return nil, errors.New("not implemented")
}

func (f *fakePantryRepository) Delete(ctx context.Context, pantryID uuid.UUID) error {
	return errors.New("not implemented")
}

func (f *fakePantryRepository) Update(ctx context.Context, pantry *pantryModel.Pantry) error {
	return errors.New("not implemented")
}

func (f *fakePantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (*pantryModel.Pantry, error) {
	return nil, errors.New("not implemented")
}

func (f *fakePantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*pantryModel.Pantry, error) {
	return nil, errors.New("not implemented")
}

func (f *fakePantryRepository) IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	if f.isUserInPantryErr != nil {
		return false, f.isUserInPantryErr
	}
	if users, ok := f.memberships[pantryID]; ok {
		return users[userID], nil
	}
	return false, nil
}

func (f *fakePantryRepository) IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	return false, errors.New("not implemented")
}

func (f *fakePantryRepository) AddUserToPantry(ctx context.Context, pantryUser *pantryModel.PantryUser) error {
	return errors.New("not implemented")
}

func (f *fakePantryRepository) RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) error {
	return errors.New("not implemented")
}

func (f *fakePantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) ([]*pantryModel.PantryUserInfo, error) {
	return nil, errors.New("not implemented")
}

func TestItemCategoryService_Create_InvalidPantry(t *testing.T) {
	repo := newFakeItemCategoryRepository()
	pantryRepo := newFakePantryRepository()
	service := NewItemCategoryService(repo, pantryRepo)

	_, err := service.Create(context.Background(), dto.CreateItemCategoryDTO{
		PantryID: "invalid",
		Name:     "Cereais",
		Color:    "#FFFFFF",
	}, uuid.New())

	if !errors.Is(err, itemDomain.ErrInvalidPantry) {
		t.Fatalf("expected ErrInvalidPantry, got %v", err)
	}
}

func TestItemCategoryService_Create_Unauthorized(t *testing.T) {
	repo := newFakeItemCategoryRepository()
	pantryRepo := newFakePantryRepository()
	service := NewItemCategoryService(repo, pantryRepo)

	pantryID := uuid.New()
	userID := uuid.New()
	pantryRepo.setMembership(pantryID, userID, false)

	_, err := service.Create(context.Background(), dto.CreateItemCategoryDTO{
		PantryID: pantryID.String(),
		Name:     "Temperos",
		Color:    "#000000",
	}, userID)

	if !errors.Is(err, itemDomain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestItemCategoryService_Create_Success(t *testing.T) {
	repo := newFakeItemCategoryRepository()
	pantryRepo := newFakePantryRepository()
	service := NewItemCategoryService(repo, pantryRepo)

	pantryID := uuid.New()
	userID := uuid.New()
	pantryRepo.setMembership(pantryID, userID, true)

	resp, err := service.Create(context.Background(), dto.CreateItemCategoryDTO{
		PantryID: pantryID.String(),
		Name:     "Grãos",
		Color:    "#F1F1F1",
	}, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatalf("expected response, got nil")
	}

	if resp.PantryID != pantryID.String() {
		t.Fatalf("expected pantry ID %s, got %s", pantryID.String(), resp.PantryID)
	}

	if resp.AddedBy != userID.String() {
		t.Fatalf("expected added_by %s, got %s", userID.String(), resp.AddedBy)
	}

	if resp.IsDefault {
		t.Fatalf("expected IsDefault to be false")
	}

	if resp.CreatedAt == "" || resp.UpdatedAt == "" {
		t.Fatalf("expected timestamps to be populated")
	}

	parsedCreated, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		t.Fatalf("created_at should be RFC3339: %v", err)
	}
	parsedUpdated, err := time.Parse(time.RFC3339, resp.UpdatedAt)
	if err != nil {
		t.Fatalf("updated_at should be RFC3339: %v", err)
	}
	if parsedCreated.IsZero() || parsedUpdated.IsZero() {
		t.Fatalf("parsed timestamps should not be zero")
	}
}
