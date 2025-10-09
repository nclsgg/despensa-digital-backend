package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Recipe represents a saved recipe in the database
type Recipe struct {
	ID                  uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID              `gorm:"type:uuid;not null;index" json:"user_id"`
	Title               string                 `gorm:"type:varchar(255);not null" json:"title"`
	Description         string                 `gorm:"type:text" json:"description"`
	Ingredients         RecipeIngredientsJSON  `gorm:"type:jsonb;not null" json:"ingredients"`
	Instructions        RecipeInstructionsJSON `gorm:"type:jsonb;not null" json:"instructions"`
	CookingTime         *int                   `gorm:"type:int" json:"cooking_time"`
	PreparationTime     *int                   `gorm:"type:int" json:"preparation_time"`
	TotalTime           *int                   `gorm:"type:int" json:"total_time"`
	ServingSize         *int                   `gorm:"type:int" json:"serving_size"`
	Difficulty          string                 `gorm:"type:varchar(50)" json:"difficulty"`
	MealType            string                 `gorm:"type:varchar(50)" json:"meal_type"`
	Cuisine             string                 `gorm:"type:varchar(100)" json:"cuisine"`
	DietaryRestrictions RecipeDietaryJSON      `gorm:"type:jsonb" json:"dietary_restrictions"`
	NutritionInfo       RecipeNutritionJSON    `gorm:"type:jsonb" json:"nutrition_info"`
	Tips                RecipeTipsJSON         `gorm:"type:jsonb" json:"tips"`
	GeneratedAt         time.Time              `gorm:"type:timestamp with time zone;not null" json:"generated_at"`
	CreatedAt           time.Time              `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt           time.Time              `gorm:"type:timestamp with time zone;not null;default:now()" json:"updated_at"`
	DeletedAt           gorm.DeletedAt         `gorm:"index" json:"-"`
}

// RecipeIngredient represents an ingredient in a recipe
type RecipeIngredient struct {
	Name        string   `json:"name"`
	Amount      *float64 `json:"amount"`
	Unit        string   `json:"unit"`
	Available   bool     `json:"available"`
	Alternative *string  `json:"alternative,omitempty"`
}

// RecipeInstruction represents an instruction step in a recipe
type RecipeInstruction struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
	Time        *int   `json:"time,omitempty"`
}

// RecipeNutrition represents nutritional information
type RecipeNutrition struct {
	Calories      *int `json:"calories,omitempty"`
	Protein       *int `json:"protein,omitempty"`
	Carbohydrates *int `json:"carbohydrates,omitempty"`
	Fat           *int `json:"fat,omitempty"`
}

// Custom JSON types for GORM
type RecipeIngredientsJSON []RecipeIngredient
type RecipeInstructionsJSON []RecipeInstruction
type RecipeDietaryJSON []string
type RecipeTipsJSON []string
type RecipeNutritionJSON RecipeNutrition

// Scan implements the sql.Scanner interface for RecipeIngredientsJSON
func (r *RecipeIngredientsJSON) Scan(value interface{}) error {
	if value == nil {
		*r = []RecipeIngredient{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface for RecipeIngredientsJSON
func (r RecipeIngredientsJSON) Value() (driver.Value, error) {
	if len(r) == 0 {
		return "[]", nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for RecipeInstructionsJSON
func (r *RecipeInstructionsJSON) Scan(value interface{}) error {
	if value == nil {
		*r = []RecipeInstruction{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface for RecipeInstructionsJSON
func (r RecipeInstructionsJSON) Value() (driver.Value, error) {
	if len(r) == 0 {
		return "[]", nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for RecipeDietaryJSON
func (r *RecipeDietaryJSON) Scan(value interface{}) error {
	if value == nil {
		*r = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface for RecipeDietaryJSON
func (r RecipeDietaryJSON) Value() (driver.Value, error) {
	if len(r) == 0 {
		return "[]", nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for RecipeTipsJSON
func (r *RecipeTipsJSON) Scan(value interface{}) error {
	if value == nil {
		*r = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface for RecipeTipsJSON
func (r RecipeTipsJSON) Value() (driver.Value, error) {
	if len(r) == 0 {
		return "[]", nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for RecipeNutritionJSON
func (r *RecipeNutritionJSON) Scan(value interface{}) error {
	if value == nil {
		*r = RecipeNutritionJSON{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface for RecipeNutritionJSON
func (r RecipeNutritionJSON) Value() (driver.Value, error) {
	return json.Marshal(r)
}
