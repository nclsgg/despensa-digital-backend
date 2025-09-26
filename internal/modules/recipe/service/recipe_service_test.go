package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	llmDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	pantrySvc "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/service"
	recipeDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/domain"
)

type stubItemRepository struct {
	items []*model.Item
	err   error
}

func (s *stubItemRepository) Create(ctx context.Context, item *model.Item) error {
	return errors.New("not implemented")
}

func (s *stubItemRepository) Update(ctx context.Context, item *model.Item) error {
	return errors.New("not implemented")
}

func (s *stubItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	return nil, errors.New("not implemented")
}

func (s *stubItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}

func (s *stubItemRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) ([]*model.Item, error) {
	return s.items, s.err
}

func (s *stubItemRepository) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters dto.ItemFilterDTO) ([]*model.Item, error) {
	return nil, errors.New("not implemented")
}

func (s *stubItemRepository) CountByPantryID(ctx context.Context, pantryID uuid.UUID) (int, error) {
	return 0, errors.New("not implemented")
}

type stubPantryService struct {
	getPantryFn func(ctx context.Context, pantryID, userID uuid.UUID) (*pantryModel.Pantry, error)
}

func (s *stubPantryService) CreatePantry(ctx context.Context, name string, ownerID uuid.UUID) (*pantryModel.Pantry, error) {
	return nil, errors.New("not implemented")
}

func (s *stubPantryService) GetPantry(ctx context.Context, pantryID, userID uuid.UUID) (*pantryModel.Pantry, error) {
	if s.getPantryFn != nil {
		return s.getPantryFn(ctx, pantryID, userID)
	}
	return &pantryModel.Pantry{}, nil
}

func (s *stubPantryService) GetPantryWithItemCount(ctx context.Context, pantryID, userID uuid.UUID) (*pantryModel.PantryWithItemCount, error) {
	return nil, errors.New("not implemented")
}

func (s *stubPantryService) ListPantriesByUser(ctx context.Context, userID uuid.UUID) ([]*pantryModel.Pantry, error) {
	return nil, errors.New("not implemented")
}

func (s *stubPantryService) ListPantriesWithItemCount(ctx context.Context, userID uuid.UUID) ([]*pantryModel.PantryWithItemCount, error) {
	return nil, errors.New("not implemented")
}

func (s *stubPantryService) DeletePantry(ctx context.Context, pantryID, userID uuid.UUID) error {
	return errors.New("not implemented")
}

func (s *stubPantryService) UpdatePantry(ctx context.Context, pantryID, userID uuid.UUID, newName string) error {
	return errors.New("not implemented")
}

func (s *stubPantryService) AddUserToPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error {
	return errors.New("not implemented")
}

func (s *stubPantryService) RemoveUserFromPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error {
	return errors.New("not implemented")
}

func (s *stubPantryService) ListUsersInPantry(ctx context.Context, pantryID, userID uuid.UUID) ([]*pantryModel.PantryUserInfo, error) {
	return nil, errors.New("not implemented")
}

func TestRecipeService_GetAvailableIngredients_Errors(t *testing.T) {
	pantryID := uuid.New()
	userID := uuid.New()

	testCases := []struct {
		name         string
		getPantryErr error
		expectErr    error
	}{
		{
			name:         "unauthorized",
			getPantryErr: pantrySvc.ErrUnauthorized,
			expectErr:    recipeDomain.ErrUnauthorized,
		},
		{
			name:         "pantry not found",
			getPantryErr: pantrySvc.ErrPantryNotFound,
			expectErr:    recipeDomain.ErrPantryNotFound,
		},
	}

	for _, tc := range testCases {
		t := t
		t.Run(tc.name, func(t *testing.T) {
			repo := &stubItemRepository{}
			pantrySvcStub := &stubPantryService{
				getPantryFn: func(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*pantryModel.Pantry, error) {
					return nil, tc.getPantryErr
				},
			}

			svc := &recipeService{
				llmService:     nil,
				itemRepository: repo,
				pantryService:  pantrySvcStub,
				promptBuilder:  nil,
			}

			_, err := svc.GetAvailableIngredients(context.Background(), pantryID, userID)
			if !errors.Is(err, tc.expectErr) {
				t.Fatalf("expected error %v, got %v", tc.expectErr, err)
			}
		})
	}
}

func TestRecipeService_GetAvailableIngredients_Success(t *testing.T) {
	pantryID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	repo := &stubItemRepository{
		items: []*model.Item{
			{
				PantryID:     pantryID,
				AddedBy:      userID,
				Name:         "  Apple  ",
				Quantity:     2,
				Unit:         "  pcs ",
				PricePerUnit: 1,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				PantryID:     pantryID,
				AddedBy:      userID,
				Name:         "Zero",
				Quantity:     0,
				Unit:         "kg",
				PricePerUnit: 1,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}

	pantrySvcStub := &stubPantryService{
		getPantryFn: func(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*pantryModel.Pantry, error) {
			return &pantryModel.Pantry{ID: id}, nil
		},
	}

	svc := &recipeService{
		llmService:     nil,
		itemRepository: repo,
		pantryService:  pantrySvcStub,
		promptBuilder:  nil,
	}

	ingredients, err := svc.GetAvailableIngredients(context.Background(), pantryID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ingredients) != 1 {
		t.Fatalf("expected 1 ingredient, got %d", len(ingredients))
	}

	ing := ingredients[0]
	if ing.Name != "Apple" {
		t.Fatalf("expected trimmed name 'Apple', got '%s'", ing.Name)
	}
	if ing.Unit != "pcs" {
		t.Fatalf("expected trimmed unit 'pcs', got '%s'", ing.Unit)
	}
	if ing.Quantity != 2 {
		t.Fatalf("expected quantity 2, got %f", ing.Quantity)
	}
}

func TestRecipeService_ValidateRecipeRequest(t *testing.T) {
	svc := &recipeService{}

	_, err := svc.validateRecipeRequest(&llmDTO.RecipeRequestDTO{})
	if !errors.Is(err, recipeDomain.ErrInvalidRequest) {
		t.Fatalf("expected invalid request error, got %v", err)
	}

	validPantry := uuid.New().String()
	_, err = svc.validateRecipeRequest(&llmDTO.RecipeRequestDTO{
		PantryID:    validPantry,
		CookingTime: 30,
		MealType:    "dinner",
		Difficulty:  "medium",
		ServingSize: 4,
	})
	if err != nil {
		t.Fatalf("expected valid request, got error %v", err)
	}

	_, err = svc.validateRecipeRequest(&llmDTO.RecipeRequestDTO{
		PantryID: validPantry,
		MealType: "invalid",
	})
	if !errors.Is(err, recipeDomain.ErrInvalidRequest) {
		t.Fatalf("expected invalid request due to meal type, got %v", err)
	}
}
