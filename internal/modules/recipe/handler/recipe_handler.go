package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	creditsDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/credits/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"
	recipeDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/domain"
	recipeDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/dto"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

// RecipeHandler handles recipe-related HTTP requests
type RecipeHandler struct {
	recipeService recipeDomain.RecipeService
	llmService    *service.LLMServiceImpl
	creditService creditsDomain.CreditService
}

func captureContextFields(c *gin.Context) map[string]any {
	if c == nil {
		return nil
	}

	fields := map[string]any{
		"client_ip": c.ClientIP(),
	}

	if req := c.Request; req != nil {
		fields["method"] = req.Method
		fields["host"] = req.Host
		fields["url_path"] = req.URL.Path
		if raw := req.URL.RawQuery; raw != "" {
			fields["query"] = raw
		}
	}

	if route := c.FullPath(); route != "" {
		fields["route"] = route
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields["request_id"] = requestID
	} else if requestID := c.GetString("requestID"); requestID != "" {
		fields["request_id"] = requestID
	}

	return fields
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(
	recipeService recipeDomain.RecipeService,
	llmService *service.LLMServiceImpl,
	creditService creditsDomain.CreditService,
) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
		llmService:    llmService,
		creditService: creditService,
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
	case errors.Is(err, recipeDomain.ErrRecipeNotFound):
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Recipe not found")
	case errors.Is(err, recipeDomain.ErrInvalidRecipeData):
		detail := extractDetail(err, recipeDomain.ErrInvalidRecipeData, "Invalid recipe data")
		response.BadRequest(c, detail)
	default:
		response.InternalError(c, "Unexpected error while processing recipe request")
	}
}

