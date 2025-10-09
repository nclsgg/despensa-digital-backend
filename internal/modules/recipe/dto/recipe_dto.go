package dto

import "time"

// AvailableIngredientDTO representa um ingrediente disponível em uma despensa
// exposto para o frontend.
type AvailableIngredientDTO struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// SaveRecipeDTO represents a recipe to be saved
type SaveRecipeDTO struct {
	ID                  string                     `json:"id" validate:"required,uuid"`
	Title               string                     `json:"title" validate:"required,min=1,max=255"`
	Description         string                     `json:"description"`
	Ingredients         []SaveRecipeIngredientDTO  `json:"ingredients" validate:"required,min=1,dive"`
	Instructions        []SaveRecipeInstructionDTO `json:"instructions" validate:"required,min=1,dive"`
	CookingTime         *int                       `json:"cooking_time"`
	PreparationTime     *int                       `json:"preparation_time"`
	TotalTime           *int                       `json:"total_time"`
	ServingSize         *int                       `json:"serving_size"`
	Difficulty          string                     `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	MealType            string                     `json:"meal_type" validate:"omitempty,oneof=breakfast lunch dinner snack dessert"`
	Cuisine             string                     `json:"cuisine" validate:"max=100"`
	DietaryRestrictions []string                   `json:"dietary_restrictions"`
	NutritionInfo       SaveRecipeNutritionDTO     `json:"nutrition_info"`
	Tips                []string                   `json:"tips"`
	GeneratedAt         string                     `json:"generated_at"`
}

// SaveRecipeIngredientDTO represents an ingredient in a recipe to be saved
type SaveRecipeIngredientDTO struct {
	Name        string   `json:"name" validate:"required"`
	Amount      *float64 `json:"amount"`
	Unit        string   `json:"unit" validate:"required"`
	Available   bool     `json:"available"`
	Alternative *string  `json:"alternative"`
}

// SaveRecipeInstructionDTO represents an instruction in a recipe to be saved
type SaveRecipeInstructionDTO struct {
	Step        int    `json:"step" validate:"required,min=1"`
	Description string `json:"description" validate:"required"`
	Time        *int   `json:"time"`
}

// SaveRecipeNutritionDTO represents nutrition info in a recipe to be saved
type SaveRecipeNutritionDTO struct {
	Calories      *int `json:"calories"`
	Protein       *int `json:"protein"`
	Carbohydrates *int `json:"carbohydrates"`
	Fat           *int `json:"fat"`
}

// RecipeDetailDTO represents a saved recipe returned to the client
type RecipeDetailDTO struct {
	ID                  string                       `json:"id"`
	Title               string                       `json:"title"`
	Description         string                       `json:"description"`
	Ingredients         []RecipeIngredientDetailDTO  `json:"ingredients"`
	Instructions        []RecipeInstructionDetailDTO `json:"instructions"`
	CookingTime         *int                         `json:"cooking_time"`
	PreparationTime     *int                         `json:"preparation_time"`
	TotalTime           *int                         `json:"total_time"`
	ServingSize         *int                         `json:"serving_size"`
	Difficulty          string                       `json:"difficulty"`
	MealType            string                       `json:"meal_type"`
	Cuisine             string                       `json:"cuisine"`
	DietaryRestrictions []string                     `json:"dietary_restrictions"`
	NutritionInfo       RecipeNutritionDetailDTO     `json:"nutrition_info"`
	Tips                []string                     `json:"tips"`
	GeneratedAt         time.Time                    `json:"generated_at"`
	CreatedAt           time.Time                    `json:"created_at"`
}

// RecipeIngredientDetailDTO represents an ingredient in a saved recipe
type RecipeIngredientDetailDTO struct {
	Name        string   `json:"name"`
	Amount      *float64 `json:"amount"`
	Unit        string   `json:"unit"`
	Available   bool     `json:"available"`
	Alternative *string  `json:"alternative,omitempty"`
}

// RecipeInstructionDetailDTO represents an instruction in a saved recipe
type RecipeInstructionDetailDTO struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
	Time        *int   `json:"time,omitempty"`
}

// RecipeNutritionDetailDTO represents nutrition info in a saved recipe
type RecipeNutritionDetailDTO struct {
	Calories      *int `json:"calories,omitempty"`
	Protein       *int `json:"protein,omitempty"`
	Carbohydrates *int `json:"carbohydrates,omitempty"`
	Fat           *int `json:"fat,omitempty"`
}
