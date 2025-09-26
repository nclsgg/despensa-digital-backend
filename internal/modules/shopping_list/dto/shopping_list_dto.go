package dto

import "github.com/google/uuid"

type CreateShoppingListDTO struct {
	Name        string                      `json:"name" binding:"required"`
	PantryID    *uuid.UUID                  `json:"pantry_id,omitempty"`
	TotalBudget float64                     `json:"total_budget" binding:"required,min=0"`
	Items       []CreateShoppingListItemDTO `json:"items"`
}

type CreateShoppingListItemDTO struct {
	Name           string  `json:"name" binding:"required"`
	Quantity       float64 `json:"quantity" binding:"required,min=0.1"`
	Unit           string  `json:"unit" binding:"required"`
	EstimatedPrice float64 `json:"estimated_price" binding:"required,min=0"`
	Category       string  `json:"category"`
	Brand          string  `json:"brand"`
	Priority       int     `json:"priority" binding:"omitempty,min=1,max=3"`
	Notes          string  `json:"notes"`
}

type UpdateShoppingListDTO struct {
	Name        *string  `json:"name,omitempty"`
	Status      *string  `json:"status,omitempty" binding:"omitempty,oneof=pending completed cancelled"`
	TotalBudget *float64 `json:"total_budget,omitempty" binding:"omitempty,min=0"`
	ActualCost  *float64 `json:"actual_cost,omitempty" binding:"omitempty,min=0"`
}

type UpdateShoppingListItemDTO struct {
	Name        *string  `json:"name,omitempty"`
	Quantity    *float64 `json:"quantity,omitempty" binding:"omitempty,min=0.1"`
	Unit        *string  `json:"unit,omitempty"`
	ActualPrice *float64 `json:"actual_price,omitempty" binding:"omitempty,min=0"`
	Category    *string  `json:"category,omitempty"`
	Brand       *string  `json:"brand,omitempty"`
	Priority    *int     `json:"priority,omitempty" binding:"omitempty,min=1,max=3"`
	Purchased   *bool    `json:"purchased,omitempty"`
	Notes       *string  `json:"notes,omitempty"`
}

type GenerateAIShoppingListDTO struct {
	Name            string    `json:"name" binding:"required"`
	PantryID        uuid.UUID `json:"pantry_id" binding:"required"`
	Prompt          string    `json:"prompt"`
	MaxBudget       *float64  `json:"max_budget,omitempty" binding:"omitempty,min=0"`
	PeopleCount     *int      `json:"people_count,omitempty" binding:"omitempty,min=1,max=20"`
	ShoppingType    string    `json:"shopping_type,omitempty" binding:"omitempty,oneof=weekly monthly stock_up emergency"`
	IncludeBasics   *bool     `json:"include_basics,omitempty"`
	ExcludeItems    []string  `json:"exclude_items,omitempty"`
	PreferredBrands []string  `json:"preferred_brands,omitempty"`
	Notes           string    `json:"notes,omitempty"`
}

type ShoppingListResponseDTO struct {
	ID            string                        `json:"id"`
	UserID        string                        `json:"user_id"`
	PantryID      *string                       `json:"pantry_id,omitempty"`
	PantryName    string                        `json:"pantry_name,omitempty"`
	Name          string                        `json:"name"`
	Status        string                        `json:"status"`
	TotalBudget   float64                       `json:"total_budget"`
	EstimatedCost float64                       `json:"estimated_cost"`
	ActualCost    float64                       `json:"actual_cost"`
	GeneratedBy   string                        `json:"generated_by"`
	Items         []ShoppingListItemResponseDTO `json:"items"`
	CreatedAt     string                        `json:"created_at"`
	UpdatedAt     string                        `json:"updated_at"`
}

type ShoppingListItemResponseDTO struct {
	ID             string  `json:"id"`
	ShoppingListID string  `json:"shopping_list_id"`
	Name           string  `json:"name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	EstimatedPrice float64 `json:"estimated_price"`
	ActualPrice    float64 `json:"actual_price"`
	Category       string  `json:"category"`
	Brand          string  `json:"brand"`
	Priority       int     `json:"priority"`
	Purchased      bool    `json:"purchased"`
	Notes          string  `json:"notes"`
	Source         string  `json:"source"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type ShoppingListSummaryDTO struct {
	ID             string  `json:"id"`
	PantryID       *string `json:"pantry_id,omitempty"`
	PantryName     string  `json:"pantry_name,omitempty"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	TotalBudget    float64 `json:"total_budget"`
	EstimatedCost  float64 `json:"estimated_cost"`
	ActualCost     float64 `json:"actual_cost"`
	GeneratedBy    string  `json:"generated_by"`
	ItemCount      int     `json:"item_count"`
	PurchasedCount int     `json:"purchased_count"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