// GenerateRecipe godoc
// @Summary Generate a recipe based on available ingredients
// @Description Generate 3 personalized recipes using LLM based on pantry ingredients and preferences
// @Tags recipes
// @Accept json
// @Produce json
// @Param request body dto.RecipeRequestDTO true "Recipe generation request"
// @Success 200 {object} response.Response{data=[]dto.RecipeResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/recipes/generate [post]
// @Security BearerAuth
func (h *RecipeHandler) GenerateRecipe(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())
	var request dto.RecipeRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("Invalid recipe generation request",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.Error(err),
		)
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

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
				logger.Error("Invalid user ID format",
					zap.String(appLogger.FieldModule, "recipe"),
					zap.String(appLogger.FieldFunction, "GenerateRecipe"),
					zap.Error(err),
				)
				response.Unauthorized(c, "user_id inválido no contexto")
				return
			}
			userID = parsed
		} else {
			response.Unauthorized(c, "user_id inválido no contexto")
			return
		}
	}

	// Generate 3 recipes as per requirements
	recipes, err := h.recipeService.GenerateMultipleRecipes(c.Request.Context(), &request, userID, 3)
	if err != nil {
		logger.Error("Failed to generate recipes",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GenerateRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		h.handleServiceError(c, err)
		return
	}

	if creditErr := h.creditService.ConsumeCredit(c.Request.Context(), userID, "AI request - recipe generation"); creditErr != nil {
		switch {
		case errors.Is(creditErr, creditsDomain.ErrInsufficientCredits):
			logger.Warn("Insufficient credits for recipe generation",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "GenerateRecipe"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			response.Fail(c, http.StatusPaymentRequired, "INSUFFICIENT_CREDITS", "You don't have enough credits to generate a recipe")
		default:
			logger.Error("Failed to consume credit",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "GenerateRecipe"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(creditErr),
			)
			response.InternalError(c, "Failed to consume credit after recipe generation")
		}
		return
	}

	logger.Info("Recipes generated successfully via handler",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GenerateRecipe"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(recipes)),
	)

	response.OK(c, map[string]interface{}{
		"message": "3 receitas foram geradas com sucesso.",
		"recipes": recipes,
		"count":   len(recipes),
	})
}

func (h *RecipeHandler) GetAvailableIngredients(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID := c.Param("id")
	if pantryID == "" {
		pantryID = c.Param("pantry_id")
	}
	if pantryID == "" {
		logger.Warn("Pantry ID not provided",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
		)
		response.BadRequest(c, "pantry_id não fornecido")
		return
	}

	pantryUUID, err := uuid.Parse(pantryID)
	if err != nil {
		logger.Warn("Invalid pantry ID format",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
			zap.String("pantry_id", pantryID),
			zap.Error(err),
		)
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
				logger.Error("Invalid user ID format",
					zap.String(appLogger.FieldModule, "recipe"),
					zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
					zap.Error(err),
				)
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
		logger.Error("Failed to get available ingredients",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryUUID.String()),
			zap.Error(err),
		)
		h.handleServiceError(c, err)
		return
	}

	logger.Info("Available ingredients retrieved via handler",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GetAvailableIngredients"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryUUID.String()),
		zap.Int(appLogger.FieldCount, len(ingredients)),
	)

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
	logger := appLogger.FromContext(c.Request.Context())
	var request dto.LLMRequestDTO

	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "userID não encontrado no contexto")
		return
	}

	var userID uuid.UUID
	switch v := userIDInterface.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			response.Unauthorized(c, "user_id inválido no contexto")
			return
		}
		userID = parsed
	default:
		response.Unauthorized(c, "user_id inválido no contexto")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("Invalid LLM chat request",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "ChatWithLLM"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	response_data, err := h.llmService.ProcessRequest(c.Request.Context(), &request)
	if err != nil {
		logger.Error("LLM request failed",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "ChatWithLLM"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Erro na requisição ao LLM: "+err.Error())
		return
	}

	if creditErr := h.creditService.ConsumeCredit(c.Request.Context(), userID, "AI request - recipe chat"); creditErr != nil {
		switch {
		case errors.Is(creditErr, creditsDomain.ErrInsufficientCredits):
			logger.Warn("Insufficient credits for recipe chat",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "ChatWithLLM"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			response.Fail(c, http.StatusPaymentRequired, "INSUFFICIENT_CREDITS", "You don't have enough credits to chat with the recipe assistant")
		default:
			logger.Error("Failed to consume credit",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "ChatWithLLM"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(creditErr),
			)
			response.InternalError(c, "Failed to consume credits for recipe chat")
		}
		return
	}

	logger.Info("LLM chat completed successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "ChatWithLLM"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)

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
	logger := appLogger.FromContext(c.Request.Context())

	providers := h.llmService.GetAvailableProviders()
	currentProvider := h.llmService.GetCurrentProvider()

	providerInfo, err := h.llmService.GetProviderInfo()
	if err != nil {
		logger.Warn("Failed to get provider info",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetLLMProviders"),
			zap.Error(err),
		)
		providerInfo = map[string]interface{}{
			"error": err.Error(),
		}
	}

	data := map[string]interface{}{
		"available_providers": providers,
		"current_provider":    currentProvider,
		"provider_info":       providerInfo,
	}

	logger.Info("LLM providers retrieved",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GetLLMProviders"),
		zap.String("current_provider", currentProvider),
	)

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
	logger := appLogger.FromContext(c.Request.Context())
	var request map[string]string

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("Invalid provider request",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SetLLMProvider"),
			zap.Error(err),
		)
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	provider, exists := request["provider"]
	if !exists || provider == "" {
		logger.Warn("Provider not specified",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SetLLMProvider"),
		)
		response.BadRequest(c, "provider não fornecido")
		return
	}

	if err := h.llmService.SetProvider(provider); err != nil {
		logger.Error("Failed to set provider",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SetLLMProvider"),
			zap.String("provider", provider),
			zap.Error(err),
		)
		response.BadRequest(c, "Erro ao definir provedor: "+err.Error())
		return
	}

	logger.Info("LLM provider changed",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "SetLLMProvider"),
		zap.String("provider", provider),
	)

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
	logger := appLogger.FromContext(c.Request.Context())
	var request map[string]string

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("Invalid token estimation request",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "EstimateTokens"),
			zap.Error(err),
		)
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	text, exists := request["text"]
	if !exists {
		logger.Warn("Text not provided for estimation",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "EstimateTokens"),
		)
		response.BadRequest(c, "text não fornecido")
		return
	}

	tokens, err := h.llmService.EstimateTokens(text)
	if err != nil {
		logger.Error("Failed to estimate tokens",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "EstimateTokens"),
			zap.Int("text_length", len(text)),
			zap.Error(err),
		)
		response.InternalError(c, "Erro ao estimar tokens: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"text":             text,
		"estimated_tokens": tokens,
		"characters":       len(text),
		"words":            len(strings.Fields(text)),
	}

	logger.Info("Tokens estimated",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "EstimateTokens"),
		zap.Int("estimated_tokens", tokens),
	)

	response.OK(c, data)
}

