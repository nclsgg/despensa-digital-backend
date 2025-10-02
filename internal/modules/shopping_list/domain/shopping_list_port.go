package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
)

type ShoppingListRepository interface {
	Create(ctx context.Context, shoppingList *model.ShoppingList) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.ShoppingList, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.ShoppingList, error)
	Update(ctx context.Context, shoppingList *model.ShoppingList) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateItem(ctx context.Context, item *model.ShoppingListItem) error
	UpdateItem(ctx context.Context, item *model.ShoppingListItem) error
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
	GetItemsByShoppingListID(ctx context.Context, shoppingListID uuid.UUID) ([]*model.ShoppingListItem, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type ShoppingListService interface {
	CreateShoppingList(ctx context.Context, userID uuid.UUID, input dto.CreateShoppingListDTO) (*dto.ShoppingListResponseDTO, error)
	GetShoppingListByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.ShoppingListResponseDTO, error)
	GetShoppingListsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*dto.ShoppingListSummaryDTO, error)
	UpdateShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID, input dto.UpdateShoppingListDTO) (*dto.ShoppingListResponseDTO, error)
	DeleteShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	CreateShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, input dto.CreateShoppingListItemDTO) (*dto.ShoppingListResponseDTO, error)
	UpdateShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID, input dto.UpdateShoppingListItemDTO) (*dto.ShoppingListItemResponseDTO, error)
	DeleteShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID) error
	GenerateAIShoppingList(ctx context.Context, userID uuid.UUID, input dto.GenerateAIShoppingListDTO) (*dto.ShoppingListResponseDTO, error)
}
