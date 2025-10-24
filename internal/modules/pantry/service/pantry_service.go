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
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
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
) domain.PantryService {
	return &pantryService{
		repo:     repo,
		userRepo: userRepo,
		itemRepo: itemRepo,
	}
}

func (s *pantryService) CreatePantry(ctx context.Context, name string, ownerID uuid.UUID) (*model.Pantry, error) {
	logger := appLogger.FromContext(ctx)

	pantry := &model.Pantry{
		ID:      uuid.New(),
		Name:    name,
		OwnerID: ownerID,
	}

	pantry, err := s.repo.Create(ctx, pantry)
	if err != nil {
		logger.Error("Failed to create pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "CreatePantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	pantryUser := &model.PantryUser{
		ID:       uuid.New(),
		PantryID: pantry.ID,
		UserID:   ownerID,
		Role:     "owner",
	}

	if err := s.repo.AddUserToPantry(ctx, pantryUser); err != nil {
		logger.Error("Failed to add owner to pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "CreatePantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantry.ID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Pantry created successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "CreatePantry"),
		zap.String(appLogger.FieldUserID, ownerID.String()),
		zap.String("pantry_id", pantry.ID.String()),
	)

	return pantry, nil
}

func (s *pantryService) GetPantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (*model.Pantry, error) {
	logger := appLogger.FromContext(ctx)

	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil || !isMember {
		logger.Warn("Unauthorized access attempt to pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, ErrUnauthorized
	}

	pantry, err := s.repo.GetByID(ctx, pantryID)
	if err != nil {
		logger.Error("Failed to get pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, ErrPantryNotFound
	}

	return pantry, nil
}

func (s *pantryService) GetPantryWithItemCount(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (*model.PantryWithItemCount, error) {
	logger := appLogger.FromContext(ctx)

	pantry, err := s.GetPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("Failed to get pantry with item count",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetPantryWithItemCount"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	// If we can't get item count, default to 0 instead of failing
	itemCount, err := s.itemRepo.CountByPantryID(ctx, pantryID)
	if err != nil {
		logger.Error("Failed to count pantry items",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetPantryWithItemCount"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		itemCount = 0
	}

	return &model.PantryWithItemCount{
		Pantry:    pantry,
		ItemCount: itemCount,
	}, nil
}

func (s *pantryService) GetMyPantry(ctx context.Context, userID uuid.UUID) (*model.PantryWithItemCount, error) {
	logger := appLogger.FromContext(ctx)

	// Get all pantries for the user
	pantries, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user pantries",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetMyPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	// Check if user has any pantries
	if len(pantries) == 0 {
		return nil, ErrPantryNotFound
	}

	// Return the first pantry with item count
	firstPantry := pantries[0]
	itemCount, err := s.itemRepo.CountByPantryID(ctx, firstPantry.ID)
	if err != nil {
		logger.Error("Failed to count items in user's main pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetMyPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", firstPantry.ID.String()),
			zap.Error(err),
		)
		// Default to 0 if we can't get the count
		itemCount = 0
	}

	logger.Info("Retrieved user's main pantry",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "GetMyPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", firstPantry.ID.String()),
		zap.Int(appLogger.FieldCount, itemCount),
	)

	return &model.PantryWithItemCount{
		Pantry:    firstPantry,
		ItemCount: itemCount,
	}, nil
}

func (s *pantryService) ListPantriesByUser(ctx context.Context, userID uuid.UUID) ([]*model.Pantry, error) {
	logger := appLogger.FromContext(ctx)

	pantries, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		logger.Error("Failed to list user pantries",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListPantriesByUser"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Listed user pantries",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "ListPantriesByUser"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(pantries)),
	)

	return pantries, nil
}

func (s *pantryService) ListPantriesWithItemCount(ctx context.Context, userID uuid.UUID) ([]*model.PantryWithItemCount, error) {
	logger := appLogger.FromContext(ctx)

	pantries, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		logger.Error("Failed to list pantries with item count",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListPantriesWithItemCount"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	var result []*model.PantryWithItemCount
	for _, pantry := range pantries {
		itemCount, err := s.itemRepo.CountByPantryID(ctx, pantry.ID)
		if err != nil {
			logger.Error("Failed to count items for pantry",
				zap.String(appLogger.FieldModule, "pantry"),
				zap.String(appLogger.FieldFunction, "ListPantriesWithItemCount"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantry.ID.String()),
				zap.Error(err),
			)
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
	logger := appLogger.FromContext(ctx)

	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, userID)
	if err != nil || !isOwner {
		logger.Warn("Unauthorized pantry update attempt",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "UpdatePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return ErrUnauthorized
	}

	pantry, err := s.repo.GetByID(ctx, pantryID)
	if err != nil {
		logger.Error("Failed to get pantry for update",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "UpdatePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return ErrPantryNotFound
	}

	pantry.Name = newName
	pantry.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, pantry); err != nil {
		logger.Error("Failed to update pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "UpdatePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("Pantry updated successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "UpdatePantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
	)

	return nil
}

func (s *pantryService) DeletePantry(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, userID)
	if err != nil || !isOwner {
		logger.Warn("Unauthorized pantry deletion attempt",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "DeletePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return ErrUnauthorized
	}

	if err := s.repo.Delete(ctx, pantryID); err != nil {
		logger.Error("Failed to delete pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "DeletePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("Pantry deleted successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "DeletePantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
	)

	return nil
}

func (s *pantryService) AddUserToPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error {
	logger := appLogger.FromContext(ctx)

	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, ownerID)
	if err != nil {
		logger.Error("Failed to verify pantry ownership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return err
	}
	if !isOwner {
		logger.Warn("Non-owner attempted to add user to pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return errors.New("only pantry owner can add users")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, targetUser)
	if err != nil {
		logger.Error("Target user not found",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(targetUser)),
			zap.Error(err),
		)
		return errors.New("user not found")
	}

	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, user.ID)
	if err != nil {
		logger.Error("Failed to check pantry membership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("target_user_id", user.ID.String()),
			zap.Error(err),
		)
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

	if err := s.repo.AddUserToPantry(ctx, pantryUser); err != nil {
		logger.Error("Failed to add user to pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("target_user_id", user.ID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("User added to pantry successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "AddUserToPantry"),
		zap.String(appLogger.FieldUserID, ownerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String("target_user_id", user.ID.String()),
	)

	return nil
}

func (s *pantryService) ListUsersInPantry(ctx context.Context, pantryID, userID uuid.UUID) ([]*model.PantryUserInfo, error) {
	logger := appLogger.FromContext(ctx)

	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, userID)
	if err != nil {
		logger.Error("Failed to verify pantry membership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}
	if !isMember {
		logger.Warn("Non-member attempted to list pantry users",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, errors.New("user is not in the pantry")
	}

	users, err := s.repo.ListUsersInPantry(ctx, pantryID)
	if err != nil {
		logger.Error("Failed to list pantry users",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	for _, info := range users {
		user, err := s.userRepo.GetUserById(ctx, info.UserID)
		if err != nil {
			logger.Error("Failed to get user details",
				zap.String(appLogger.FieldModule, "pantry"),
				zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
				zap.String("target_user_id", info.UserID.String()),
				zap.Error(err),
			)
			continue
		}
		info.FirstName = user.FirstName
		info.LastName = user.LastName
	}

	logger.Info("Listed pantry users",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.Int(appLogger.FieldCount, len(users)),
	)

	return users, nil
}

func (s *pantryService) RemoveUserFromPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error {
	logger := appLogger.FromContext(ctx)

	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, ownerID)
	if err != nil {
		logger.Error("Failed to verify pantry ownership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return err
	}
	if !isOwner {
		logger.Warn("Non-owner attempted to remove user from pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return errors.New("only pantry owner can remove users")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, targetUser)
	if err != nil {
		logger.Error("Target user not found",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(targetUser)),
			zap.Error(err),
		)
		return errors.New("user not found")
	}

	if ownerID == user.ID {
		return errors.New("owner cannot remove themselves")
	}

	if err := s.repo.RemoveUserFromPantry(ctx, pantryID, user.ID); err != nil {
		logger.Error("Failed to remove user from pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("target_user_id", user.ID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("User removed from pantry successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
		zap.String(appLogger.FieldUserID, ownerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String("target_user_id", user.ID.String()),
	)

	return nil
}

func (s *pantryService) RemoveSpecificUserFromPantry(ctx context.Context, pantryID, ownerID, targetUserID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	// Verify that the requester is the owner
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, ownerID)
	if err != nil {
		logger.Error("Failed to verify pantry ownership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return err
	}
	if !isOwner {
		logger.Warn("Non-owner attempted to remove specific user from pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return errors.New("only pantry owner can remove users")
	}

	// Owner cannot remove themselves
	if ownerID == targetUserID {
		return errors.New("owner cannot remove themselves")
	}

	// Verify that the target user is in the pantry
	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, targetUserID)
	if err != nil {
		logger.Error("Failed to verify pantry membership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("target_user_id", targetUserID.String()),
			zap.Error(err),
		)
		return err
	}
	if !isMember {
		return errors.New("user is not in the pantry")
	}

	if err := s.repo.RemoveUserFromPantry(ctx, pantryID, targetUserID); err != nil {
		logger.Error("Failed to remove specific user from pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("target_user_id", targetUserID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("Specific user removed from pantry successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
		zap.String(appLogger.FieldUserID, ownerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String("target_user_id", targetUserID.String()),
	)

	return nil
}

func (s *pantryService) TransferOwnership(ctx context.Context, pantryID, currentOwnerID, newOwnerID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	// Verify that the requester is the current owner
	isOwner, err := s.repo.IsUserOwner(ctx, pantryID, currentOwnerID)
	if err != nil {
		logger.Error("Failed to verify current ownership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return err
	}
	if !isOwner {
		logger.Warn("Non-owner attempted ownership transfer",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return errors.New("only pantry owner can transfer ownership")
	}

	// Cannot transfer to self
	if currentOwnerID == newOwnerID {
		return errors.New("cannot transfer ownership to yourself")
	}

	// Verify that the new owner is a member of the pantry
	isMember, err := s.repo.IsUserInPantry(ctx, pantryID, newOwnerID)
	if err != nil {
		logger.Error("Failed to verify new owner membership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("new_owner_id", newOwnerID.String()),
			zap.Error(err),
		)
		return err
	}
	if !isMember {
		return errors.New("new owner must be a member of the pantry")
	}

	// Update pantry owner_id
	pantry, err := s.repo.GetByID(ctx, pantryID)
	if err != nil {
		logger.Error("Failed to get pantry for ownership transfer",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return ErrPantryNotFound
	}

	pantry.OwnerID = newOwnerID
	pantry.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, pantry); err != nil {
		logger.Error("Failed to update pantry ownership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("new_owner_id", newOwnerID.String()),
			zap.Error(err),
		)
		return err
	}

	// Update current owner role to member
	if err := s.repo.UpdatePantryUserRole(ctx, pantryID, currentOwnerID, "member"); err != nil {
		logger.Error("Failed to update previous owner role",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return err
	}

	// Update new owner role to owner
	if err := s.repo.UpdatePantryUserRole(ctx, pantryID, newOwnerID, "owner"); err != nil {
		logger.Error("Failed to update new owner role",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("new_owner_id", newOwnerID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("Pantry ownership transferred successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "TransferOwnership"),
		zap.String(appLogger.FieldUserID, currentOwnerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String("new_owner_id", newOwnerID.String()),
	)

	return nil
}
