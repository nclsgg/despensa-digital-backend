package dto

type CreateItemCategoryDTO struct {
	PantryID string `json:"pantry_id" binding:"required,uuid"`
	Name     string `json:"name" binding:"required"`
	Color    string `json:"color" binding:"required"`
}

type CreateDefaultItemCategoryDTO struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type UpdateItemCategoryDTO struct {
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

type ItemCategoryResponse struct {
	ID        string  `json:"id"`
	PantryID  string  `json:"pantry_id"`
	AddedBy   string  `json:"added_by"`
	Name      string  `json:"name"`
	Color     string  `json:"color"`
	IsDefault bool    `json:"is_default"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}
