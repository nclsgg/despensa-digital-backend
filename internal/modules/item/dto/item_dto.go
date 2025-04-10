package dto

import "time"

type CreateItemDTO struct {
	PantryID     string  `json:"pantry_id" binding:"required,uuid"`
	Name         string  `json:"name" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gte=0"`
	PricePerUnit float64 `json:"price_per_unit" binding:"required,gte=0"`
	Unit         string  `json:"unit" binding:"required"`
	CategoryID   *string `json:"category_id,omitempty"`
	ExpiresAt    string  `json:"expires_at,omitempty"`
}

type UpdateItemDTO struct {
	Name         *string  `json:"name,omitempty"`
	Quantity     *float64 `json:"quantity,omitempty"`
	PricePerUnit *float64 `json:"price_per_unit,omitempty"`
	Unit         *string  `json:"unit,omitempty"`
	CategoryID   *string  `json:"category_id,omitempty"`
	ExpiresAt    string   `json:"expires_at,omitempty"`
}

type ItemResponse struct {
	ID           string     `json:"id"`
	PantryID     string     `json:"pantry_id"`
	AddedBy      string     `json:"added_by"`
	Name         string     `json:"name"`
	Quantity     float64    `json:"quantity"`
	TotalPrice   float64    `json:"total_price"`
	PricePerUnit float64    `json:"price_per_unit"`
	Unit         string     `json:"unit"`
	CategoryID   string     `json:"category_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}
