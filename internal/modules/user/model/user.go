package model

import "github.com/google/uuid"

type User struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	Role             string    `json:"role"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	ProfileCompleted bool      `json:"profile_completed"`
}
