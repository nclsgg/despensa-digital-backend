package dto

import "time"

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
	ID        string     `json:"id"`
	PantryID  string     `json:"pantry_id"`
	AddedBy   string     `json:"added_by"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
