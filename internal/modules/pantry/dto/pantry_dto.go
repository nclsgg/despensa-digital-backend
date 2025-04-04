package dto

import "github.com/google/uuid"

type CreatePantryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdatePantryRequest struct {
	Name string `json:"name" binding:"required"`
}

type PantryResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id"`
}

type ModifyPantryUserRequest struct {
	Email string `json:"email" binding:"required"`
}

type PantryUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}
