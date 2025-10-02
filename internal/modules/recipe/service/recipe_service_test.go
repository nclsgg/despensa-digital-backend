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
	"go.uber.org/zap"
)

type stubItemRepository struct {
	items []*model.Item
	err   error
}

func (s *stubItemRepository) Create(ctx context.Context, item *model.Item) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubItemRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubItemRepository.Create"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubItemRepository) Update(ctx context.Context, item *model.Item) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubItemRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubItemRepository.Update"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubItemRepository) FindByID(ctx context.Context, id uuid.UUID) (result0 *model.Item, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubItemRepository.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubItemRepository.FindByID"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (s *stubItemRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubItemRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubItemRepository.Delete"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubItemRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 []*model.Item, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubItemRepository.ListByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubItemRepository.ListByPantryID"), zap.Any("params", __logParams))
	result0 = s.items
	result1 = s.err
	return
}

func (s *stubItemRepository) FilterByPantryID(ctx context.Context, pantryID uuid.UUID, filters dto.ItemFilterDTO) (result0 []*model.Item, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "filters": filters}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubItemRepository.FilterByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubItemRepository.FilterByPantryID"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (s *stubItemRepository) CountByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 int, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubItemRepository.CountByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubItemRepository.CountByPantryID"), zap.Any("params", __logParams))
	result0 = 0
	result1 = errors.New("not implemented")
	return
}

type stubPantryService struct {
	getPantryFn func(ctx context.Context, pantryID, userID uuid.UUID) (*pantryModel.Pantry, error)
}

func (s *stubPantryService) CreatePantry(ctx context.Context, name string, ownerID uuid.UUID) (result0 *pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "name": name, "ownerID": ownerID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.CreatePantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.CreatePantry"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (s *stubPantryService) GetPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 *pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.GetPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.GetPantry"), zap.Any("params", __logParams))
	if s.getPantryFn != nil {
		result0, result1 = s.getPantryFn(ctx, pantryID, userID)
		return
	}
	result0 = &pantryModel.Pantry{}
	result1 = nil
	return
}

func (s *stubPantryService) GetPantryWithItemCount(ctx context.Context, pantryID, userID uuid.UUID) (result0 *pantryModel.PantryWithItemCount, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.GetPantryWithItemCount"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.GetPantryWithItemCount"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (s *stubPantryService) ListPantriesByUser(ctx context.Context, userID uuid.UUID) (result0 []*pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.ListPantriesByUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.ListPantriesByUser"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (s *stubPantryService) ListPantriesWithItemCount(ctx context.Context, userID uuid.UUID) (result0 []*pantryModel.PantryWithItemCount, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.ListPantriesWithItemCount"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.ListPantriesWithItemCount"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (s *stubPantryService) DeletePantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.DeletePantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.DeletePantry"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubPantryService) UpdatePantry(ctx context.Context, pantryID, userID uuid.UUID, newName string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID, "newName": newName}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.UpdatePantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.UpdatePantry"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubPantryService) AddUserToPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "ownerID": ownerID, "targetUser": targetUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.AddUserToPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.AddUserToPantry"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubPantryService) RemoveUserFromPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "ownerID": ownerID, "targetUser": targetUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.RemoveUserFromPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.RemoveUserFromPantry"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubPantryService) RemoveSpecificUserFromPantry(ctx context.Context, pantryID, ownerID, targetUserID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "ownerID": ownerID, "targetUserID": targetUserID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.RemoveSpecificUserFromPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.RemoveSpecificUserFromPantry"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubPantryService) TransferOwnership(ctx context.Context, pantryID, currentOwnerID, newOwnerID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "currentOwnerID": currentOwnerID, "newOwnerID": newOwnerID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.TransferOwnership"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.TransferOwnership"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (s *stubPantryService) ListUsersInPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 []*pantryModel.PantryUserInfo, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*stubPantryService.ListUsersInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*stubPantryService.ListUsersInPantry"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func TestRecipeService_GetAvailableIngredients_Errors(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestRecipeService_GetAvailableIngredients_Errors"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestRecipeService_GetAvailableIngredients_Errors"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestRecipeService_GetAvailableIngredients_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestRecipeService_GetAvailableIngredients_Success"), zap.Any("params", __logParams))
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
		zap.L().Error("function.error", zap.String("func", "TestRecipeService_GetAvailableIngredients_Success"), zap.Error(err), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestRecipeService_ValidateRecipeRequest"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestRecipeService_ValidateRecipeRequest"), zap.Any("params", __logParams))
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
		zap.L().Error("function.error", zap.String("func", "TestRecipeService_ValidateRecipeRequest"), zap.Error(err), zap.Any("params", __logParams))
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
