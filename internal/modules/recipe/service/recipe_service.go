package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	llmSvc "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"
	pantryDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
)

// RecipeServiceImpl implementa a interface RecipeService
type RecipeServiceImpl struct {
	llmService     *llmSvc.LLMServiceImpl
	itemRepository domain.ItemRepository
	pantryService  pantryDomain.PantryService
	promptBuilder  *llmSvc.PromptBuilderImpl
}

// NewRecipeService cria uma nova instância do serviço de receitas
func NewRecipeService(
	llmService *llmSvc.LLMServiceImpl,
	itemRepository domain.ItemRepository,
	pantryService pantryDomain.PantryService,
) *RecipeServiceImpl {
	return &RecipeServiceImpl{
		llmService:     llmService,
		itemRepository: itemRepository,
		pantryService:  pantryService,
		promptBuilder:  llmSvc.NewPromptBuilder(),
	}
}

// GenerateRecipe gera uma receita baseada nos parâmetros
func (rs *RecipeServiceImpl) GenerateRecipe(ctx context.Context, request *dto.RecipeRequestDTO, userID string) (*dto.RecipeResponseDTO, error) {
	// Valida a requisição
	if err := rs.ValidateRecipeRequest(request); err != nil {
		return nil, fmt.Errorf("requisição inválida: %w", err)
	}

	// Obtém ingredientes disponíveis na despensa
	availableIngredients, err := rs.GetAvailableIngredients(ctx, request.PantryID, userID)
	fmt.Println("Available Ingredients:", availableIngredients)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter ingredientes: %w", err)
	}

	if len(availableIngredients) == 0 {
		return nil, fmt.Errorf("nenhum ingrediente disponível na despensa")
	}

	// Prepara variáveis para o prompt
	variables := rs.buildPromptVariables(request, availableIngredients)

	// Obtém templates de prompt
	templates := llmSvc.GetRecipePromptTemplates()

	// Constrói prompts
	systemPrompt, err := rs.promptBuilder.BuildSystemPrompt(templates.SystemPrompt, variables)
	fmt.Println("System Prompt:", systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("erro ao construir system prompt: %w", err)
	}

	userPrompt, err := rs.promptBuilder.BuildUserPrompt(templates.UserPrompt, variables)
	fmt.Println("User Prompt:", userPrompt)
	if err != nil {
		return nil, fmt.Errorf("erro ao construir user prompt: %w", err)
	}

	// Cria opções para o LLM
	options := map[string]interface{}{
		"max_tokens":  2000,
		"temperature": 0.7,
		"top_p":       0.9,
	}

	// Cria requisição de chat
	llmRequest := rs.llmService.CreateChatRequest(systemPrompt, userPrompt, options)
	fmt.Println("LLM Request Messages:", llmRequest)

	// Executa requisição ao LLM usando o provider especificado ou o ativo
	var llmResponse *dto.LLMResponseDTO

	if request.Provider != "" {
		// Usa provider específico da requisição
		llmResponse, err = rs.llmService.ProcessRequestWithProvider(ctx, llmRequest, request.Provider)
	} else {
		// Usa provider ativo
		llmResponse, err = rs.llmService.ProcessRequest(ctx, llmRequest)
	}

	if err != nil {
		return nil, fmt.Errorf("erro na requisição ao LLM: %w", err)
	}

	// Parseia a resposta JSON do LLM
	recipe, err := rs.parseRecipeResponse(llmResponse.Response)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta do LLM: %w", err)
	}

	// Enriquece a receita com informações adicionais
	recipe.ID = uuid.New().String()
	recipe.GeneratedAt = time.Now().Format(time.RFC3339)

	// Marca ingredientes como disponíveis ou não
	rs.markIngredientAvailability(recipe, availableIngredients)

	return recipe, nil
}

// IngredientInfo contém informações detalhadas de um ingrediente
type IngredientInfo struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// GetAvailableIngredients obtém ingredientes disponíveis na despensa com quantidades
func (rs *RecipeServiceImpl) GetAvailableIngredients(ctx context.Context, pantryID string, userID string) ([]IngredientInfo, error) {
	pantryUUID, err := uuid.Parse(pantryID)
	if err != nil {
		return nil, fmt.Errorf("ID da despensa inválido: %w", err)
	}

	// Obtém itens da despensa
	items, err := rs.itemRepository.ListByPantryID(ctx, pantryUUID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter itens da despensa: %w", err)
	}

	// Extrai informações detalhadas dos ingredientes
	ingredients := make([]IngredientInfo, 0, len(items))
	for _, item := range items {
		if item.Quantity > 0 { // Só inclui itens com quantidade disponível
			ingredients = append(ingredients, IngredientInfo{
				Name:     item.Name,
				Quantity: item.Quantity,
				Unit:     item.Unit,
			})
		}
	}

	return ingredients, nil
}

