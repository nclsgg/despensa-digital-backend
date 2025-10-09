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
) (result0 recipeDomain.RecipeService) {
	__logParams := map[string]any{"llmService": llmService, "itemRepository": itemRepository, "pantryService": pantryService, "recipeRepository": recipeRepository}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewRecipeService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewRecipeService"), zap.Any("params", __logParams))
	result0 = &recipeService{
		llmService:       llmService,
		itemRepository:   itemRepository,
		pantryService:    pantryService,
		recipeRepository: recipeRepository,
		promptBuilder:    llmSvc.NewPromptBuilder(),
	}
	return
}

func (rs *recipeService) GenerateRecipe(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID) (result0 *llmDTO.RecipeResponseDTO, result1 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "request": request, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.GenerateRecipe"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.GenerateRecipe"), zap.Any("params", __logParams))
	if request == nil {
		result0 = nil
		result1 = fmt.Errorf("%w: request payload is required", recipeDomain.ErrInvalidRequest)
		return
	}

	request.SetDefaults()

	pantryID, err := rs.validateRecipeRequest(request)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GenerateRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	availableIngredients, err := rs.GetAvailableIngredients(ctx, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GenerateRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	if len(availableIngredients) == 0 {
		result0 = nil
		result1 = recipeDomain.ErrNoIngredients
		return
	}

	variables := rs.buildPromptVariables(request, availableIngredients)
	templates := llmSvc.GetRecipePromptTemplates()

	systemPrompt, err := rs.promptBuilder.BuildSystemPrompt(templates.SystemPrompt, variables)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GenerateRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRequest, err)
		return
	}

	userPrompt, err := rs.promptBuilder.BuildUserPrompt(templates.UserPrompt, variables)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GenerateRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRequest, err)
		return
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
		zap.L().Error("function.error", zap.String("func", "*recipeService.GenerateRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("%w: %v", recipeDomain.ErrLLMRequest, err)
		return
	}

	recipe, err := rs.parseRecipeResponse(llmResponse.Response)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GenerateRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("%w: %v", recipeDomain.ErrInvalidLLMResponse, err)
		return
	}

	recipe.ID = uuid.New().String()
	recipe.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	rs.markIngredientAvailability(recipe, availableIngredients)
	result0 = recipe
	result1 = nil
	return
}

func (rs *recipeService) GetAvailableIngredients(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) (result0 []recipeDTO.AvailableIngredientDTO, result1 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "pantryID": pantryID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.GetAvailableIngredients"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.GetAvailableIngredients"), zap.Any("params", __logParams))
	if _, err := rs.pantryService.GetPantry(ctx, pantryID, userID); err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GetAvailableIngredients"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, pantrySvc.ErrUnauthorized):
			result0 = nil
			result1 = recipeDomain.ErrUnauthorized
			return
		case errors.Is(err, pantrySvc.ErrPantryNotFound):
			result0 = nil
			result1 = recipeDomain.ErrPantryNotFound
			return
		default:
			result0 = nil
			result1 = err
			return
		}
	}

	items, err := rs.itemRepository.ListByPantryID(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GetAvailableIngredients"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
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
	result0 = ingredients
	result1 = nil
	return
}

func (rs *recipeService) SearchRecipesByIngredients(ctx context.Context, ingredients []string, filters map[string]string) (result0 []llmDTO.RecipeResponseDTO, result1 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "ingredients": ingredients, "filters": filters}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.SearchRecipesByIngredients"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.SearchRecipesByIngredients"), zap.Any("params", __logParams))
	result0 = nil
	result1 = fmt.Errorf("%w: search by ingredients not implemented", recipeDomain.ErrInvalidRequest)
	return
}

