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
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type fakeItemCategoryRepository struct {
	store     map[uuid.UUID]*model.ItemCategory
	createErr error
}

func newFakeItemCategoryRepository() (result0 *fakeItemCategoryRepository) {
	__logParams := map[string]any{}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "newFakeItemCategoryRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "newFakeItemCategoryRepository"), zap.Any("params", __logParams))
	result0 = &fakeItemCategoryRepository{store: make(map[uuid.UUID]*model.ItemCategory)}
	return
}

func (f *fakeItemCategoryRepository) Create(ctx context.Context, category *model.ItemCategory) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "category": category}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakeItemCategoryRepository.Create"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakeItemCategoryRepository.Create"), zap.Any("params", __logParams))
	if f.createErr != nil {
		result0 = f.createErr
		return
	}
	clone := *category
	f.store[category.ID] = &clone
	result0 = nil
	return
}

func (f *fakeItemCategoryRepository) Update(ctx context.Context, category *model.ItemCategory) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "category": category}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakeItemCategoryRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakeItemCategoryRepository.Update"), zap.Any("params", __logParams))
	if _, ok := f.store[category.ID]; !ok {
		result0 = errors.New("not found")
		return
	}
	clone := *category
	f.store[category.ID] = &clone
	result0 = nil
	return
}

func (f *fakeItemCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (result0 *model.ItemCategory, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakeItemCategoryRepository.FindByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakeItemCategoryRepository.FindByID"), zap.Any("params", __logParams))
	if cat, ok := f.store[id]; ok {
		clone := *cat
		result0 = &clone
		result1 = nil
		return
	}
	result0 = nil
	result1 = gorm.ErrRecordNotFound
	return
}

func (f *fakeItemCategoryRepository) Delete(ctx context.Context, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakeItemCategoryRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakeItemCategoryRepository.Delete"), zap.Any("params", __logParams))
	if _, ok := f.store[id]; !ok {
		result0 = gorm.ErrRecordNotFound
		return
	}
	delete(f.store, id)
	result0 = nil
	return
}

func (f *fakeItemCategoryRepository) ListByPantryID(ctx context.Context, pantryID uuid.UUID) (result0 []*model.ItemCategory, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakeItemCategoryRepository.ListByPantryID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakeItemCategoryRepository.ListByPantryID"), zap.Any("params", __logParams))
	var result []*model.ItemCategory
	for _, cat := range f.store {
		if cat.PantryID == pantryID {
			clone := *cat
			result = append(result, &clone)
		}
	}
	result0 = result
	result1 = nil
	return
}

func (f *fakeItemCategoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) (result0 []*model.ItemCategory, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakeItemCategoryRepository.ListByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakeItemCategoryRepository.ListByUserID"), zap.Any("params", __logParams))
	var result []*model.ItemCategory
	for _, cat := range f.store {
		if cat.AddedBy == userID {
			clone := *cat
			result = append(result, &clone)
		}
	}
	result0 = result
	result1 = nil
	return
}

type fakePantryRepository struct {
	memberships       map[uuid.UUID]map[uuid.UUID]bool
	isUserInPantryErr error
}

func newFakePantryRepository() (result0 *fakePantryRepository) {
	__logParams := map[string]any{}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "newFakePantryRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "newFakePantryRepository"), zap.Any("params", __logParams))
	result0 = &fakePantryRepository{memberships: make(map[uuid.UUID]map[uuid.UUID]bool)}
	return
}

func (f *fakePantryRepository) setMembership(pantryID, userID uuid.UUID, isMember bool) {
	__logParams := map[string]any{"f": f, "pantryID": pantryID, "userID": userID, "isMember": isMember}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.setMembership"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.setMembership"), zap.Any("params", __logParams))
	if _, ok := f.memberships[pantryID]; !ok {
		f.memberships[pantryID] = make(map[uuid.UUID]bool)
	}
	f.memberships[pantryID][userID] = isMember
}

func (f *fakePantryRepository) Create(ctx context.Context, pantry *pantryModel.Pantry) (result0 *pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantry": pantry}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.Create"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.Create"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) Delete(ctx context.Context, pantryID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.Delete"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.Delete"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) Update(ctx context.Context, pantry *pantryModel.Pantry) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantry": pantry}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.Update"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.Update"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (result0 *pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.GetByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.GetByID"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) (result0 []*pantryModel.Pantry, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.GetByUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.GetByUser"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 bool, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.IsUserInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.IsUserInPantry"), zap.Any("params", __logParams))
	if f.isUserInPantryErr != nil {
		result0 = false
		result1 = f.isUserInPantryErr
		return
	}
	if users, ok := f.memberships[pantryID]; ok {
		result0 = users[userID]
		result1 = nil
		return
	}
	result0 = false
	result1 = nil
	return
}

func (f *fakePantryRepository) IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (result0 bool, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.IsUserOwner"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.IsUserOwner"), zap.Any("params", __logParams))
	result0 = false
	result1 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) AddUserToPantry(ctx context.Context, pantryUser *pantryModel.PantryUser) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryUser": pantryUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.AddUserToPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.AddUserToPantry"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.RemoveUserFromPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.RemoveUserFromPantry"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) UpdatePantryUserRole(ctx context.Context, pantryID, userID uuid.UUID, newRole string) (result0 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID, "userID": userID, "newRole": newRole}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.UpdatePantryUserRole"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.UpdatePantryUserRole"), zap.Any("params", __logParams))
	result0 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) GetPantryUser(ctx context.Context, pantryID, userID uuid.UUID) (result0 *pantryModel.PantryUser, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.GetPantryUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.GetPantryUser"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func (f *fakePantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) (result0 []*pantryModel.PantryUserInfo, result1 error) {
	__logParams := map[string]any{"f": f, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*fakePantryRepository.ListUsersInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*fakePantryRepository.ListUsersInPantry"), zap.Any("params", __logParams))
	result0 = nil
	result1 = errors.New("not implemented")
	return
}

func TestItemCategoryService_Create_InvalidPantry(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestItemCategoryService_Create_InvalidPantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestItemCategoryService_Create_InvalidPantry"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestItemCategoryService_Create_Unauthorized"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestItemCategoryService_Create_Unauthorized"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestItemCategoryService_Create_Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestItemCategoryService_Create_Success"), zap.Any("params", __logParams))
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
		zap.L().Error("function.error", zap.String("func", "TestItemCategoryService_Create_Success"), zap.Error(err), zap.Any("params", __logParams))
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
		zap.L().Error("function.error", zap.String("func", "TestItemCategoryService_Create_Success"), zap.Error(err), zap.Any("params", __logParams))
		t.Fatalf("created_at should be RFC3339: %v", err)
	}
	parsedUpdated, err := time.Parse(time.RFC3339, resp.UpdatedAt)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "TestItemCategoryService_Create_Success"), zap.Error(err), zap.Any("params", __logParams))
		t.Fatalf("updated_at should be RFC3339: %v", err)
	}
	if parsedCreated.IsZero() || parsedUpdated.IsZero() {
		t.Fatalf("parsed timestamps should not be zero")
	}
}