// SearchRecipesByIngredients busca receitas por ingredientes (placeholder)
func (rs *RecipeServiceImpl) SearchRecipesByIngredients(ctx context.Context, ingredients []string, filters map[string]string) ([]dto.RecipeResponseDTO, error) {
	// TODO: Implementar busca mais inteligente de receitas similares
	return nil, fmt.Errorf("busca de receitas por ingredientes não implementada")
}

// ValidateRecipeRequest valida uma requisição de receita
func (rs *RecipeServiceImpl) ValidateRecipeRequest(request *dto.RecipeRequestDTO) error {
	if request.PantryID == "" {
		return fmt.Errorf("ID da despensa é obrigatório")
	}

	if _, err := uuid.Parse(request.PantryID); err != nil {
		return fmt.Errorf("ID da despensa inválido")
	}

	if request.CookingTime < 0 || request.CookingTime > 480 {
		return fmt.Errorf("tempo de cozimento deve estar entre 0 e 480 minutos")
	}

	if request.ServingSize < 0 || request.ServingSize > 20 {
		return fmt.Errorf("número de porções deve estar entre 1 e 20")
	}

	validMealTypes := map[string]bool{
		"breakfast": true,
		"lunch":     true,
		"dinner":    true,
		"snack":     true,
		"dessert":   true,
		"":          true, // vazio é válido
	}

	if !validMealTypes[request.MealType] {
		return fmt.Errorf("tipo de refeição inválido")
	}

	validDifficulties := map[string]bool{
		"easy":   true,
		"medium": true,
		"hard":   true,
		"":       true, // vazio é válido
	}

	if !validDifficulties[request.Difficulty] {
		return fmt.Errorf("dificuldade inválida")
	}

	return nil
}

// EnrichRecipeWithNutrition adiciona informações nutricionais (placeholder)
func (rs *RecipeServiceImpl) EnrichRecipeWithNutrition(ctx context.Context, recipe *dto.RecipeResponseDTO) error {
	// TODO: Implementar cálculo nutricional real
	return nil
}

// buildPromptVariables constrói as variáveis para o prompt
func (rs *RecipeServiceImpl) buildPromptVariables(request *dto.RecipeRequestDTO, ingredients []IngredientInfo) map[string]string {
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
		variables["meal_type"] = request.MealType
	} else {
		variables["meal_type"] = "qualquer"
	}

	if request.Difficulty != "" {
		variables["difficulty"] = request.Difficulty
	} else {
		variables["difficulty"] = "qualquer"
	}

	if request.ServingSize > 0 {
		variables["serving_size"] = fmt.Sprintf("%d", request.ServingSize)
	} else {
		variables["serving_size"] = "4"
	}

	if request.Cuisine != "" {
		variables["cuisine"] = request.Cuisine
	}

	if len(request.DietaryRestrictions) > 0 {
		variables["dietary_restrictions"] = strings.Join(request.DietaryRestrictions, ", ")
	}

	if request.Purpose != "" {
		variables["purpose"] = request.Purpose
	}

	if request.AdditionalNotes != "" {
		variables["additional_notes"] = request.AdditionalNotes
	}

	return variables
}

// parseRecipeResponse parseia a resposta JSON do LLM
func (rs *RecipeServiceImpl) parseRecipeResponse(response string) (*dto.RecipeResponseDTO, error) {
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
	var recipe dto.RecipeResponseDTO
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
func (rs *RecipeServiceImpl) fixCommonJSONIssues(jsonStr string) string {
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
func (rs *RecipeServiceImpl) markIngredientAvailability(recipe *dto.RecipeResponseDTO, availableIngredients []IngredientInfo) {
	// Cria um mapa para busca rápida
	availableMap := make(map[string]bool)
	for _, ingredient := range availableIngredients {
		availableMap[strings.ToLower(ingredient.Name)] = true
	}

	// Marca disponibilidade para cada ingrediente da receita
	for i := range recipe.Ingredients {
		ingredientName := strings.ToLower(recipe.Ingredients[i].Name)
		recipe.Ingredients[i].Available = availableMap[ingredientName]
	}
}
