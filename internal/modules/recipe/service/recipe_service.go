package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	itemDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	llmDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	llmSvc "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"
	pantryDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	pantrySvc "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/service"
	recipeDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/domain"
	recipeDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/dto"
	recipeModel "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/model"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type recipeService struct {
	llmService       *llmSvc.LLMServiceImpl
	itemRepository   itemDomain.ItemRepository
	pantryService    pantryDomain.PantryService
	recipeRepository recipeDomain.RecipeRepository
	promptBuilder    *llmSvc.PromptBuilderImpl
}

func NewRecipeService(
	llmService *llmSvc.LLMServiceImpl,
	itemRepository itemDomain.ItemRepository,
	pantryService pantryDomain.PantryService,
	recipeRepository recipeDomain.RecipeRepository,
) recipeDomain.RecipeService {
	return &recipeService{
		llmService:       llmService,
		itemRepository:   itemRepository,
		pantryService:    pantryService,
		recipeRepository: recipeRepository,
		promptBuilder:    llmSvc.NewPromptBuilder(),
	}
}

func (rs *recipeService) GenerateRecipe(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID) (*llmDTO.RecipeResponseDTO, error) {
	logger := appLogger.FromContext(ctx)

	if request == nil {
		logger.Warn("Recipe generation request is nil",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
		)
		return nil, fmt.Errorf("%w: request payload is required", recipeDomain.ErrInvalidRequest)
	}

	request.SetDefaults()

	pantryID, err := rs.validateRecipeRequest(request)
	if err != nil {
		logger.Warn("Invalid recipe request",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	availableIngredients, err := rs.GetAvailableIngredients(ctx, pantryID, userID)
	if err != nil {
		logger.Error("Failed to get available ingredients",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if len(availableIngredients) == 0 {
		logger.Warn("No ingredients available in pantry",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
		)
		return nil, recipeDomain.ErrNoIngredients
	}

	variables := rs.buildPromptVariables(request, availableIngredients)
	templates := llmSvc.GetRecipePromptTemplates()

	systemPrompt, err := rs.promptBuilder.BuildSystemPrompt(templates.SystemPrompt, variables)
	if err != nil {
		logger.Error("Failed to build system prompt",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRequest, err)
	}

	userPrompt, err := rs.promptBuilder.BuildUserPrompt(templates.UserPrompt, variables)
	if err != nil {
		logger.Error("Failed to build user prompt",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRequest, err)
	}

	options := map[string]interface{}{
		"max_tokens":  2000,
		"temperature": 0.7,
		"top_p":       0.9,
	}

	llmRequest := rs.llmService.CreateChatRequest(systemPrompt, userPrompt, options)

	var llmResponse *llmDTO.LLMResponseDTO
	if request.Provider != "" {
		llmResponse, err = rs.llmService.ProcessRequestWithProvider(ctx, llmRequest, request.Provider)
	} else {
		llmResponse, err = rs.llmService.ProcessRequest(ctx, llmRequest)
	}

	if err != nil {
		logger.Error("LLM request failed",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("provider", request.Provider),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%w: %v", recipeDomain.ErrLLMRequest, err)
	}

	recipe, err := rs.parseRecipeResponse(llmResponse.Response)
	if err != nil {
		logger.Error("Failed to parse LLM response",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%w: %v", recipeDomain.ErrInvalidLLMResponse, err)
	}

	recipe.ID = uuid.New().String()
	recipe.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	rs.markIngredientAvailability(recipe, availableIngredients)

	logger.Info("Recipe generated successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GenerateRecipe"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("recipe_id", recipe.ID),
		zap.String("title", recipe.Title),
		zap.Int("ingredient_count", len(recipe.Ingredients)),
	)

	return recipe, nil
}

func (rs *recipeService) GetAvailableIngredients(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]recipeDTO.AvailableIngredientDTO, error) {
	logger := appLogger.FromContext(ctx)

	if _, err := rs.pantryService.GetPantry(ctx, pantryID, userID); err != nil {
		switch {
		case errors.Is(err, pantrySvc.ErrUnauthorized):
			logger.Warn("Unauthorized access to pantry",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
			)
			return nil, recipeDomain.ErrUnauthorized
		case errors.Is(err, pantrySvc.ErrPantryNotFound):
			logger.Warn("Pantry not found",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
			)
			return nil, recipeDomain.ErrPantryNotFound
		default:
			logger.Error("Failed to get pantry",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
				zap.Error(err),
			)
			return nil, err
		}
	}

	items, err := rs.itemRepository.ListByPantryID(ctx, pantryID)
	if err != nil {
		logger.Error("Failed to list items from pantry",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	ingredients := make([]recipeDTO.AvailableIngredientDTO, 0, len(items))
	for _, item := range items {
		if item.Quantity > 0 {
			ingredients = append(ingredients, recipeDTO.AvailableIngredientDTO{
				Name:     strings.TrimSpace(item.Name),
				Quantity: item.Quantity,
				Unit:     strings.TrimSpace(item.Unit),
			})
		}
	}

	logger.Info("Available ingredients retrieved",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.Int(appLogger.FieldCount, len(ingredients)),
	)

	return ingredients, nil
}

func (rs *recipeService) SearchRecipesByIngredients(ctx context.Context, ingredients []string, filters map[string]string) ([]llmDTO.RecipeResponseDTO, error) {
	logger := appLogger.FromContext(ctx)

	logger.Warn("Search by ingredients not implemented",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "SearchRecipesByIngredients"),
	)

	return nil, fmt.Errorf("%w: search by ingredients not implemented", recipeDomain.ErrInvalidRequest)
}

func (rs *recipeService) validateRecipeRequest(request *llmDTO.RecipeRequestDTO) (uuid.UUID, error) {
	if request.PantryID == "" {
		return uuid.Nil, fmt.Errorf("%w: pantry_id is required", recipeDomain.ErrInvalidRequest)
	}

	pantryID, err := uuid.Parse(request.PantryID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: pantry_id must be a valid UUID", recipeDomain.ErrInvalidRequest)
	}

	if request.CookingTime != 0 && (request.CookingTime < 5 || request.CookingTime > 480) {
		return uuid.Nil, fmt.Errorf("%w: cooking_time must be between 5 and 480 minutes", recipeDomain.ErrInvalidRequest)
	}

	if request.ServingSize < 0 || request.ServingSize > 20 {
		return uuid.Nil, fmt.Errorf("%w: serving_size must be between 1 and 20", recipeDomain.ErrInvalidRequest)
	}

	validMealTypes := map[string]bool{
		"breakfast": true,
		"lunch":     true,
		"dinner":    true,
		"snack":     true,
		"dessert":   true,
		"":          true,
	}

	if !validMealTypes[strings.ToLower(strings.TrimSpace(request.MealType))] {
		return uuid.Nil, fmt.Errorf("%w: meal_type is invalid", recipeDomain.ErrInvalidRequest)
	}

	validDifficulties := map[string]bool{
		"easy":   true,
		"medium": true,
		"hard":   true,
		"":       true,
	}

	if !validDifficulties[strings.ToLower(strings.TrimSpace(request.Difficulty))] {
		return uuid.Nil, fmt.Errorf("%w: difficulty is invalid", recipeDomain.ErrInvalidRequest)
	}

	return pantryID, nil
}

// EnrichRecipeWithNutrition adiciona informações nutricionais (placeholder)
func (rs *recipeService) EnrichRecipeWithNutrition(ctx context.Context, recipe *llmDTO.RecipeResponseDTO) error {
	// TODO: Implementar cálculo nutricional real
	return nil
}

// buildPromptVariables constrói as variáveis para o prompt
func (rs *recipeService) buildPromptVariables(request *llmDTO.RecipeRequestDTO, ingredients []recipeDTO.AvailableIngredientDTO) map[string]string {
	var formattedIngredients []string
	for _, ingredient := range ingredients {
		formatted := fmt.Sprintf("%s (%.1f %s)", ingredient.Name, ingredient.Quantity, ingredient.Unit)
		formattedIngredients = append(formattedIngredients, formatted)
	}

	variables := map[string]string{
		"available_ingredients": strings.Join(formattedIngredients, ", "),
	}

	if request.CookingTime > 0 {
		variables["cooking_time"] = fmt.Sprintf("%d", request.CookingTime)
	} else {
		variables["cooking_time"] = "não especificado"
	}

	if request.MealType != "" {
		variables["meal_type"] = strings.ToLower(strings.TrimSpace(request.MealType))
	} else {
		variables["meal_type"] = "qualquer"
	}

	if request.Difficulty != "" {
		variables["difficulty"] = strings.ToLower(strings.TrimSpace(request.Difficulty))
	} else {
		variables["difficulty"] = "qualquer"
	}

	if request.ServingSize > 0 {
		variables["serving_size"] = fmt.Sprintf("%d", request.ServingSize)
	} else {
		variables["serving_size"] = "4"
	}

	if request.Cuisine != "" {
		variables["cuisine"] = strings.TrimSpace(request.Cuisine)
	}

	if len(request.DietaryRestrictions) > 0 {
		variables["dietary_restrictions"] = strings.Join(cleanStringSlice(request.DietaryRestrictions), ", ")
	}

	if request.Purpose != "" {
		variables["purpose"] = strings.TrimSpace(request.Purpose)
	}

	if request.AdditionalNotes != "" {
		variables["additional_notes"] = strings.TrimSpace(request.AdditionalNotes)
	}

	return variables
}

// parseRecipeResponse parseia a resposta JSON do LLM
func (rs *recipeService) parseRecipeResponse(response string) (*llmDTO.RecipeResponseDTO, error) {
	// Remove possíveis caracteres extras antes e depois do JSON
	response = strings.TrimSpace(response)

	// Procura pelo início e fim do JSON
	startIndex := strings.Index(response, "{")
	if startIndex == -1 {
		return nil, fmt.Errorf("JSON não encontrado na resposta")
	}

	endIndex := strings.LastIndex(response, "}")
	if endIndex == -1 {
		return nil, fmt.Errorf("JSON malformado na resposta")
	}

	jsonResponse := response[startIndex : endIndex+1]

	// Primeiro, tenta parsear diretamente
	var recipe llmDTO.RecipeResponseDTO
	if err := json.Unmarshal([]byte(jsonResponse), &recipe); err != nil {
		// Se falhar, tenta corrigir problemas comuns
		correctedJSON := rs.fixCommonJSONIssues(jsonResponse)
		if err2 := json.Unmarshal([]byte(correctedJSON), &recipe); err2 != nil {
			return nil, fmt.Errorf("erro ao parsear JSON da receita: %w, resposta: %s", err, jsonResponse)
		}
	}

	return &recipe, nil
}

// fixCommonJSONIssues corrige problemas comuns no JSON retornado pelo LLM
func (rs *recipeService) fixCommonJSONIssues(jsonStr string) string {
	// Substitui frações matemáticas por decimais
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 1/2`, `"amount": 0.5`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 1/3`, `"amount": 0.33`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 2/3`, `"amount": 0.67`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 1/4`, `"amount": 0.25`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 3/4`, `"amount": 0.75`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 1/8`, `"amount": 0.125`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 3/8`, `"amount": 0.375`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 5/8`, `"amount": 0.625`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": 7/8`, `"amount": 0.875`)

	// Substitui strings "a gosto" por null
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": "a gosto"`, `"amount": null`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": "à gosto"`, `"amount": null`)
	jsonStr = strings.ReplaceAll(jsonStr, `"amount": "ao gosto"`, `"amount": null`)

	// Corrige valores numéricos em campos de temperature que devem ser string
	jsonStr = strings.ReplaceAll(jsonStr, `"temperature": null`, `"temperature": null`)

	// Remove trailing commas antes de fechar objetos/arrays (comum em JSON gerado por LLMs)
	jsonStr = strings.ReplaceAll(jsonStr, ",\n  }", "\n  }")
	jsonStr = strings.ReplaceAll(jsonStr, ",\n]", "\n]")
	jsonStr = strings.ReplaceAll(jsonStr, ", }", " }")
	jsonStr = strings.ReplaceAll(jsonStr, ", ]", " ]")

	return jsonStr
}

// markIngredientAvailability marca quais ingredientes estão disponíveis
func (rs *recipeService) markIngredientAvailability(recipe *llmDTO.RecipeResponseDTO, availableIngredients []recipeDTO.AvailableIngredientDTO) {
	// Cria um mapa para busca rápida
	availableMap := make(map[string]bool)
	for _, ingredient := range availableIngredients {
		availableMap[strings.ToLower(strings.TrimSpace(ingredient.Name))] = true
	}

	// Marca disponibilidade para cada ingrediente da receita
	for i := range recipe.Ingredients {
		ingredientName := strings.ToLower(strings.TrimSpace(recipe.Ingredients[i].Name))
		recipe.Ingredients[i].Available = availableMap[ingredientName]
	}
}

func cleanStringSlice(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return cleaned
}

// GenerateMultipleRecipes generates multiple recipes based on pantry ingredients
func (rs *recipeService) GenerateMultipleRecipes(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID, count int) ([]*llmDTO.RecipeResponseDTO, error) {
	logger := appLogger.FromContext(ctx)

	if count < 1 || count > 10 {
		logger.Warn("Invalid recipe count",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateMultipleRecipes"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Int(appLogger.FieldCount, count),
		)
		return nil, fmt.Errorf("%w: count must be between 1 and 10", recipeDomain.ErrInvalidRequest)
	}

	recipes := make([]*llmDTO.RecipeResponseDTO, 0, count)
	for i := 0; i < count; i++ {
		recipe, err := rs.GenerateRecipe(ctx, request, userID)
		if err != nil {
			logger.Error("Failed to generate recipe in batch",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "GenerateMultipleRecipes"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Int("recipe_index", i),
				zap.Error(err),
			)
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	logger.Info("Multiple recipes generated successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GenerateMultipleRecipes"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(recipes)),
	)

	return recipes, nil
}

// SaveRecipe saves a single recipe to the database
func (rs *recipeService) SaveRecipe(ctx context.Context, recipeDTO *recipeDTO.SaveRecipeDTO, userID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	if err := rs.validateSaveRecipeDTO(recipeDTO); err != nil {
		logger.Warn("Invalid recipe data for save",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return err
	}

	recipe, err := rs.convertSaveRecipeDTOToModel(recipeDTO, userID)
	if err != nil {
		logger.Error("Failed to convert recipe DTO to model",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRecipeData, err)
	}

	if err := rs.recipeRepository.Create(ctx, recipe); err != nil {
		logger.Error("Failed to create recipe in database",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("recipe_id", recipe.ID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("Recipe saved successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "SaveRecipe"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("recipe_id", recipe.ID.String()),
		zap.String("title", recipe.Title),
	)

	return nil
}

// SaveMultipleRecipes saves multiple recipes to the database atomically
func (rs *recipeService) SaveMultipleRecipes(ctx context.Context, recipeDTOs []*recipeDTO.SaveRecipeDTO, userID uuid.UUID) error {
	logger := appLogger.FromContext(ctx)

	if len(recipeDTOs) == 0 {
		logger.Warn("Empty recipe list for batch save",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveMultipleRecipes"),
			zap.String(appLogger.FieldUserID, userID.String()),
		)
		return fmt.Errorf("%w: at least one recipe is required", recipeDomain.ErrInvalidRequest)
	}

	recipes := make([]*recipeModel.Recipe, 0, len(recipeDTOs))
	for _, dto := range recipeDTOs {
		if err := rs.validateSaveRecipeDTO(dto); err != nil {
			logger.Warn("Invalid recipe data in batch",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "SaveMultipleRecipes"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			return err
		}

		recipe, err := rs.convertSaveRecipeDTOToModel(dto, userID)
		if err != nil {
			logger.Error("Failed to convert recipe DTO in batch",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "SaveMultipleRecipes"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			return fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRecipeData, err)
		}
		recipes = append(recipes, recipe)
	}

	if err := rs.recipeRepository.CreateMany(ctx, recipes); err != nil {
		logger.Error("Failed to create multiple recipes in database",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveMultipleRecipes"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Int(appLogger.FieldCount, len(recipes)),
			zap.Error(err),
		)
		return err
	}

	logger.Info("Multiple recipes saved successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "SaveMultipleRecipes"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(recipes)),
	)

	return nil
}

// GetRecipeByID retrieves a single recipe by ID
func (rs *recipeService) GetRecipeByID(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID) (*recipeDTO.RecipeDetailDTO, error) {
	logger := appLogger.FromContext(ctx)

	recipe, err := rs.recipeRepository.FindByID(ctx, recipeID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Recipe not found",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "GetRecipeByID"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("recipe_id", recipeID.String()),
			)
			return nil, recipeDomain.ErrRecipeNotFound
		}

		logger.Error("Failed to get recipe from database",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetRecipeByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("recipe_id", recipeID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Recipe retrieved successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GetRecipeByID"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("recipe_id", recipeID.String()),
	)

	return rs.convertModelToRecipeDetailDTO(recipe), nil
}

// GetUserRecipes retrieves all recipes for a user
func (rs *recipeService) GetUserRecipes(ctx context.Context, userID uuid.UUID) ([]*recipeDTO.RecipeDetailDTO, error) {
	logger := appLogger.FromContext(ctx)

	recipes, err := rs.recipeRepository.FindByUserID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user recipes from database",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetUserRecipes"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	recipeDTOs := make([]*recipeDTO.RecipeDetailDTO, 0, len(recipes))
	for _, recipe := range recipes {
		recipeDTOs = append(recipeDTOs, rs.convertModelToRecipeDetailDTO(recipe))
	}

	logger.Info("User recipes retrieved successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GetUserRecipes"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(recipeDTOs)),
	)

	return recipeDTOs, nil
}

// Helper methods

func (rs *recipeService) validateSaveRecipeDTO(dto *recipeDTO.SaveRecipeDTO) error {
	if dto == nil {
		return fmt.Errorf("%w: recipe data is required", recipeDomain.ErrInvalidRequest)
	}

	if dto.ID == "" {
		return fmt.Errorf("%w: recipe ID is required", recipeDomain.ErrInvalidRequest)
	}

	if _, err := uuid.Parse(dto.ID); err != nil {
		return fmt.Errorf("%w: invalid recipe ID format", recipeDomain.ErrInvalidRequest)
	}

	if dto.Title == "" {
		return fmt.Errorf("%w: title is required", recipeDomain.ErrInvalidRequest)
	}

	if len(dto.Ingredients) == 0 {
		return fmt.Errorf("%w: at least one ingredient is required", recipeDomain.ErrInvalidRequest)
	}

	if len(dto.Instructions) == 0 {
		return fmt.Errorf("%w: at least one instruction is required", recipeDomain.ErrInvalidRequest)
	}

	return nil
}

func (rs *recipeService) convertSaveRecipeDTOToModel(dto *recipeDTO.SaveRecipeDTO, userID uuid.UUID) (*recipeModel.Recipe, error) {
	recipeID, err := uuid.Parse(dto.ID)
	if err != nil {
		return nil, err
	}

	// Convert ingredients
	ingredients := make([]recipeModel.RecipeIngredient, 0, len(dto.Ingredients))
	for _, ing := range dto.Ingredients {
		ingredients = append(ingredients, recipeModel.RecipeIngredient{
			Name:        ing.Name,
			Amount:      ing.Amount,
			Unit:        ing.Unit,
			Available:   ing.Available,
			Alternative: ing.Alternative,
		})
	}

	// Convert instructions
	instructions := make([]recipeModel.RecipeInstruction, 0, len(dto.Instructions))
	for _, inst := range dto.Instructions {
		instructions = append(instructions, recipeModel.RecipeInstruction{
			Step:        inst.Step,
			Description: inst.Description,
			Time:        inst.Time,
		})
	}

	// Convert dietary restrictions
	dietaryRestrictions := dto.DietaryRestrictions
	if dietaryRestrictions == nil {
		dietaryRestrictions = []string{}
	}

	// Convert tips
	tips := dto.Tips
	if tips == nil {
		tips = []string{}
	}

	// Convert nutrition info
	nutritionInfo := recipeModel.RecipeNutrition{
		Calories:      dto.NutritionInfo.Calories,
		Protein:       dto.NutritionInfo.Protein,
		Carbohydrates: dto.NutritionInfo.Carbohydrates,
		Fat:           dto.NutritionInfo.Fat,
	}

	// Parse generated_at timestamp
	var generatedAt time.Time
	if dto.GeneratedAt != "" {
		generatedAt, err = time.Parse(time.RFC3339, dto.GeneratedAt)
		if err != nil {
			// Try alternative formats
			generatedAt, err = time.Parse("2006-01-02T15:04:05Z07:00", dto.GeneratedAt)
			if err != nil {
				generatedAt = time.Now().UTC()
			}
		}
	} else {
		generatedAt = time.Now().UTC()
	}

	recipe := &recipeModel.Recipe{
		ID:                  recipeID,
		UserID:              userID,
		Title:               dto.Title,
		Description:         dto.Description,
		Ingredients:         recipeModel.RecipeIngredientsJSON(ingredients),
		Instructions:        recipeModel.RecipeInstructionsJSON(instructions),
		CookingTime:         dto.CookingTime,
		PreparationTime:     dto.PreparationTime,
		TotalTime:           dto.TotalTime,
		ServingSize:         dto.ServingSize,
		Difficulty:          dto.Difficulty,
		MealType:            dto.MealType,
		Cuisine:             dto.Cuisine,
		DietaryRestrictions: recipeModel.RecipeDietaryJSON(dietaryRestrictions),
		NutritionInfo:       recipeModel.RecipeNutritionJSON(nutritionInfo),
		Tips:                recipeModel.RecipeTipsJSON(tips),
		GeneratedAt:         generatedAt,
	}

	return recipe, nil
}

func (rs *recipeService) convertModelToRecipeDetailDTO(recipe *recipeModel.Recipe) *recipeDTO.RecipeDetailDTO {
	// Convert ingredients
	ingredients := make([]recipeDTO.RecipeIngredientDetailDTO, 0, len(recipe.Ingredients))
	for _, ing := range recipe.Ingredients {
		ingredients = append(ingredients, recipeDTO.RecipeIngredientDetailDTO{
			Name:        ing.Name,
			Amount:      ing.Amount,
			Unit:        ing.Unit,
			Available:   ing.Available,
			Alternative: ing.Alternative,
		})
	}

	// Convert instructions
	instructions := make([]recipeDTO.RecipeInstructionDetailDTO, 0, len(recipe.Instructions))
	for _, inst := range recipe.Instructions {
		instructions = append(instructions, recipeDTO.RecipeInstructionDetailDTO{
			Step:        inst.Step,
			Description: inst.Description,
			Time:        inst.Time,
		})
	}

	// Convert dietary restrictions
	dietaryRestrictions := []string(recipe.DietaryRestrictions)
	if dietaryRestrictions == nil {
		dietaryRestrictions = []string{}
	}

	// Convert tips
	tips := []string(recipe.Tips)
	if tips == nil {
		tips = []string{}
	}

	// Convert nutrition info
	nutritionInfo := recipeDTO.RecipeNutritionDetailDTO{
		Calories:      recipe.NutritionInfo.Calories,
		Protein:       recipe.NutritionInfo.Protein,
		Carbohydrates: recipe.NutritionInfo.Carbohydrates,
		Fat:           recipe.NutritionInfo.Fat,
	}

	return &recipeDTO.RecipeDetailDTO{
		ID:                  recipe.ID.String(),
		Title:               recipe.Title,
		Description:         recipe.Description,
		Ingredients:         ingredients,
		Instructions:        instructions,
		CookingTime:         recipe.CookingTime,
		PreparationTime:     recipe.PreparationTime,
		TotalTime:           recipe.TotalTime,
		ServingSize:         recipe.ServingSize,
		Difficulty:          recipe.Difficulty,
		MealType:            recipe.MealType,
		Cuisine:             recipe.Cuisine,
		DietaryRestrictions: dietaryRestrictions,
		NutritionInfo:       nutritionInfo,
		Tips:                tips,
		GeneratedAt:         recipe.GeneratedAt,
		CreatedAt:           recipe.CreatedAt,
	}
}