func (rs *recipeService) validateRecipeRequest(request *llmDTO.RecipeRequestDTO) (result0 uuid.UUID, result1 error) {
	__logParams := map[string]any{"rs": rs, "request": request}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.validateRecipeRequest"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.validateRecipeRequest"), zap.Any("params", __logParams))
	if request.PantryID == "" {
		result0 = uuid.Nil
		result1 = fmt.Errorf("%w: pantry_id is required", recipeDomain.ErrInvalidRequest)
		return
	}

	pantryID, err := uuid.Parse(request.PantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.validateRecipeRequest"), zap.Error(err), zap.Any("params", __logParams))
		result0 = uuid.Nil
		result1 = fmt.Errorf("%w: pantry_id must be a valid UUID", recipeDomain.ErrInvalidRequest)
		return
	}

	if request.CookingTime != 0 && (request.CookingTime < 5 || request.CookingTime > 480) {
		result0 = uuid.Nil
		result1 = fmt.Errorf("%w: cooking_time must be between 5 and 480 minutes", recipeDomain.ErrInvalidRequest)
		return
	}

	if request.ServingSize < 0 || request.ServingSize > 20 {
		result0 = uuid.Nil
		result1 = fmt.Errorf("%w: serving_size must be between 1 and 20", recipeDomain.ErrInvalidRequest)
		return
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
		result0 = uuid.Nil
		result1 = fmt.Errorf("%w: meal_type is invalid", recipeDomain.ErrInvalidRequest)
		return
	}

	validDifficulties := map[string]bool{
		"easy":   true,
		"medium": true,
		"hard":   true,
		"":       true,
	}

	if !validDifficulties[strings.ToLower(strings.TrimSpace(request.Difficulty))] {
		result0 = uuid.Nil
		result1 = fmt.Errorf("%w: difficulty is invalid", recipeDomain.ErrInvalidRequest)
		return
	}
	result0 = pantryID
	result1 = nil
	return
}

// EnrichRecipeWithNutrition adiciona informações nutricionais (placeholder)
func (rs *recipeService) EnrichRecipeWithNutrition(ctx context.Context, recipe *llmDTO.RecipeResponseDTO) (result0 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "recipe": recipe}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.EnrichRecipeWithNutrition"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.EnrichRecipeWithNutrition"), zap.Any("params", __logParams))

	// TODO: Implementar cálculo nutricional real
	result0 = nil
	return
}

// buildPromptVariables constrói as variáveis para o prompt
func (rs *recipeService) buildPromptVariables(request *llmDTO.RecipeRequestDTO, ingredients []recipeDTO.AvailableIngredientDTO) (result0 map[string]string) {
	__logParams := map[string]any{"rs": rs, "request": request, "ingredients": ingredients}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.buildPromptVariables"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.buildPromptVariables"), zap.Any("params", __logParams))

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
	result0 = variables
	return
}

// parseRecipeResponse parseia a resposta JSON do LLM
func (rs *recipeService) parseRecipeResponse(response string) (result0 *llmDTO.RecipeResponseDTO, result1 error) {
	__logParams :=
		// Remove possíveis caracteres extras antes e depois do JSON
		map[string]any{"rs": rs, "response": response}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit",

			// Procura pelo início e fim do JSON
			zap.String("func", "*recipeService.parseRecipeResponse"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.parseRecipeResponse"), zap.Any("params", __logParams))

	response = strings.TrimSpace(response)

	startIndex := strings.Index(response, "{")
	if startIndex == -1 {
		result0 = nil
		result1 = fmt.Errorf("JSON não encontrado na resposta")
		return
	}

	endIndex := strings.LastIndex(response, "}")
	if endIndex == -1 {
		result0 = nil
		result1 = fmt.Errorf("JSON malformado na resposta")
		return
	}

	jsonResponse := response[startIndex : endIndex+1]

	// Primeiro, tenta parsear diretamente
	var recipe llmDTO.RecipeResponseDTO
	if err := json.Unmarshal([]byte(jsonResponse), &recipe); err != nil {
		zap.L(
		// Se falhar, tenta corrigir problemas comuns
		).Error("function.error", zap.String("func", "*recipeService.parseRecipeResponse"), zap.Error(err), zap.Any("params", __logParams))

		correctedJSON := rs.fixCommonJSONIssues(jsonResponse)
		if err2 := json.Unmarshal([]byte(correctedJSON), &recipe); err2 != nil {
			zap.L().Error("function.error", zap.String("func", "*recipeService.parseRecipeResponse"), zap.Error(err2), zap.Any("params", __logParams))
			result0 = nil
			result1 = fmt.Errorf("erro ao parsear JSON da receita: %w, resposta: %s", err, jsonResponse)
			return
		}
	}
	result0 = &recipe
	result1 = nil
	return
}

// fixCommonJSONIssues corrige problemas comuns no JSON retornado pelo LLM
func (rs *recipeService) fixCommonJSONIssues(jsonStr string) (result0 string) {
	__logParams :=
		// Substitui frações matemáticas por decimais
		map[string]any{"rs": rs, "jsonStr": jsonStr}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.fixCommonJSONIssues"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.fixCommonJSONIssues"), zap.Any("params", __logParams))

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

	result0 = jsonStr
	return
}

