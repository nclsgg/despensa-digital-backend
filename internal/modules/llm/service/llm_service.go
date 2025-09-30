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
	"go.uber.org/zap"
)

// LLMServiceImpl implementa a interface LLMService
type LLMServiceImpl struct {
	factory        *provider.ProviderFactory
	promptBuilder  domain.PromptBuilder
	configs        map[string]*model.LLMConfig
	activeProvider string
}

// NewLLMService cria uma nova instância do serviço LLM
func NewLLMService() (result0 *LLMServiceImpl) {
	__logParams := map[string]any{}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewLLMService"), zap.Any("result", result0), zap.Duration("duration", time.

			// Auto-configure OpenAI if environment variables are available
			Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewLLMService"), zap.Any("params", __logParams))
	service := &LLMServiceImpl{
		factory:       provider.NewProviderFactory(),
		promptBuilder: NewPromptBuilder(),
		configs:       make(map[string]*model.LLMConfig),
	}

	service.initializeDefaultProviders()
	result0 = service
	return
}

// AddProviderConfig adiciona uma configuração de provedor
func (s *LLMServiceImpl) AddProviderConfig(providerName string, config *model.LLMConfig) (result0 error) {
	__logParams :=
		// Valida se o provedor é suportado
		map[string]any{"s": s, "providerName": providerName, "config": config}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.AddProviderConfig"),

			// Testa a configuração criando um provedor temporário
			zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.AddProviderConfig"), zap.Any("params", __logParams))

	if !s.factory.IsProviderSupported(config.Provider) {
		result0 = fmt.Errorf("provedor '%s' não é suportado", config.Provider)
		return
	}

	_, err := s.factory.CreateProvider(config)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*LLMServiceImpl.AddProviderConfig"), zap.Error(err), zap.

			// Define como ativo se for o primeiro
			Any("params", __logParams))
		result0 = fmt.Errorf("configuração inválida: %w", err)
		return
	}

	s.configs[providerName] = config

	if s.activeProvider == "" {
		s.activeProvider = providerName
	}
	result0 = nil
	return
}