// SaveRecipe godoc
// @Summary Save a generated recipe
// @Description Save one or more generated recipes to the database
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipes body []dto.SaveRecipeDTO true "Recipe(s) to save"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/recipes/save [post]
// @Security BearerAuth
func (h *RecipeHandler) SaveRecipe(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

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
				logger.Error("Invalid user ID format",
					zap.String(appLogger.FieldModule, "recipe"),
					zap.String(appLogger.FieldFunction, "SaveRecipe"),
					zap.Error(err),
				)
				response.Unauthorized(c, "user_id inválido no contexto")
				return
			}
			userID = parsed
		} else {
			response.Unauthorized(c, "user_id inválido no contexto")
			return
		}
	}

	// Try to parse as array first
	var recipesArray []recipeDTO.SaveRecipeDTO
	if err := c.ShouldBindJSON(&recipesArray); err == nil && len(recipesArray) > 0 {
		// Multiple recipes
		recipePtrs := make([]*recipeDTO.SaveRecipeDTO, len(recipesArray))
		for i := range recipesArray {
			recipePtrs[i] = &recipesArray[i]
		}

		if err := h.recipeService.SaveMultipleRecipes(c.Request.Context(), recipePtrs, userID); err != nil {
			logger.Error("Failed to save multiple recipes",
				zap.String(appLogger.FieldModule, "recipe"),
				zap.String(appLogger.FieldFunction, "SaveRecipe"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Int(appLogger.FieldCount, len(recipesArray)),
				zap.Error(err),
			)
			h.handleServiceError(c, err)
			return
		}

		logger.Info("Multiple recipes saved successfully",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Int(appLogger.FieldCount, len(recipesArray)),
		)

		message := "Receitas salvas com sucesso."
		if len(recipesArray) == 1 {
			message = "Receita salva com sucesso."
		}
		response.OK(c, map[string]interface{}{
			"message": message,
			"count":   len(recipesArray),
		})
		return
	}

	// Try to parse as single recipe
	var singleRecipe recipeDTO.SaveRecipeDTO
	if err := c.ShouldBindJSON(&singleRecipe); err != nil {
		logger.Warn("Invalid recipe save request",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	if err := h.recipeService.SaveRecipe(c.Request.Context(), &singleRecipe, userID); err != nil {
		logger.Error("Failed to save recipe",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "SaveRecipe"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		h.handleServiceError(c, err)
		return
	}

	logger.Info("Recipe saved successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "SaveRecipe"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)

	response.OK(c, map[string]interface{}{
		"message": "Receita salva com sucesso.",
	})
}

// GetRecipes godoc
// @Summary Get all saved recipes
// @Description Get all recipes saved by the logged-in user
// @Tags recipes
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.RecipeDetailDTO}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/recipes [get]
// @Security BearerAuth
func (h *RecipeHandler) GetRecipes(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

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
				logger.Error("Invalid user ID format",
					zap.String(appLogger.FieldModule, "recipe"),
					zap.String(appLogger.FieldFunction, "GetRecipes"),
					zap.Error(err),
				)
				response.Unauthorized(c, "user_id inválido no contexto")
				return
			}
			userID = parsed
		} else {
			response.Unauthorized(c, "user_id inválido no contexto")
			return
		}
	}

	recipes, err := h.recipeService.GetUserRecipes(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user recipes",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetRecipes"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		h.handleServiceError(c, err)
		return
	}

	if len(recipes) == 0 {
		logger.Info("No recipes found for user",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetRecipes"),
			zap.String(appLogger.FieldUserID, userID.String()),
		)
		response.OK(c, map[string]interface{}{
			"message": "Nenhuma receita encontrada.",
			"recipes": []interface{}{},
		})
		return
	}

	logger.Info("User recipes retrieved successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GetRecipes"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(recipes)),
	)

	response.OK(c, recipes)
}

// GetRecipeByID godoc
// @Summary Get a specific recipe
// @Description Get a specific recipe by ID belonging to the logged-in user
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} response.Response{data=dto.RecipeDetailDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/recipes/{id} [get]
// @Security BearerAuth
func (h *RecipeHandler) GetRecipeByID(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	recipeID := c.Param("id")
	if recipeID == "" {
		logger.Warn("Recipe ID not provided",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetRecipeByID"),
		)
		response.BadRequest(c, "ID da receita não fornecido")
		return
	}

	recipeUUID, err := uuid.Parse(recipeID)
	if err != nil {
		logger.Warn("Invalid recipe ID format",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetRecipeByID"),
			zap.String("recipe_id", recipeID),
			zap.Error(err),
		)
		response.BadRequest(c, "ID da receita inválido: "+err.Error())
		return
	}

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
				logger.Error("Invalid user ID format",
					zap.String(appLogger.FieldModule, "recipe"),
					zap.String(appLogger.FieldFunction, "GetRecipeByID"),
					zap.Error(err),
				)
				response.Unauthorized(c, "user_id inválido no contexto")
				return
			}
			userID = parsed
		} else {
			response.Unauthorized(c, "user_id inválido no contexto")
			return
		}
	}

	recipe, err := h.recipeService.GetRecipeByID(c.Request.Context(), recipeUUID, userID)
	if err != nil {
		logger.Error("Failed to get recipe by ID",
			zap.String(appLogger.FieldModule, "recipe"),
			zap.String(appLogger.FieldFunction, "GetRecipeByID"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("recipe_id", recipeUUID.String()),
			zap.Error(err),
		)
		h.handleServiceError(c, err)
		return
	}

	logger.Info("Recipe retrieved by ID successfully",
		zap.String(appLogger.FieldModule, "recipe"),
		zap.String(appLogger.FieldFunction, "GetRecipeByID"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("recipe_id", recipeUUID.String()),
	)

	response.OK(c, recipe)
}
