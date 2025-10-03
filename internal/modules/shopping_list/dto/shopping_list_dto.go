package dto

import "github.com/google/uuid"

type ShoppingListPreferencesOverrideDTO struct {
	HouseholdSize       *int     `json:"household_size,omitempty" binding:"omitempty,min=1,max=20"`
	MonthlyIncome       *float64 `json:"monthly_income,omitempty" binding:"omitempty,min=0"`
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
}

type ShoppingListPreferencesDTO struct {
	HouseholdSize       int      `json:"household_size"`
	MonthlyIncome       float64  `json:"monthly_income"`
	DietaryRestrictions []string `json:"dietary_restrictions"`
}

type CreateShoppingListDTO struct {
	Name        string                              `json:"name" binding:"required"`
	PantryID    *uuid.UUID                          `json:"pantry_id,omitempty"`
	TotalBudget float64                             `json:"total_budget" binding:"required,min=0"`
	Items       []CreateShoppingListItemDTO         `json:"items"`
	Preferences *ShoppingListPreferencesOverrideDTO `json:"preferences,omitempty"`
}

type CreateShoppingListItemDTO struct {
	Name           string     `json:"name" binding:"required"`
	Quantity       float64    `json:"quantity" binding:"required,min=0.1"`
	Unit           string     `json:"unit" binding:"required"`
	EstimatedPrice float64    `json:"estimated_price" binding:"required,min=0"`
	Category       string     `json:"category"`
	Priority       int        `json:"priority" binding:"omitempty,min=1,max=3"`
	PantryItemID   *uuid.UUID `json:"pantry_item_id,omitempty"`
}

type UpdateShoppingListDTO struct {
	Name        *string                             `json:"name,omitempty"`
	Status      *string                             `json:"status,omitempty" binding:"omitempty,oneof=pending completed cancelled"`
	TotalBudget *float64                            `json:"total_budget,omitempty" binding:"omitempty,min=0"`
	ActualCost  *float64                            `json:"actual_cost,omitempty" binding:"omitempty,min=0"`
	Preferences *ShoppingListPreferencesOverrideDTO `json:"preferences,omitempty"`
}

type UpdateShoppingListItemDTO struct {
	Name           *string    `json:"name,omitempty"`
	Quantity       *float64   `json:"quantity,omitempty" binding:"omitempty,min=0.1"`
	Unit           *string    `json:"unit,omitempty"`
	EstimatedPrice *float64   `json:"estimated_price,omitempty" binding:"omitempty,min=0"`
	ActualPrice    *float64   `json:"actual_price,omitempty" binding:"omitempty,min=0"`
	Category       *string    `json:"category,omitempty"`
	Priority       *int       `json:"priority,omitempty" binding:"omitempty,min=1,max=3"`
	Purchased      *bool      `json:"purchased,omitempty"`
	PantryItemID   *uuid.UUID `json:"pantry_item_id,omitempty"`
}

type GenerateAIShoppingListDTO struct {
	Name          string                              `json:"name" binding:"required"`
	PantryID      uuid.UUID                           `json:"pantry_id" binding:"required"`
	Prompt        string                              `json:"prompt"`
	MaxBudget     *float64                            `json:"max_budget,omitempty" binding:"omitempty,min=0"`
	PeopleCount   *int                                `json:"people_count,omitempty" binding:"omitempty,min=1,max=20"`
	ShoppingType  string                              `json:"shopping_type,omitempty" binding:"omitempty,oneof=weekly monthly stock_up emergency"`
	IncludeBasics *bool                               `json:"include_basics,omitempty"`
	ExcludeItems  []string                            `json:"exclude_items,omitempty"`
	Notes         string                              `json:"notes,omitempty"`
	Preferences   *ShoppingListPreferencesOverrideDTO `json:"preferences,omitempty"`
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
	Preferences   ShoppingListPreferencesDTO    `json:"preferences"`
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
	Priority       int     `json:"priority"`
	Purchased      bool    `json:"purchased"`
	Source         string  `json:"source"`
	PantryItemID   *string `json:"pantry_item_id,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type ShoppingListSummaryDTO struct {
	ID             string                     `json:"id"`
	PantryID       *string                    `json:"pantry_id,omitempty"`
	PantryName     string                     `json:"pantry_name,omitempty"`
	Name           string                     `json:"name"`
	Status         string                     `json:"status"`
	TotalBudget    float64                    `json:"total_budget"`
	EstimatedCost  float64                    `json:"estimated_cost"`
	ActualCost     float64                    `json:"actual_cost"`
	GeneratedBy    string                     `json:"generated_by"`
	ItemCount      int                        `json:"item_count"`
	PurchasedCount int                        `json:"purchased_count"`
	Preferences    ShoppingListPreferencesDTO `json:"preferences"`
	CreatedAt      string                     `json:"created_at"`
	UpdatedAt      string                     `json:"updated_at"`
}
