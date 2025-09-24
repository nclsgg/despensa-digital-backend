package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/provider"
)

// LLMServiceImpl implementa a interface LLMService
type LLMServiceImpl struct {
	factory        *provider.ProviderFactory
	promptBuilder  domain.PromptBuilder
	configs        map[string]*model.LLMConfig
	activeProvider string
}

// NewLLMService cria uma nova instância do serviço LLM
func NewLLMService() *LLMServiceImpl {
	service := &LLMServiceImpl{
		factory:       provider.NewProviderFactory(),
		promptBuilder: NewPromptBuilder(),
		configs:       make(map[string]*model.LLMConfig),
	}

	// Auto-configure OpenAI if environment variables are available
	service.initializeDefaultProviders()

	return service
}

// AddProviderConfig adiciona uma configuração de provedor
func (s *LLMServiceImpl) AddProviderConfig(providerName string, config *model.LLMConfig) error {
	// Valida se o provedor é suportado
	if !s.factory.IsProviderSupported(config.Provider) {
		return fmt.Errorf("provedor '%s' não é suportado", config.Provider)
	}

	// Testa a configuração criando um provedor temporário
	_, err := s.factory.CreateProvider(config)
	if err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	s.configs[providerName] = config

	// Define como ativo se for o primeiro
	if s.activeProvider == "" {
		s.activeProvider = providerName
	}

	return nil
}

// ProcessRequest processa uma requisição genérica de LLM
func (s *LLMServiceImpl) ProcessRequest(ctx context.Context, request *dto.LLMRequestDTO) (*dto.LLMResponseDTO, error) {
	// Se um provedor específico foi especificado, usa ele
	if request.Provider != "" {
		return s.ProcessRequestWithProvider(ctx, request, request.Provider)
	}

	// Converte DTO para modelo interno
	llmRequest := request.ToLLMRequest()

	// Obtém o provedor ativo
	provider, err := s.getActiveProvider()
	if err != nil {
		return nil, err
	}

	// Define o modelo se não especificado
	if llmRequest.Model == "" {
		llmRequest.Model = provider.GetModel()
	}

	// Executa a requisição
	response, err := provider.Chat(ctx, llmRequest)
	if err != nil {
		return nil, err
	}

	// Converte resposta para DTO
	return dto.FromLLMResponse(response), nil
}

// ProcessRequestWithProvider processa uma requisição usando um provedor específico
func (s *LLMServiceImpl) ProcessRequestWithProvider(ctx context.Context, request *dto.LLMRequestDTO, providerName string) (*dto.LLMResponseDTO, error) {
	// Verifica se o provedor existe e está configurado
	config, exists := s.configs[providerName]
	if !exists {
		return nil, fmt.Errorf("provedor '%s' não está configurado", providerName)
	}

	// Cria o provedor temporariamente
	provider, err := s.factory.CreateProvider(config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar provedor '%s': %w", providerName, err)
	}

	// Converte DTO para modelo interno
	llmRequest := request.ToLLMRequest()

	// Define o modelo se não especificado
	if llmRequest.Model == "" {
		llmRequest.Model = provider.GetModel()
	}

	// Executa a requisição
	response, err := provider.Chat(ctx, llmRequest)
	if err != nil {
		return nil, err
	}

	// Converte resposta para DTO e adiciona informação do provedor usado
	responseDTO := dto.FromLLMResponse(response)
	if responseDTO.Metadata == nil {
		responseDTO.Metadata = make(map[string]string)
	}
	responseDTO.Metadata["used_provider"] = providerName

	return responseDTO, nil
}

// ProcessChatRequest processa uma requisição de chat simples
func (s *LLMServiceImpl) ProcessChatRequest(ctx context.Context, request *dto.ChatRequestDTO) (*dto.ChatResponseDTO, error) {
	// Cria uma requisição LLM baseada no chat
	llmRequest := &dto.LLMRequestDTO{
		Messages: []dto.MessageDTO{
			{
				Role:    "user",
				Content: request.Message,
			},
		},
		Provider: request.Provider,
	}

	// Adiciona contexto se fornecido
	if request.Context != "" {
		// Adiciona mensagem de sistema com contexto
		systemMessage := dto.MessageDTO{
			Role:    "system",
			Content: "Contexto: " + request.Context,
		}
		llmRequest.Messages = append([]dto.MessageDTO{systemMessage}, llmRequest.Messages...)
	}

	// Processa a requisição
	var response *dto.LLMResponseDTO
	var err error
	var usedProvider string

	if request.Provider != "" {
		response, err = s.ProcessRequestWithProvider(ctx, llmRequest, request.Provider)
		usedProvider = request.Provider
	} else {
		response, err = s.ProcessRequest(ctx, llmRequest)
		// Obtém o provedor ativo para incluir na resposta
		if s.activeProvider != "" {
			usedProvider = s.activeProvider
		}
	}

	if err != nil {
		return nil, err
	}

	// Converte para ChatResponseDTO
	return &dto.ChatResponseDTO{
		Response: response.Response,
		Provider: usedProvider,
		Model:    response.Model,
		Usage:    response.Usage,
	}, nil
}

