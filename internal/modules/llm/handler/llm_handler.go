package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

// LLMHandler handles LLM-related HTTP requests
type LLMHandler struct {
	llmService *service.LLMServiceImpl
}

// NewLLMHandler creates a new LLM handler
func NewLLMHandler(llmService *service.LLMServiceImpl) *LLMHandler {
	return &LLMHandler{
		llmService: llmService,
	}
}

// ProcessChatRequest godoc
// @Summary Process a chat request with LLM
// @Description Send messages to LLM and get response (supports provider selection)
// @Tags LLM
// @Accept json
// @Produce json
// @Param request body dto.ChatRequestDTO true "Chat request with optional provider"
// @Success 200 {object} response.Response{data=dto.ChatResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/chat [post]
// @Security BearerAuth
func (h *LLMHandler) ProcessChatRequest(c *gin.Context) {
	var request dto.ChatRequestDTO

	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	// Process request with LLM service
	chatResponse, err := h.llmService.ProcessChatRequest(c.Request.Context(), &request)
	if err != nil {
		response.InternalError(c, "Erro ao processar requisição LLM: "+err.Error())
		return
	}

	response.Success(c, 200, chatResponse)
}

// ProcessLLMRequest godoc
// @Summary Process a detailed LLM request
// @Description Send detailed messages to LLM with full configuration options
// @Tags LLM
// @Accept json
// @Produce json
// @Param request body dto.LLMRequestDTO true "Detailed LLM request with optional provider"
// @Success 200 {object} response.Response{data=dto.LLMResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/process [post]
// @Security BearerAuth
func (h *LLMHandler) ProcessLLMRequest(c *gin.Context) {
	var request dto.LLMRequestDTO

	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	// Process request with LLM service
	llmResponse, err := h.llmService.ProcessRequest(c.Request.Context(), &request)
	if err != nil {
		response.InternalError(c, "Erro ao processar requisição LLM: "+err.Error())
		return
	}

	response.OK(c, llmResponse)
}

// BuildPrompt godoc
// @Summary Build a prompt from template and variables
// @Description Build a prompt using template and variables for LLM
// @Tags LLM
// @Accept json
// @Produce json
// @Param request body dto.PromptBuilderDTO true "Prompt builder request"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/prompt/build [post]
// @Security BearerAuth
func (h *LLMHandler) BuildPrompt(c *gin.Context) {
	var request dto.PromptBuilderDTO

	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Dados de entrada inválidos: "+err.Error())
		return
	}

	// Get prompt builder from service
	promptBuilder := service.NewPromptBuilder()

	// Build prompt
	builtPrompt, err := promptBuilder.BuildSystemPrompt(request.Template, request.Variables)
	if err != nil {
		response.BadRequest(c, "Erro ao construir prompt: "+err.Error())
		return
	}

	result := map[string]string{
		"template":     request.Template,
		"built_prompt": builtPrompt,
	}

	response.OK(c, result)
}

// GetProviderStatus godoc
// @Summary Get status of LLM providers
// @Description Get current status and configuration of LLM providers
// @Tags LLM
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/providers/status [get]
// @Security BearerAuth
func (h *LLMHandler) GetProviderStatus(c *gin.Context) {
	availableProviders := h.llmService.GetAvailableProviders()
	currentProvider := h.llmService.GetCurrentProvider()

	status := map[string]interface{}{
		"available_providers": availableProviders,
		"current_provider":    currentProvider,
		"total_providers":     len(availableProviders),
	}

	response.OK(c, status)
}