// ProcessRequest processa uma requisição genérica de LLM
func (s *LLMServiceImpl) ProcessRequest(ctx context.Context, request *dto.LLMRequestDTO) (result0 *dto.LLMResponseDTO, result1 error) {
	__logParams :=
		// Se um provedor específico foi especificado, usa ele
		map[string]any{"s": s, "ctx": ctx, "request": request}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.ProcessRequest"),

			// Converte DTO para modelo interno
			zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.

				// Obtém o provedor ativo
				Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.ProcessRequest"),

		// Define o modelo se não especificado
		zap.Any("params", __logParams))

	if request.Provider != "" {
		result0, result1 = s.ProcessRequestWithProvider(ctx, request, request.Provider)
		return
	}

	llmRequest := request.ToLLMRequest()

	provider, err := s.getActiveProvider()
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*LLMServiceImpl.ProcessRequest"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	if llmRequest.Model == "" {
		llmRequest.Model = provider.GetModel()
	}

	// Executa a requisição
	response, err := provider.Chat(ctx, llmRequest)
	if err != nil {
		zap.L().Error("function.error",

			// Converte resposta para DTO
			zap.String("func", "*LLMServiceImpl.ProcessRequest"), zap.Error(err), zap.Any("params",

				// ProcessRequestWithProvider processa uma requisição usando um provedor específico
				__logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = dto.FromLLMResponse(response)
	result1 = nil
	return
}

func (s *LLMServiceImpl) ProcessRequestWithProvider(ctx context.Context, request *dto.LLMRequestDTO, providerName string) (result0 *dto.LLMResponseDTO, result1 error) {
	__logParams :=
		// Verifica se o provedor existe e está configurado
		map[string]any{"s": s, "ctx": ctx, "request": request, "providerName": providerName}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.ProcessRequestWithProvider"),

			// Cria o provedor temporariamente
			zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.ProcessRequestWithProvider"),

		// Converte DTO para modelo interno
		zap.Any("params", __logParams))

	config, exists := s.configs[providerName]
	if !exists {
		result0 = nil
		result1 = fmt.Errorf("provedor '%s' não está configurado", providerName)
		return
	}

	provider, err := s.factory.CreateProvider(config)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*LLMServiceImpl.ProcessRequestWithProvider"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao criar provedor '%s': %w", providerName, err)
		return
	}

	llmRequest := request.ToLLMRequest()

	// Define o modelo se não especificado
	if llmRequest.Model == "" {
		llmRequest.Model = provider.GetModel()
	}

	// Executa a requisição
	response, err := provider.Chat(ctx, llmRequest)
	if err != nil {
		zap.L().Error("function.error",

			// Converte resposta para DTO e adiciona informação do provedor usado
			zap.String("func", "*LLMServiceImpl.ProcessRequestWithProvider"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	responseDTO := dto.FromLLMResponse(response)
	if responseDTO.Metadata == nil {
		responseDTO.Metadata = make(map[string]string)
	}
	responseDTO.Metadata["used_provider"] = providerName
	result0 = responseDTO
	result1 = nil
	return
}

// ProcessChatRequest processa uma requisição de chat simples
func (s *LLMServiceImpl) ProcessChatRequest(ctx context.Context, request *dto.ChatRequestDTO) (result0 *dto.ChatResponseDTO, result1 error) {
	__logParams :=
		// Cria uma requisição LLM baseada no chat
		map[string]any{"s": s, "ctx": ctx, "request": request}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.ProcessChatRequest"), zap.Any("result", map[string]any{"result0": result0,

			// Adiciona contexto se fornecido
			"result1": result1}), zap.Duration("duration", time.Since(__logStart)))

		// Adiciona mensagem de sistema com contexto
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.ProcessChatRequest"), zap.Any("params", __logParams))

	llmRequest := &dto.LLMRequestDTO{
		Messages: []dto.MessageDTO{
			{
				Role:    "user",
				Content: request.Message,
			},
		},
		Provider: request.Provider,
	}

	if request.Context != "" {

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
		zap.L().Error("function.error",

			// Converte para ChatResponseDTO
			zap.String("func", "*LLMServiceImpl.ProcessChatRequest"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &dto.ChatResponseDTO{
		Response: response.Response,
		Provider: usedProvider,
		Model:    response.Model,
		Usage:    response.Usage,
	}
	result1 = nil
	return
}

// GenerateText gera texto baseado em um prompt
func (s *LLMServiceImpl) GenerateText(ctx context.Context, prompt string, options map[string]interface{}) (result0 *dto.LLMResponseDTO, result1 error) {
	__logParams :=
		// Cria mensagens básicas
		map[string]any{"s": s, "ctx": ctx, "prompt": prompt, "options": options}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit",

			// Cria requisição
			zap.String("func", "*LLMServiceImpl.GenerateText"), zap.Any("result", map[string]any{

				// Aplica opções
				"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.GenerateText"), zap.Any("params", __logParams))

	messages := []model.Message{
		{
			Role:    model.RoleUser,
			Content: prompt,
		},
	}

	request := &model.LLMRequest{
		Messages: messages,
	}

	if maxTokens, ok := options["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	}
	if temperature, ok := options["temperature"].(float64); ok {
		request.Temperature = temperature
	}
	if topP, ok := options["top_p"].(float64); ok {
		request.TopP = topP
	}
	if responseFormat, ok := options["response_format"].(string); ok && responseFormat != "" {
		request.ResponseFormat = responseFormat
	}

	// Obtém o provedor ativo
	provider, err := s.getActiveProvider()
	if err != nil {
		zap.L().Error("function.error",

			// Executa a requisição
			zap.String("func", "*LLMServiceImpl.GenerateText"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	response, err := provider.Chat(ctx, request)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*LLMServiceImpl.GenerateText"), zap.Error(err), zap.Any("params", __logParams))

		// BuildPrompt constrói um prompt usando um template
		result0 = nil
		result1 = fmt.Errorf("erro na geração de texto: %w", err)
		return
	}
	result0 = dto.FromLLMResponse(response)
	result1 = nil
	return
}

func (s *LLMServiceImpl) BuildPrompt(ctx context.Context, templateID string, variables map[string]string) (result0 string, result1 error) {
	__logParams :=
		// TODO: Implementar busca de template no repositório
		// Por enquanto, retorna erro
		map[string]any{"s": s, "ctx": ctx, "templateID": templateID, "variables": variables}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String(

			// GetAvailableProviders retorna os provedores disponíveis
			"func", "*LLMServiceImpl.BuildPrompt"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().

		// SetProvider define o provedor ativo
		Info("function.entry", zap.String("func", "*LLMServiceImpl.BuildPrompt"), zap.Any("params", __logParams))
	result0 = ""
	result1 = fmt.Errorf("busca de templates não implementada")
	return
}

func (s *LLMServiceImpl) GetAvailableProviders() (result0 []string) {
	__logParams := map[string]any{"s": s}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.GetAvailableProviders"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.GetAvailableProviders"), zap.Any("params",

		// GetCurrentProvider retorna o provedor atual
		__logParams))
	result0 = s.factory.GetSupportedProviders()
	return
}

func (s *LLMServiceImpl) SetProvider(providerName string) (result0 error) {
	__logParams := map[string]any{"s": s, "providerName": providerName}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.SetProvider"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.SetProvider"),

		// getActiveProvider obtém o provedor ativo
		zap.Any("params", __logParams))
	if _, exists := s.configs[providerName]; !exists {
		result0 = fmt.Errorf("provedor '%s' não configurado", providerName)
		return
	}

	s.activeProvider = providerName
	result0 = nil
	return
}

func (s *LLMServiceImpl) GetCurrentProvider() (result0 string) {
	__logParams := map[string]any{"s": s}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.GetCurrentProvider"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.GetCurrentProvider"), zap.Any("params", __logParams))
	result0 = s.activeProvider
	return
}