// markIngredientAvailability marca quais ingredientes estão disponíveis
func (rs *recipeService) markIngredientAvailability(recipe *llmDTO.RecipeResponseDTO, availableIngredients []recipeDTO.AvailableIngredientDTO) {
	__logParams :=
		// Cria um mapa para busca rápida
		map[string]any{"rs": rs, "recipe": recipe, "availableIngredients": availableIngredients}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.markIngredientAvailability"),

			// Marca disponibilidade para cada ingrediente da receita
			zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.markIngredientAvailability"), zap.Any("params", __logParams))

	availableMap := make(map[string]bool)
	for _, ingredient := range availableIngredients {
		availableMap[strings.ToLower(strings.TrimSpace(ingredient.Name))] = true
	}

	for i := range recipe.Ingredients {
		ingredientName := strings.ToLower(strings.TrimSpace(recipe.Ingredients[i].Name))
		recipe.Ingredients[i].Available = availableMap[ingredientName]
	}
}

func cleanStringSlice(values []string) (result0 []string) {
	__logParams := map[string]any{"values": values}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "cleanStringSlice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "cleanStringSlice"), zap.Any("params", __logParams))
	if len(values) == 0 {
		result0 = []string{}
		return
	}
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	result0 = cleaned
	return
}

// GenerateMultipleRecipes generates multiple recipes based on pantry ingredients
func (rs *recipeService) GenerateMultipleRecipes(ctx context.Context, request *llmDTO.RecipeRequestDTO, userID uuid.UUID, count int) (result0 []*llmDTO.RecipeResponseDTO, result1 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "request": request, "userID": userID, "count": count}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.GenerateMultipleRecipes"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.GenerateMultipleRecipes"), zap.Any("params", __logParams))

	if count < 1 || count > 10 {
		result0 = nil
		result1 = fmt.Errorf("%w: count must be between 1 and 10", recipeDomain.ErrInvalidRequest)
		return
	}

	recipes := make([]*llmDTO.RecipeResponseDTO, 0, count)
	for i := 0; i < count; i++ {
		recipe, err := rs.GenerateRecipe(ctx, request, userID)
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "*recipeService.GenerateMultipleRecipes"), zap.Error(err), zap.Any("params", __logParams))
			result0 = nil
			result1 = err
			return
		}
		recipes = append(recipes, recipe)
	}

	result0 = recipes
	result1 = nil
	return
}

// SaveRecipe saves a single recipe to the database
func (rs *recipeService) SaveRecipe(ctx context.Context, recipeDTO *recipeDTO.SaveRecipeDTO, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "recipeDTO": recipeDTO, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.SaveRecipe"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.SaveRecipe"), zap.Any("params", __logParams))

	if err := rs.validateSaveRecipeDTO(recipeDTO); err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.SaveRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}

	recipe, err := rs.convertSaveRecipeDTOToModel(recipeDTO, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.SaveRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRecipeData, err)
		return
	}

	if err := rs.recipeRepository.Create(ctx, recipe); err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.SaveRecipe"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}

	result0 = nil
	return
}

// SaveMultipleRecipes saves multiple recipes to the database atomically
func (rs *recipeService) SaveMultipleRecipes(ctx context.Context, recipeDTOs []*recipeDTO.SaveRecipeDTO, userID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "recipeDTOs": recipeDTOs, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.SaveMultipleRecipes"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.SaveMultipleRecipes"), zap.Any("params", __logParams))

	if len(recipeDTOs) == 0 {
		result0 = fmt.Errorf("%w: at least one recipe is required", recipeDomain.ErrInvalidRequest)
		return
	}

	recipes := make([]*recipeModel.Recipe, 0, len(recipeDTOs))
	for _, dto := range recipeDTOs {
		if err := rs.validateSaveRecipeDTO(dto); err != nil {
			zap.L().Error("function.error", zap.String("func", "*recipeService.SaveMultipleRecipes"), zap.Error(err), zap.Any("params", __logParams))
			result0 = err
			return
		}

		recipe, err := rs.convertSaveRecipeDTOToModel(dto, userID)
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "*recipeService.SaveMultipleRecipes"), zap.Error(err), zap.Any("params", __logParams))
			result0 = fmt.Errorf("%w: %v", recipeDomain.ErrInvalidRecipeData, err)
			return
		}
		recipes = append(recipes, recipe)
	}

	if err := rs.recipeRepository.CreateMany(ctx, recipes); err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.SaveMultipleRecipes"), zap.Error(err), zap.Any("params", __logParams))
		result0 = err
		return
	}

	result0 = nil
	return
}

