package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"
	recipeDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/domain"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

// RecipeHandler handles recipe-related HTTP requests
type RecipeHandler struct {
	recipeService recipeDomain.RecipeService
	llmService    *service.LLMServiceImpl
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(
	recipeService recipeDomain.RecipeService,
	llmService *service.LLMServiceImpl,
) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
		llmService:    llmService,
	}
}

func extractDetail(err error, base error, fallback string) string {
	if err == nil {
		return fallback
	}
	if base == nil || !errors.Is(err, base) {
		return fallback
	}
	message := err.Error()
	prefix := base.Error() + ": "
	if strings.HasPrefix(message, prefix) {
		detail := strings.TrimPrefix(message, prefix)
		if detail != "" {
			return detail
		}
	}
	return fallback
}

func (h *RecipeHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, recipeDomain.ErrInvalidRequest):
		detail := extractDetail(err, recipeDomain.ErrInvalidRequest, "Invalid recipe request")
		response.BadRequest(c, detail)
	case errors.Is(err, recipeDomain.ErrUnauthorized):
		response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
	case errors.Is(err, recipeDomain.ErrPantryNotFound):
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Pantry not found")
	case errors.Is(err, recipeDomain.ErrNoIngredients):
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "No ingredients available for this pantry")
	case errors.Is(err, recipeDomain.ErrInvalidLLMResponse):
		response.InternalError(c, "Received invalid response from recipe generator")
	case errors.Is(err, recipeDomain.ErrLLMRequest):
		response.InternalError(c, "Failed to process recipe request")
	default:
		response.InternalError(c, "Unexpected error while processing recipe request")
	}
}

// GenerateRecipe godoc
// @Summary Generate a recipe based on available ingredients
// @Description Generate a personalized recipe using LLM based on pantry ingredients and preferences
// @Tags recipes
// @Accept json
// @Produce json
// @Param request body dto.RecipeRequestDTO true "Recipe generation request"
// @Success 200 {object} response.Response{data=dto.RecipeResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/recipes/generate [post]
// @Security BearerAuth
func (h *RecipeHandler) GenerateRecipe(c *gin.Context) {
	var request dto.RecipeRequestDTO

	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "userID não encontrado no contexto")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		if userString, ok := userIDInterface.(string); ok {
			parsed, err := uuid.Parse(userString)
			if err != nil {
				response.Unauthorized(c, "user_id inválido no contexto")
				return
			}
			userID = parsed
		} else {
			response.Unauthorized(c, "user_id inválido no contexto")
			return
		}
	}

	recipe, err := h.recipeService.GenerateRecipe(c.Request.Context(), &request, userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.OK(c, recipe)
}

// GetAvailableIngredients godoc
// @Summary Get available ingredients from pantry with quantities
// @Description Get list of ingredients available in a specific pantry with quantity and unit information
// @Tags recipes
// @Produce json
// @Param pantry_id path string true "Pantry ID" format(uuid)
// @Success 200 {object} response.Response{data=[]service.IngredientInfo}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/recipes/pantries/{pantry_id}/ingredients [get]
// @Security BearerAuth
func (h *RecipeHandler) GetAvailableIngredients(c *gin.Context) {
	// Get pantry ID from URL
	pantryID := c.Param("id")
	if pantryID == "" {
		pantryID = c.Param("pantry_id") // fallback for recipe-specific endpoints
	}
	if pantryID == "" {
		response.BadRequest(c, "pantry_id não fornecido")
		return
	}

	// Validate pantry ID format
	pantryUUID, err := uuid.Parse(pantryID)
	if err != nil {
		response.BadRequest(c, "ID da despensa inválido: "+err.Error())
		return
	}

	// Get user ID from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "userID não encontrado no contexto")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		if userString, ok := userIDInterface.(string); ok {
			parsed, err := uuid.Parse(userString)
			if err != nil {
				response.Unauthorized(c, "user_id inválido no contexto")
				return
			}
			userID = parsed
		} else {
			response.Unauthorized(c, "user_id inválido no contexto")
			return
		}
	}

	ingredients, err := h.recipeService.GetAvailableIngredients(c.Request.Context(), pantryUUID, userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.OK(c, ingredients)
}

// ChatWithLLM godoc
// @Summary Chat with LLM for recipe suggestions
// @Description Send a direct message to LLM for recipe advice and cooking tips
// @Tags recipes,llm
// @Accept json
// @Produce json
// @Param request body dto.LLMRequestDTO true "LLM chat request"
// @Success 200 {object} response.Response{data=dto.LLMResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/recipes/chat [post]
// @Security BearerAuth
func (h *RecipeHandler) ChatWithLLM(c *gin.Context) {
	var request dto.LLMRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	// Process LLM request
	response_data, err := h.llmService.ProcessRequest(c.Request.Context(), &request)
	if err != nil {
		response.InternalError(c, "Erro na requisição ao LLM: "+err.Error())
		return
	}

	response.OK(c, response_data)
}

// GetLLMProviders godoc
// @Summary Get available LLM providers
// @Description Get list of available LLM providers and current active provider
// @Tags llm
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/providers [get]
// @Security BearerAuth
func (h *RecipeHandler) GetLLMProviders(c *gin.Context) {
	providers := h.llmService.GetAvailableProviders()
	currentProvider := h.llmService.GetCurrentProvider()

	providerInfo, err := h.llmService.GetProviderInfo()
	if err != nil {
		providerInfo = map[string]interface{}{
			"error": err.Error(),
		}
	}

	data := map[string]interface{}{
		"available_providers": providers,
		"current_provider":    currentProvider,
		"provider_info":       providerInfo,
	}

	response.OK(c, data)
}

// SetLLMProvider godoc
// @Summary Set active LLM provider
// @Description Change the active LLM provider
// @Tags llm
// @Accept json
// @Produce json
// @Param request body map[string]string true "Provider name" example({"provider": "openai"})
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/providers [put]
// @Security BearerAuth
func (h *RecipeHandler) SetLLMProvider(c *gin.Context) {
	var request map[string]string

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	provider, exists := request["provider"]
	if !exists || provider == "" {
		response.BadRequest(c, "provider não fornecido")
		return
	}

	if err := h.llmService.SetProvider(provider); err != nil {
		response.BadRequest(c, "Erro ao definir provedor: "+err.Error())
		return
	}

	response.OK(c, map[string]string{
		"provider": provider,
	})
}

// EstimateTokens godoc
// @Summary Estimate tokens for text
// @Description Estimate the number of tokens a text will use
// @Tags llm
// @Accept json
// @Produce json
// @Param request body map[string]string true "Text to analyze" example({"text": "Hello, world!"})
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/estimate-tokens [post]
// @Security BearerAuth
func (h *RecipeHandler) EstimateTokens(c *gin.Context) {
	var request map[string]string

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	text, exists := request["text"]
	if !exists {
		response.BadRequest(c, "text não fornecido")
		return
	}

	tokens, err := h.llmService.EstimateTokens(text)
	if err != nil {
		response.InternalError(c, "Erro ao estimar tokens: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"text":             text,
		"estimated_tokens": tokens,
		"characters":       len(text),
		"words":            len(strings.Fields(text)),
	}

	response.OK(c, data)
}