func (s *LLMServiceImpl) getActiveProvider() (result0 domain.LLMProvider, result1 error) {
	__logParams := map[string]any{"s": s}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.getActiveProvider"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func",

		// CreateChatRequest cria uma requisição de chat com system e user prompts
		"*LLMServiceImpl.getActiveProvider"), zap.Any("params", __logParams))
	if s.activeProvider == "" {
		result0 = nil
		result1 = fmt.Errorf("nenhum provedor ativo configurado")
		return
	}

	config, exists := s.configs[s.activeProvider]
	if !exists {
		result0 = nil
		result1 = fmt.Errorf("configuração não encontrada para provedor '%s'", s.activeProvider)
		return
	}
	result0, result1 = s.factory.CreateProvider(config)
	return
}

func (s *LLMServiceImpl) CreateChatRequest(systemPrompt, userPrompt string, options map[string]interface{}) (result0 *dto.LLMRequestDTO) {
	__logParams := map[string]any{"s": s, "systemPrompt": systemPrompt, "userPrompt": userPrompt, "options": options}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.CreateChatRequest"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.CreateChatRequest"), zap.Any("params",

		// Aplica opções
		__logParams))
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

	if maxTokens, ok := options["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	}
	if temperature, ok := options["temperature"].(float64); ok {
		request.Temperature = temperature
	}
	if topP, ok := options["top_p"].(float64); ok {
		request.TopP = topP
	}
	result0 = request
	return
}

// EstimateTokens estima tokens usando o provedor ativo
func (s *LLMServiceImpl) EstimateTokens(text string) (result0 int, result1 error) {
	__logParams := map[string]any{"s": s, "text": text}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.EstimateTokens"),

			// GenerateID gera um ID único para requisições
			zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))

		// GetProviderInfo retorna informações sobre o provedor ativo
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.EstimateTokens"), zap.Any("params", __logParams))
	provider, err := s.getActiveProvider()
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*LLMServiceImpl.EstimateTokens"), zap.Error(err), zap.Any("params", __logParams))
		result0 = 0
		result1 = err
		return
	}
	result0 = provider.EstimateTokens(text)
	result1 = nil
	return
}

func (s *LLMServiceImpl) GenerateID() (result0 string) {
	__logParams := map[string]any{"s": s}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.GenerateID"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.GenerateID"), zap.Any("params", __logParams))
	result0 = uuid.New().String()
	return
}

func (s *LLMServiceImpl) GetProviderInfo() (result0 map[string]interface{}, result1 error) {
	__logParams := map[string]any{"s": s}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.GetProviderInfo"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.GetProviderInfo"), zap.Any("params", __logParams))
	if s.activeProvider == "" {
		result0 = nil
		result1 = fmt.Errorf("nenhum provedor ativo")
		return
	}

	config, exists := s.configs[s.activeProvider]
	if !exists {
		result0 = nil
		result1 = fmt.Errorf("configuração não encontrada")
		return
	}

	provider, err := s.getActiveProvider()
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*LLMServiceImpl.GetProviderInfo"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = map[string]interface{}{
		"name":     provider.GetProviderName(),
		"model":    provider.GetModel(),
		"provider": s.activeProvider,
		"timeout":  config.Timeout,
	}
	result1 = nil
	return
}

// initializeDefaultProviders configura automaticamente provedores com base nas variáveis de ambiente
func (s *LLMServiceImpl) initializeDefaultProviders() {
	__logParams :=
		// Verifica se Gemini está configurado (configurado primeiro para ser o padrão)
		map[string]any{"s": s}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*LLMServiceImpl.initializeDefaultProviders"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*LLMServiceImpl.initializeDefaultProviders"), zap.Any("params", __logParams))

	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		geminiConfig := &model.LLMConfig{
			Provider:    model.ProviderGemini,
			APIKey:      apiKey,
			Model:       "gemini-2.0-flash",
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

	// Aqui você pode adicionar configurações automáticas para outros provedores
	// como Anthropic Claude, etc.
}
