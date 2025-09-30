package dto

type CreateItemDTO struct {
	PantryID      string  `json:"pantry_id" binding:"required,uuid"`
	Name          string  `json:"name" binding:"required"`
	Quantity      float64 `json:"quantity" binding:"required,gte=0"`
	PricePerUnit  float64 `json:"price_per_unit" binding:"required,gte=0"`
	PriceQuantity float64 `json:"price_quantity" binding:"omitempty,min=0.001"`
	Unit          string  `json:"unit" binding:"required"`
	CategoryID    *string `json:"category_id,omitempty"`
	ExpiresAt     string  `json:"expires_at,omitempty"`
}

type UpdateItemDTO struct {
	Name          *string  `json:"name,omitempty"`
	Quantity      *float64 `json:"quantity,omitempty"`
	PricePerUnit  *float64 `json:"price_per_unit,omitempty"`
	PriceQuantity *float64 `json:"price_quantity,omitempty"`
	Unit          *string  `json:"unit,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	ExpiresAt     string   `json:"expires_at,omitempty"`
}

type ItemFilterDTO struct {
	MinPrice      *float64 `json:"min_price,omitempty"`
	MaxPrice      *float64 `json:"max_price,omitempty"`
	ExpiresUntil  string   `json:"expires_until,omitempty"`
	Name          *string  `json:"name,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	SortBy        *string  `json:"sort_by,omitempty"`        // "price", "expires_at", "category", "name"
	SortDirection *string  `json:"sort_direction,omitempty"` // "asc", "desc"
}

type ItemResponse struct {
	ID            string  `json:"id"`
	PantryID      string  `json:"pantry_id"`
	AddedBy       string  `json:"added_by"`
	Name          string  `json:"name"`
	Quantity      float64 `json:"quantity"`
	Unit          string  `json:"unit"`
	PricePerUnit  float64 `json:"price_per_unit"`
	PriceQuantity float64 `json:"price_quantity"`
	TotalPrice    float64 `json:"total_price"`
	CategoryID    *string `json:"category_id,omitempty"`
	ExpiresAt     *string `json:"expires_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