// GenerateText gera texto baseado em um prompt
func (s *LLMServiceImpl) GenerateText(ctx context.Context, prompt string, options map[string]interface{}) (*dto.LLMResponseDTO, error) {
	// Cria mensagens básicas
	messages := []model.Message{
		{
			Role:    model.RoleUser,
			Content: prompt,
		},
	}

	// Cria requisição
	request := &model.LLMRequest{
		Messages: messages,
	}

	// Aplica opções
	if maxTokens, ok := options["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	}
	if temperature, ok := options["temperature"].(float64); ok {
		request.Temperature = temperature
	}
	if topP, ok := options["top_p"].(float64); ok {
		request.TopP = topP
	}

	// Obtém o provedor ativo
	provider, err := s.getActiveProvider()
	if err != nil {
		return nil, err
	}

	// Executa a requisição
	response, err := provider.Chat(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("erro na geração de texto: %w", err)
	}

	return dto.FromLLMResponse(response), nil
}

// BuildPrompt constrói um prompt usando um template
func (s *LLMServiceImpl) BuildPrompt(ctx context.Context, templateID string, variables map[string]string) (string, error) {
	// TODO: Implementar busca de template no repositório
	// Por enquanto, retorna erro
	return "", fmt.Errorf("busca de templates não implementada")
}

// GetAvailableProviders retorna os provedores disponíveis
func (s *LLMServiceImpl) GetAvailableProviders() []string {
	return s.factory.GetSupportedProviders()
}

// SetProvider define o provedor ativo
func (s *LLMServiceImpl) SetProvider(providerName string) error {
	if _, exists := s.configs[providerName]; !exists {
		return fmt.Errorf("provedor '%s' não configurado", providerName)
	}

	s.activeProvider = providerName
	return nil
}

// GetCurrentProvider retorna o provedor atual
func (s *LLMServiceImpl) GetCurrentProvider() string {
	return s.activeProvider
}

// getActiveProvider obtém o provedor ativo
func (s *LLMServiceImpl) getActiveProvider() (domain.LLMProvider, error) {
	if s.activeProvider == "" {
		return nil, fmt.Errorf("nenhum provedor ativo configurado")
	}

	config, exists := s.configs[s.activeProvider]
	if !exists {
		return nil, fmt.Errorf("configuração não encontrada para provedor '%s'", s.activeProvider)
	}

	return s.factory.CreateProvider(config)
}

// CreateChatRequest cria uma requisição de chat com system e user prompts
func (s *LLMServiceImpl) CreateChatRequest(systemPrompt, userPrompt string, options map[string]interface{}) *dto.LLMRequestDTO {
	messages := []dto.MessageDTO{}

	if systemPrompt != "" {
		messages = append(messages, dto.MessageDTO{
			Role:    string(model.RoleSystem),
			Content: systemPrompt,
		})
	}

	if userPrompt != "" {
		messages = append(messages, dto.MessageDTO{
			Role:    string(model.RoleUser),
			Content: userPrompt,
		})
	}

	request := &dto.LLMRequestDTO{
		Messages: messages,
	}

	// Aplica opções
	if maxTokens, ok := options["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	}
	if temperature, ok := options["temperature"].(float64); ok {
		request.Temperature = temperature
	}
	if topP, ok := options["top_p"].(float64); ok {
		request.TopP = topP
	}

	return request
}

// EstimateTokens estima tokens usando o provedor ativo
func (s *LLMServiceImpl) EstimateTokens(text string) (int, error) {
	provider, err := s.getActiveProvider()
	if err != nil {
		return 0, err
	}

	return provider.EstimateTokens(text), nil
}

// GenerateID gera um ID único para requisições
func (s *LLMServiceImpl) GenerateID() string {
	return uuid.New().String()
}

// GetProviderInfo retorna informações sobre o provedor ativo
func (s *LLMServiceImpl) GetProviderInfo() (map[string]interface{}, error) {
	if s.activeProvider == "" {
		return nil, fmt.Errorf("nenhum provedor ativo")
	}

	config, exists := s.configs[s.activeProvider]
	if !exists {
		return nil, fmt.Errorf("configuração não encontrada")
	}

	provider, err := s.getActiveProvider()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":     provider.GetProviderName(),
		"model":    provider.GetModel(),
		"provider": s.activeProvider,
		"timeout":  config.Timeout,
	}, nil
}

// initializeDefaultProviders configura automaticamente provedores com base nas variáveis de ambiente
func (s *LLMServiceImpl) initializeDefaultProviders() {
	// Verifica se OpenAI está configurado
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		openaiConfig := &model.LLMConfig{
			Provider:    model.ProviderOpenAI,
			APIKey:      apiKey,
			Model:       "gpt-3.5-turbo",
			MaxTokens:   2000,
			Temperature: 0.7,
			Timeout:     30 * time.Second,
		}

		err := s.AddProviderConfig("openai", openaiConfig)
		if err == nil {
			log.Println("OpenAI provider configurado automaticamente")
		} else {
			log.Printf("Erro ao configurar OpenAI provider automaticamente: %v", err)
		}
	}

	// Verifica se Gemini está configurado
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		geminiConfig := &model.LLMConfig{
			Provider:    model.ProviderGemini,
			APIKey:      apiKey,
			Model:       "gemini-1.5-flash",
			MaxTokens:   2000,
			Temperature: 0.7,
			Timeout:     30 * time.Second,
		}

		err := s.AddProviderConfig("gemini", geminiConfig)
		if err == nil {
			log.Println("Gemini provider configurado automaticamente")
		} else {
			log.Printf("Erro ao configurar Gemini provider automaticamente: %v", err)
		}
	}

	// Aqui você pode adicionar configurações automáticas para outros provedores
	// como Anthropic Claude, etc.
}
