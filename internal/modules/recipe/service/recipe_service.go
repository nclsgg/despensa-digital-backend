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
)

type recipeService struct {
	llmService     *llmSvc.LLMServiceImpl
	itemRepository itemDomain.ItemRepository
	pantryService  pantryDomain.PantryService
	promptBuilder  *llmSvc.PromptBuilderImpl
}

func NewRecipeService(
	llmService *llmSvc.LLMServiceImpl,
	itemRepository itemDomain.ItemRepository,
	pantryService pantryDomain.PantryService,
) recipeDomain.RecipeService {
	return &recipeService{
		llmService:     llmService,
		itemRepository: itemRepository,
		pantryService:  pantryService,
		promptBuilder:  llmSvc.NewPromptBuilder(),
	}
}

func (rs *recipeService) GenerateRecipe(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID) (*llmDTO.RecipeResponseDTO, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request payload is required", recipeDomain.ErrInvalidRequest)
	}

	request.SetDefaults()

	pantryID, err := rs.validateRecipeRequest(request)
	if err != nil {
		return nil, err
	}

	availableIngredients, err := rs.GetAvailableIngredients(ctx, pantryID, userID)
	if err != nil {
		return nil, err
	}

	if len(availableIngredients) == 0 {
		return nil, recipeDomain.ErrNoIngredients
	}

	variables := rs.buildPromptVariables(request, availableIngredients)
	templates := llmSvc.GetRecipePromptTemplates()

	systemPrompt, err := rs.promptBuilder.BuildSystemPrompt(templates.SystemPrompt, variables)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRequest, err)
	}

	userPrompt, err := rs.promptBuilder.BuildUserPrompt(templates.UserPrompt, variables)
	if err != nil {
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
		return nil, fmt.Errorf("%w: %v", recipeDomain.ErrLLMRequest, err)
	}

	recipe, err := rs.parseRecipeResponse(llmResponse.Response)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", recipeDomain.ErrInvalidLLMResponse, err)
	}

	recipe.ID = uuid.New().String()
	recipe.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	rs.markIngredientAvailability(recipe, availableIngredients)

	return recipe, nil
}

func (rs *recipeService) GetAvailableIngredients(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]recipeDTO.AvailableIngredientDTO, error) {
	if _, err := rs.pantryService.GetPantry(ctx, pantryID, userID); err != nil {
		switch {
		case errors.Is(err, pantrySvc.ErrUnauthorized):
			return nil, recipeDomain.ErrUnauthorized
		case errors.Is(err, pantrySvc.ErrPantryNotFound):
			return nil, recipeDomain.ErrPantryNotFound
		default:
			return nil, err
		}
	}

	items, err := rs.itemRepository.ListByPantryID(ctx, pantryID)
	if err != nil {
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

	return ingredients, nil
}

func (rs *recipeService) SearchRecipesByIngredients(ctx context.Context, ingredients []string, filters map[string]string) ([]llmDTO.RecipeResponseDTO, error) {
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
	// Formata ingredientes com quantidade e unidade
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
