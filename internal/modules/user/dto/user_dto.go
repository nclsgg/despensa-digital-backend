package dto

type UserResponse struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Role             string `json:"role"`
	ProfileCompleted bool   `json:"profile_completed"`
}

type CompleteProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
