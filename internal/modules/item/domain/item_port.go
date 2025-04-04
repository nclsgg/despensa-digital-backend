package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
)

type ItemService interface {
	Create(ctx context.Context, dto dto.CreateItemDTO, userID uuid.UUID) (*model.Item, error)
	Update(ctx context.Context, id uuid.UUID, dto dto.UpdateItemDTO, userID uuid.UUID) (*model.Item, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Item, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]*model.Item, error)
}

type ItemRepository interface {
	Create(ctx context.Context, item *model.Item) error
	Update(ctx context.Context, item *model.Item) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByPantryID(ctx context.Context, pantryID uuid.UUID) ([]*model.Item, error)
}

type ItemCategoryService interface {
	Create(ctx context.Context, dto dto.CreateItemCategoryDTO, userID uuid.UUID) (*model.ItemCategory, error)
	CreateDefault(ctx context.Context, dto dto.CreateDefaultItemCategoryDTO, userID uuid.UUID) (*model.ItemCategory, error)
	CloneDefaultCategoryToPantry(ctx context.Context, defaultCategoryID, pantryID uuid.UUID, userID uuid.UUID) (*model.ItemCategory, error)
	Update(ctx context.Context, id uuid.UUID, dto dto.UpdateItemCategoryDTO, userID uuid.UUID) (*model.ItemCategory, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.ItemCategory, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ListByPantryID(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]*model.ItemCategory, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ItemCategory, error)
}

type ItemCategoryRepository interface {
	Create(ctx context.Context, itemCategory *model.ItemCategory) error
	Update(ctx context.Context, itemCategory *model.ItemCategory) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.ItemCategory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByPantryID(ctx context.Context, pantryID uuid.UUID) ([]*model.ItemCategory, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ItemCategory, error)
}

type ItemHandler interface {
	CreateItem(ctx *gin.Context)
	UpdateItem(ctx *gin.Context)
	GetItem(ctx *gin.Context)
	DeleteItem(ctx *gin.Context)
	ListItems(ctx *gin.Context)
}

type ItemCategoryHandler interface {
	CreateItemCategory(ctx *gin.Context)
	CreateDefaultItemCategory(ctx *gin.Context)
	CloneDefaultCategoryToPantry(ctx *gin.Context)
	UpdateItemCategory(ctx *gin.Context)
	GetItemCategory(ctx *gin.Context)
	DeleteItemCategory(ctx *gin.Context)
	ListItemCategoriesByPantry(ctx *gin.Context)
	ListItemCategoriesByUser(ctx *gin.Context)
}
