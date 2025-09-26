package dto

type CreatePantryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdatePantryRequest struct {
	Name string `json:"name" binding:"required"`
}

type PantrySummaryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	ItemCount int    `json:"item_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PantryDetailResponse struct {
	PantrySummaryResponse
	Items []PantryItemResponse `json:"items"`
}

type PantryItemResponse struct {
	ID             string  `json:"id"`
	PantryID       string  `json:"pantry_id"`
	Name           string  `json:"name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	PricePerUnit   float64 `json:"price_per_unit"`
	TotalPrice     float64 `json:"total_price"`
	AddedBy        string  `json:"added_by"`
	CategoryID     *string `json:"category_id,omitempty"`
	ExpirationDate *string `json:"expiration_date,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type ModifyPantryUserRequest struct {
	Email string `json:"email" binding:"required"`
}

type PantryUserResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	PantryID string `json:"pantry_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}