// GetRecipeByID retrieves a single recipe by ID
func (rs *recipeService) GetRecipeByID(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID) (result0 *recipeDTO.RecipeDetailDTO, result1 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "recipeID": recipeID, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.GetRecipeByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.GetRecipeByID"), zap.Any("params", __logParams))

	recipe, err := rs.recipeRepository.FindByID(ctx, recipeID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = recipeDomain.ErrRecipeNotFound
			return
		}
		zap.L().Error("function.error", zap.String("func", "*recipeService.GetRecipeByID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	result0 = rs.convertModelToRecipeDetailDTO(recipe)
	result1 = nil
	return
}

// GetUserRecipes retrieves all recipes for a user
func (rs *recipeService) GetUserRecipes(ctx context.Context, userID uuid.UUID) (result0 []*recipeDTO.RecipeDetailDTO, result1 error) {
	__logParams := map[string]any{"rs": rs, "ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.GetUserRecipes"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.GetUserRecipes"), zap.Any("params", __logParams))

	recipes, err := rs.recipeRepository.FindByUserID(ctx, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*recipeService.GetUserRecipes"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	recipeDTOs := make([]*recipeDTO.RecipeDetailDTO, 0, len(recipes))
	for _, recipe := range recipes {
		recipeDTOs = append(recipeDTOs, rs.convertModelToRecipeDetailDTO(recipe))
	}

	result0 = recipeDTOs
	result1 = nil
	return
}

// Helper methods

func (rs *recipeService) validateSaveRecipeDTO(dto *recipeDTO.SaveRecipeDTO) (result0 error) {
	__logParams := map[string]any{"rs": rs, "dto": dto}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.validateSaveRecipeDTO"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.validateSaveRecipeDTO"), zap.Any("params", __logParams))

	if dto == nil {
		result0 = fmt.Errorf("%w: recipe data is required", recipeDomain.ErrInvalidRequest)
		return
	}

	if dto.ID == "" {
		result0 = fmt.Errorf("%w: recipe ID is required", recipeDomain.ErrInvalidRequest)
		return
	}

	if _, err := uuid.Parse(dto.ID); err != nil {
		result0 = fmt.Errorf("%w: invalid recipe ID format", recipeDomain.ErrInvalidRequest)
		return
	}

	if dto.Title == "" {
		result0 = fmt.Errorf("%w: title is required", recipeDomain.ErrInvalidRequest)
		return
	}

	if len(dto.Ingredients) == 0 {
		result0 = fmt.Errorf("%w: at least one ingredient is required", recipeDomain.ErrInvalidRequest)
		return
	}

	if len(dto.Instructions) == 0 {
		result0 = fmt.Errorf("%w: at least one instruction is required", recipeDomain.ErrInvalidRequest)
		return
	}

	result0 = nil
	return
}

func (rs *recipeService) convertSaveRecipeDTOToModel(dto *recipeDTO.SaveRecipeDTO, userID uuid.UUID) (result0 *recipeModel.Recipe, result1 error) {
	__logParams := map[string]any{"rs": rs, "dto": dto, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.convertSaveRecipeDTOToModel"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.convertSaveRecipeDTOToModel"), zap.Any("params", __logParams))

	recipeID, err := uuid.Parse(dto.ID)
	if err != nil {
		result0 = nil
		result1 = err
		return
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

	result0 = recipe
	result1 = nil
	return
}

func (rs *recipeService) convertModelToRecipeDetailDTO(recipe *recipeModel.Recipe) (result0 *recipeDTO.RecipeDetailDTO) {
	__logParams := map[string]any{"rs": rs, "recipe": recipe}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*recipeService.convertModelToRecipeDetailDTO"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*recipeService.convertModelToRecipeDetailDTO"), zap.Any("params", __logParams))

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

	result0 = &recipeDTO.RecipeDetailDTO{
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
	return
}
