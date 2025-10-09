package domain

import (
	"context"

	"github.com/google/uuid"
	llmDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	recipeDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/dto"
	recipeModel "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/model"
)

type RecipeService interface {
	GenerateRecipe(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID) (*llmDTO.RecipeResponseDTO, error)
	GenerateMultipleRecipes(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID, count int) ([]*llmDTO.RecipeResponseDTO, error)
	SaveRecipe(ctx context.Context, recipe *recipeDTO.SaveRecipeDTO, userID uuid.UUID) error
	SaveMultipleRecipes(ctx context.Context, recipes []*recipeDTO.SaveRecipeDTO, userID uuid.UUID) error
	GetRecipeByID(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID) (*recipeDTO.RecipeDetailDTO, error)
	GetUserRecipes(ctx context.Context, userID uuid.UUID) ([]*recipeDTO.RecipeDetailDTO, error)
	GetAvailableIngredients(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]recipeDTO.AvailableIngredientDTO, error)
	SearchRecipesByIngredients(ctx context.Context, ingredients []string, filters map[string]string) ([]llmDTO.RecipeResponseDTO, error)
}

type RecipeRepository interface {
	Create(ctx context.Context, recipe *recipeModel.Recipe) error
	CreateMany(ctx context.Context, recipes []*recipeModel.Recipe) error
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*recipeModel.Recipe, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*recipeModel.Recipe, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
