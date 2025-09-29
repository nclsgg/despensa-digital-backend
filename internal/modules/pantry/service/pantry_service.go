package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	itemDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	userDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"go.uber.org/zap"
)

type pantryService struct {
	repo     domain.PantryRepository
	userRepo userDomain.UserRepository
	itemRepo itemDomain.ItemRepository
}

var (
	ErrUnauthorized   = errors.New("user not authorized for this operation")
	ErrPantryNotFound = errors.New("pantry not found")
)

func NewPantryService(
	repo domain.PantryRepository,
	userRepo userDomain.UserRepository,
	itemRepo itemDomain.ItemRepository,
) (result0 domain.PantryService) {
	__logParams := map[string]any{"repo": repo, "userRepo": userRepo, "itemRepo": itemRepo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewPantryService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewPantryService"), zap.Any("params", __logParams))
	result0 = &pantryService{
		repo:     repo,
		userRepo: userRepo,
		itemRepo: itemRepo,
	}
	return
}

func (s *pantryService) CreatePantry(ctx context.Context, name string, ownerID uuid.UUID) (result0 *model.Pantry, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "name": name, "ownerID": ownerID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.CreatePantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.CreatePantry"), zap.Any("params", __logParams))
	pantry := &model.Pantry{
		ID:      uuid.New(),
		Name:    name,
		OwnerID: ownerID,
	}

	pantry, err := s.repo.Create(ctx, pantry)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.CreatePantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	pantryUser := &model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantry.ID,
		UserID:   ownerID,
		Role:     "owner",
	}

	if err := s.repo.AddUserToPantry(ctx, pantryUser); err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.CreatePantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = pantry
	result1 = nil
	return
}

func (s *pantryService) GetPantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (result0 *model.Pantry, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.GetPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.GetPantry"), zap.Any("params", __logParams))
	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		result0 = nil
		result1 = ErrUnauthorized
		return
	}

	pantry, err := s.repo.GetByID(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.GetPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = ErrPantryNotFound
		return
	}
	result0 = pantry
	result1 = nil
	return
}

func (s *pantryService) GetPantryWithItemCount(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (result0 *model.PantryWithItemCount, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.GetPantryWithItemCount"),

			// If we can't get item count, default to 0 instead of failing
			zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.GetPantryWithItemCount"), zap.Any("params", __logParams))
	pantry, err := s.GetPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.GetPantryWithItemCount"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	itemCount, err := s.itemRepo.CountByPantryID(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.GetPantryWithItemCount"), zap.Error(err), zap.Any("params", __logParams))

		itemCount = 0
	}
	result0 = &model.PantryWithItemCount{
		Pantry:    pantry,
		ItemCount: itemCount,
	}
	result1 = nil
	return
}

func (s *pantryService) ListPantriesByUser(ctx context.Context, userID uuid.UUID) (result0 []*model.Pantry, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.ListPantriesByUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.ListPantriesByUser"), zap.Any("params", __logParams))
	result0, result1 = s.repo.GetByUser(ctx, userID)
	if result1 != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.ListPantriesByUser"), zap.Error(result1), zap.Any("params", __logParams))
		result0 = nil
		return
	}
	return
}

func (s *pantryService) ListPantriesWithItemCount(ctx context.Context, userID uuid.UUID) (result0 []*model.PantryWithItemCount, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.ListPantriesWithItemCount"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.

			// If we can't get item count, default to 0 instead of failing
			Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.ListPantriesWithItemCount"), zap.Any("params", __logParams))
	pantries, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.ListPantriesWithItemCount"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	var result []*model.PantryWithItemCount
	for _, pantry := range pantries {
		itemCount, err := s.itemRepo.CountByPantryID(ctx, pantry.ID)
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "*pantryService.ListPantriesWithItemCount"), zap.Error(err), zap.Any("params", __logParams))

			itemCount = 0
		}
		result = append(result, &model.PantryWithItemCount{
			Pantry:    pantry,
			ItemCount: itemCount,
		})
	}
	result0 = result
	result1 = nil
	return
}

func (s *pantryService) UpdatePantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID, newName string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID, "newName": newName}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.UpdatePantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.UpdatePantry"), zap.Any("params", __logParams))
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, userID)
	if err != nil || !isOwner {
		result0 = ErrUnauthorized
		return
	}

	pantry, err := s.repo.GetByID(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.UpdatePantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = ErrPantryNotFound
		return
	}

	pantry.Name = newName
	pantry.UpdatedAt = time.Now()
	result0 = s.repo.Update(ctx, pantry)
	return
}

func (s *pantryService) DeletePantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.DeletePantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.DeletePantry"), zap.Any("params", __logParams))
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, userID)
	if err != nil || !isOwner {
		result0 = ErrUnauthorized
		return
	}
	result0 = s.repo.Delete(ctx, pantryID)
	return
}

func (s *pantryService) AddUserToPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "ownerID": ownerID, "targetUser": targetUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.AddUserToPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.AddUserToPantry"), zap.Any("params", __logParams))
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, ownerID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.AddUserToPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}
	if !isOwner {
		result0 = errors.New("only pantry owner can add users")
		return
	}

	user, err := s.userRepo.GetUserByEmail(ctx, targetUser)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.AddUserToPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = errors.New("user not found")
		return
	}

	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, user.ID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.AddUserToPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}
	if isMember {
		result0 = errors.New("user already in pantry")
		return
	}

	pantryUser := &model.PantryUser{
		PantryID: pantryID,
		UserID:   user.ID,
		Role:     "member",
	}
	result0 = s.repo.AddUserToPantry(ctx, pantryUser)
	return
}

func (s *pantryService) ListUsersInPantry(ctx context.Context, pantryID, userID uuid.UUID) (result0 []*model.PantryUserInfo, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.ListUsersInPantry"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.ListUsersInPantry"), zap.Any("params", __logParams))
	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.ListUsersInPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !isMember {
		result0 = nil
		result1 = errors.New("user is not in the pantry")
		return
	}

	users, err := s.repo.ListUsersInPantry(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.ListUsersInPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	for _, info := range users {
		user, err := s.userRepo.GetUserById(ctx, info.UserID)
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "*pantryService.ListUsersInPantry"), zap.Error(err), zap.Any("params", __logParams))
			continue
		}
		info.FirstName = user.FirstName
		info.LastName = user.LastName
	}
	result0 = users
	result1 = nil
	return
}

func (s *pantryService) RemoveUserFromPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID, "ownerID": ownerID, "targetUser": targetUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*pantryService.RemoveUserFromPantry"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*pantryService.RemoveUserFromPantry"), zap.Any("params", __logParams))
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, ownerID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.RemoveUserFromPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}
	if !isOwner {
		result0 = errors.New("only pantry owner can remove users")
		return
	}

	user, err := s.userRepo.GetUserByEmail(ctx, targetUser)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*pantryService.RemoveUserFromPantry"), zap.Error(err), zap.Any("params", __logParams))
		result0 = errors.New("user not found")
		return
	}

	if ownerID == user.ID {
		result0 = errors.New("owner cannot remove themselves")
		return
	}
	result0 = s.repo.RemoveUserFromPantry(ctx, pantryID, user.ID)
	return
}