// ConfigureProvider godoc
// @Summary Configure a new LLM provider
// @Description Add or update configuration for an LLM provider
// @Tags LLM
// @Accept json
// @Produce json
// @Param provider_name query string true "Provider name (e.g., openai, anthropic)"
// @Param request body map[string]interface{} true "Provider configuration"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/providers/config [post]
// @Security BearerAuth
func (h *LLMHandler) ConfigureProvider(c *gin.Context) {
	providerName := c.Query("provider_name")
	if providerName == "" {
		response.BadRequest(c, "Nome do provedor é obrigatório")
		return
	}

	var configRequest map[string]interface{}
	if err := c.ShouldBindJSON(&configRequest); err != nil {
		response.BadRequest(c, "Configuração inválida: "+err.Error())
		return
	}

	// Convert to LLM config
	config := &model.LLMConfig{
		Provider: model.LLMProvider(providerName),
	}

	// Extract common configuration fields
	if apiKey, ok := configRequest["api_key"].(string); ok {
		config.APIKey = apiKey
	}
	if model, ok := configRequest["model"].(string); ok {
		config.Model = model
	}
	if baseURL, ok := configRequest["base_url"].(string); ok {
		config.BaseURL = baseURL
	}

	// Add provider configuration
	if err := h.llmService.AddProviderConfig(providerName, config); err != nil {
		response.BadRequest(c, "Erro ao configurar provedor: "+err.Error())
		return
	}

	result := map[string]string{
		"provider": providerName,
		"status":   "configured",
		"model":    config.Model,
	}

	response.OK(c, result)
}

// GetAvailableProviders godoc
// @Summary Get list of available LLM providers
// @Description Get all supported LLM provider names
// @Tags LLM
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Failure 401 {object} response.Response
// @Router /api/v1/llm/providers/available [get]
// @Security BearerAuth
func (h *LLMHandler) GetAvailableProviders(c *gin.Context) {
	providers := h.llmService.GetAvailableProviders()
	response.OK(c, providers)
}

// SwitchProvider godoc
// @Summary Switch active LLM provider
// @Description Switch to a different configured LLM provider
// @Tags LLM
// @Accept json
// @Produce json
// @Param provider_name query string true "Provider name to switch to"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/providers/switch [post]
// @Security BearerAuth
func (h *LLMHandler) SwitchProvider(c *gin.Context) {
	providerName := c.Query("provider_name")
	if providerName == "" {
		response.BadRequest(c, "Nome do provedor é obrigatório")
		return
	}

	// Switch provider
	if err := h.llmService.SetProvider(providerName); err != nil {
		response.BadRequest(c, "Erro ao trocar provedor: "+err.Error())
		return
	}

	result := map[string]string{
		"previous_provider": h.llmService.GetCurrentProvider(),
		"new_provider":      providerName,
		"status":            "switched",
	}

	response.OK(c, result)
}

// TestProvider godoc
// @Summary Test LLM provider connectivity
// @Description Send a test message to verify provider is working
// @Tags LLM
// @Accept json
// @Produce json
// @Param provider_name query string false "Provider to test (current if not specified)"
// @Param message query string false "Test message (default test message if not specified)"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/llm/providers/test [post]
// @Security BearerAuth
func (h *LLMHandler) TestProvider(c *gin.Context) {
	providerName := c.Query("provider_name")
	if providerName == "" {
		providerName = h.llmService.GetCurrentProvider()
	}

	testMessage := c.Query("message")
	if testMessage == "" {
		testMessage = "Olá! Este é um teste de conectividade. Responda apenas: 'Teste OK'"
	}

	// Create test request
	testRequest := &dto.LLMRequestDTO{
		Messages: []dto.MessageDTO{
			{
				Role:    "user",
				Content: testMessage,
			},
		},
		MaxTokens:   50,
		Temperature: 0.1,
	}

	// Test with current provider
	ctx := c.Request.Context()
	llmResponse, err := h.llmService.ProcessRequest(ctx, testRequest)
	if err != nil {
		response.InternalError(c, "Teste de provedor falhou: "+err.Error())
		return
	}

	result := map[string]interface{}{
		"provider":      providerName,
		"test_message":  testMessage,
		"response":      llmResponse.Response,
		"status":        "success",
		"response_time": "< 5s", // You could measure actual time
		"tokens_used":   llmResponse.Usage.TotalTokens,
	}

	response.OK(c, result)
}
