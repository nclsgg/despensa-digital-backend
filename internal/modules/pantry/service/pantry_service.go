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
) domain.PantryService {
	return &pantryService{
		repo:     repo,
		userRepo: userRepo,
		itemRepo: itemRepo,
	}
}

func (s *pantryService) CreatePantry(ctx context.Context, name string, ownerID uuid.UUID) (*model.Pantry, error) {
	pantry := &model.Pantry{
		ID:      uuid.New(),
		Name:    name,
		OwnerID: ownerID,
	}

	pantry, err := s.repo.Create(ctx, pantry)
	if err != nil {
		return nil, err
	}

	pantryUser := &model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantry.ID,
		UserID:   ownerID,
		Role:     "owner",
	}

	if err := s.repo.AddUserToPantry(ctx, pantryUser); err != nil {
		return nil, err
	}

	return pantry, nil
}

func (s *pantryService) GetPantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (*model.Pantry, error) {
	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		return nil, ErrUnauthorized
	}

	pantry, err := s.repo.GetByID(ctx, pantryID)
	if err != nil {
		return nil, ErrPantryNotFound
	}

	return pantry, nil
}

func (s *pantryService) GetPantryWithItemCount(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (*model.PantryWithItemCount, error) {
	pantry, err := s.GetPantry(ctx, pantryID, userID)
	if err != nil {
		return nil, err
	}

	itemCount, err := s.itemRepo.CountByPantryID(ctx, pantryID)
	if err != nil {
		// If we can't get item count, default to 0 instead of failing
		itemCount = 0
	}

	return &model.PantryWithItemCount{
		Pantry:    pantry,
		ItemCount: itemCount,
	}, nil
}

func (s *pantryService) ListPantriesByUser(ctx context.Context, userID uuid.UUID) ([]*model.Pantry, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *pantryService) ListPantriesWithItemCount(ctx context.Context, userID uuid.UUID) ([]*model.PantryWithItemCount, error) {
	pantries, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []*model.PantryWithItemCount
	for _, pantry := range pantries {
		itemCount, err := s.itemRepo.CountByPantryID(ctx, pantry.ID)
		if err != nil {
			// If we can't get item count, default to 0 instead of failing
			itemCount = 0
		}
		result = append(result, &model.PantryWithItemCount{
			Pantry:    pantry,
			ItemCount: itemCount,
		})
	}

	return result, nil
}

func (s *pantryService) UpdatePantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID, newName string) error {
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, userID)
	if err != nil || !isOwner {
		return ErrUnauthorized
	}

	pantry, err := s.repo.GetByID(ctx, pantryID)
	if err != nil {
		return ErrPantryNotFound
	}

	pantry.Name = newName
	pantry.UpdatedAt = time.Now()

	return s.repo.Update(ctx, pantry)
}

func (s *pantryService) DeletePantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) error {
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, userID)
	if err != nil || !isOwner {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, pantryID)
}

func (s *pantryService) AddUserToPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error {
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("only pantry owner can add users")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, targetUser)
	if err != nil {
		return errors.New("user not found")
	}

	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, user.ID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user already in pantry")
	}

	pantryUser := &model.PantryUser{
		PantryID: pantryID,
		UserID:   user.ID,
		Role:     "member",
	}

	return s.repo.AddUserToPantry(ctx, pantryUser)
}

func (s *pantryService) ListUsersInPantry(ctx context.Context, pantryID, userID uuid.UUID) ([]*model.PantryUserInfo, error) {
	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not in the pantry")
	}

	return s.repo.ListUsersInPantry(ctx, pantryID)
}

func (s *pantryService) RemoveUserFromPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error {
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("only pantry owner can remove users")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, targetUser)
	if err != nil {
		return errors.New("user not found")
	}

	if ownerID == user.ID {
		return errors.New("owner cannot remove themselves")
	}

	return s.repo.RemoveUserFromPantry(ctx, pantryID, user.ID)
}
