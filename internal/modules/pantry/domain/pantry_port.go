package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
)

type PantryService interface {
	CreatePantry(ctx context.Context, name string, ownerID uuid.UUID) (*model.Pantry, error)
	GetPantry(ctx context.Context, pantryID, userID uuid.UUID) (*model.Pantry, error)
	GetPantryWithItemCount(ctx context.Context, pantryID, userID uuid.UUID) (*model.PantryWithItemCount, error)
	ListPantriesByUser(ctx context.Context, userID uuid.UUID) ([]*model.Pantry, error)
	ListPantriesWithItemCount(ctx context.Context, userID uuid.UUID) ([]*model.PantryWithItemCount, error)
	DeletePantry(ctx context.Context, pantryID, userID uuid.UUID) error
	UpdatePantry(ctx context.Context, pantryID, userID uuid.UUID, newName string) error
	AddUserToPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error
	RemoveUserFromPantry(ctx context.Context, pantryID, ownerID uuid.UUID, targetUser string) error
	RemoveSpecificUserFromPantry(ctx context.Context, pantryID, ownerID, targetUserID uuid.UUID) error
	TransferOwnership(ctx context.Context, pantryID, currentOwnerID, newOwnerID uuid.UUID) error
	ListUsersInPantry(ctx context.Context, pantryID, userID uuid.UUID) ([]*model.PantryUserInfo, error)
}

type PantryRepository interface {
	Create(ctx context.Context, pantry *model.Pantry) (*model.Pantry, error)
	Delete(ctx context.Context, pantryID uuid.UUID) error
	Update(ctx context.Context, pantry *model.Pantry) error
	GetByID(ctx context.Context, pantryID uuid.UUID) (*model.Pantry, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*model.Pantry, error)
	IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (bool, error)
	IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (bool, error)
	AddUserToPantry(ctx context.Context, pantryUser *model.PantryUser) error
	RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) error
	UpdatePantryUserRole(ctx context.Context, pantryID, userID uuid.UUID, newRole string) error
	GetPantryUser(ctx context.Context, pantryID, userID uuid.UUID) (*model.PantryUser, error)
	ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) ([]*model.PantryUserInfo, error)
}

type PantryHandler interface {
	CreatePantry(ctx *gin.Context)
	ListPantries(ctx *gin.Context)
	GetPantry(ctx *gin.Context)
	DeletePantry(ctx *gin.Context)
	UpdatePantry(ctx *gin.Context)
	AddUserToPantry(ctx *gin.Context)
	RemoveUserFromPantry(ctx *gin.Context)
	RemoveSpecificUserFromPantry(ctx *gin.Context)
	TransferOwnership(ctx *gin.Context)
	ListUsersInPantry(c *gin.Context)
}
