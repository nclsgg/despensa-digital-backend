package dto

type CreateProfileDTO struct {
	MonthlyIncome       float64  `json:"monthly_income" binding:"required,min=0"`
	PreferredBudget     float64  `json:"preferred_budget" binding:"required,min=0"`
	HouseholdSize       int      `json:"household_size" binding:"required,min=1"`
	DietaryRestrictions []string `json:"dietary_restrictions"`
	PreferredBrands     []string `json:"preferred_brands"`
	ShoppingFrequency   string   `json:"shopping_frequency" binding:"required,oneof=weekly biweekly monthly"`
}

type UpdateProfileDTO struct {
	MonthlyIncome       *float64  `json:"monthly_income,omitempty" binding:"omitempty,min=0"`
	PreferredBudget     *float64  `json:"preferred_budget,omitempty" binding:"omitempty,min=0"`
	HouseholdSize       *int      `json:"household_size,omitempty" binding:"omitempty,min=1"`
	DietaryRestrictions *[]string `json:"dietary_restrictions,omitempty"`
	PreferredBrands     *[]string `json:"preferred_brands,omitempty"`
	ShoppingFrequency   *string   `json:"shopping_frequency,omitempty" binding:"omitempty,oneof=weekly biweekly monthly"`
}

type ProfileResponseDTO struct {
	ID                  string   `json:"id"`
	UserID              string   `json:"user_id"`
	MonthlyIncome       float64  `json:"monthly_income"`
	PreferredBudget     float64  `json:"preferred_budget"`
	HouseholdSize       int      `json:"household_size"`
	DietaryRestrictions []string `json:"dietary_restrictions"`
	PreferredBrands     []string `json:"preferred_brands"`
	ShoppingFrequency   string   `json:"shopping_frequency"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
}
