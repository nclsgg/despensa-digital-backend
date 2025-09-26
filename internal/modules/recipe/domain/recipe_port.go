package domain

import (
	"context"

	"github.com/google/uuid"
	llmDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	recipeDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/dto"
)

type RecipeService interface {
	GenerateRecipe(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID) (*llmDTO.RecipeResponseDTO, error)
	GetAvailableIngredients(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]recipeDTO.AvailableIngredientDTO, error)
	SearchRecipesByIngredients(ctx context.Context, ingredients []string, filters map[string]string) ([]llmDTO.RecipeResponseDTO, error)
}
