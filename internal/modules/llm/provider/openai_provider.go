package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	"go.uber.org/zap"
)

// OpenAIProvider implementa o provedor OpenAI
type OpenAIProvider struct {
	config     *model.LLMConfig
	httpClient *http.Client
}

// NewOpenAIProvider cria uma nova instância do provedor OpenAI
func NewOpenAIProvider(config *model.LLMConfig) (result0 *OpenAIProvider) {
	__logParams := map[string]any{"config": config}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewOpenAIProvider"), zap.Any("result", result0), zap.Duration("duration", time.

			// Chat realiza uma conversa com o OpenAI
			Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewOpenAIProvider"), zap.Any("params", __logParams))
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	result0 = &OpenAIProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return
}

func (p *OpenAIProvider) Chat(ctx context.Context, request *model.LLMRequest) (result0 *model.LLMResponse, result1 error) {
	__logParams :=
		// Prepara a requisição para a API da OpenAI
		map[string]any{"p": p, "ctx": ctx, "request": request}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*OpenAIProvider.Chat"), zap.Any("result", map[string]any{

			// Adiciona parâmetros opcionais
			"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.Chat"), zap.Any("params", __logParams))

	openAIRequest := map[string]interface{}{
		"model":    p.getModel(request.Model),
		"messages": p.convertMessages(request.Messages),
	}

	if request.MaxTokens > 0 {
		openAIRequest["max_tokens"] = request.MaxTokens
	}
	if request.Temperature > 0 {
		openAIRequest["temperature"] = request.Temperature
	}
	if request.TopP > 0 {
		openAIRequest["top_p"] = request.TopP
	}
	if request.FrequencyPenalty != 0 {
		openAIRequest["frequency_penalty"] = request.FrequencyPenalty
	}
	if request.PresencePenalty != 0 {
		openAIRequest["presence_penalty"] = request.PresencePenalty
	}
	if len(request.Stop) > 0 {
		openAIRequest["stop"] = request.Stop
	}
	if request.Stream {
		openAIRequest["stream"] = request.Stream
	}
	if format := strings.TrimSpace(request.ResponseFormat); format != "" {
		payload := map[string]string{}
		switch strings.ToLower(format) {
		case "json", "json_object":
			payload["type"] = "json_object"
		default:
			payload["type"] = format
		}
		openAIRequest["response_format"] = payload
	}

	// Serializa a requisição
	jsonData, err := json.Marshal(openAIRequest)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*OpenAIProvider.Chat"), zap.Error(

			// Cria a requisição HTTP
			err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao serializar requisição: %w", err)
		return
	}

	baseURL := p.config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*OpenAIProvider.Chat"), zap.Error(

			// Adiciona headers
			err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao criar requisição HTTP: %w", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	// Adiciona headers personalizados
	for key, value := range p.config.DefaultHeaders {
		req.Header.Set(key, value)
	}

	// Executa a requisição com retry
	var response *http.Response
	var lastErr error

	maxRetries := p.config.RetryAttempts
	if maxRetries == 0 {
		maxRetries = 3
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		response, lastErr = p.httpClient.Do(req)
		if lastErr == nil && response.StatusCode < 500 {
			break
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	if lastErr != nil {
		zap.L().Error("function.error", zap.String("func", "*OpenAIProvider.Chat"), zap.Error(lastErr), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro na requisição HTTP após %d tentativas: %w", maxRetries+1, lastErr)
		return
	}
	defer response.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*OpenAIProvider.Chat"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao ler resposta: %w", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		result0 = nil
		result1 = fmt.Errorf("API retornou erro %d: %s", response.StatusCode, string(body))
		return
	}

	// Parseia a resposta
	var openAIResponse map[string]interface{}
	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		zap.L().Error("function.error", zap.String("func", "*OpenAIProvider.Chat"),

			// Converte para o modelo padrão
			zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao parsear resposta: %w", err)
		return
	}
	result0, result1 = p.convertResponse(openAIResponse)
	return
}

// GetModel retorna o modelo atual
func (p *OpenAIProvider) GetModel() (result0 string) {
	__logParams := map[string]any{"p":

	// GetProviderName retorna o nome do provedor
	p}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*OpenAIProvider.GetModel"), zap.Any("result", result0), zap.Duration(

			// ValidateConfig valida a configuração
			"duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.GetModel"), zap.Any("params", __logParams))
	result0 = p.config.Model
	return
}

func (p *OpenAIProvider) GetProviderName() (result0 string) {
	__logParams := map[string]any{"p": p}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*OpenAIProvider.GetProviderName"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.GetProviderName"), zap.Any("params", __logParams))
	result0 = string(model.ProviderOpenAI)
	return
}

func (p *OpenAIProvider) ValidateConfig() (result0 error) {
	__logParams := map[string]any{"p": p}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*OpenAIProvider.ValidateConfig"), zap.Any("result", result0), zap.Duration("duration", time.Since(

			// EstimateTokens estima o número de tokens
			__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.ValidateConfig"), zap.

		// Estimativa aproximada: ~4 caracteres por token em inglês
		// Para português, pode ser ligeiramente diferente
		Any("params", __logParams))
	if p.config.APIKey == "" {
		result0 = fmt.Errorf("API key é obrigatória para OpenAI")
		return
	}
	if p.config.Model == "" {
		result0 = fmt.Errorf("modelo é obrigatório para OpenAI")
		return
	}
	result0 = nil
	return
}

func (p *OpenAIProvider) EstimateTokens(text string) (result0 int) {
	__logParams := map[string]any{"p": p, "text": text}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*OpenAIProvider.EstimateTokens"), zap.Any("result",

			// getModel retorna o modelo a ser usado
			result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.EstimateTokens"), zap.Any("params", __logParams))
	result0 = len(strings.Fields(text)) + len(text)/4
	return
}

func (p *OpenAIProvider) getModel(requestModel string) (result0 string) {
	__logParams := map[string]any{"p": p, "requestModel": requestModel}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func",

			// modelo padrão
			"*OpenAIProvider.getModel"),

			// convertMessages converte mensagens para o formato da OpenAI
			zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.getModel"), zap.Any("params", __logParams))
	if requestModel != "" {
		result0 = requestModel
		return
	}
	if p.config.Model != "" {
		result0 = p.config.Model
		return
	}
	result0 = "gpt-3.5-turbo"
	return
}

func (p *OpenAIProvider) convertMessages(messages []model.Message) (result0 []map[string]string) {
	__logParams := map[string]any{"p": p, "messages": messages}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*OpenAIProvider.convertMessages"), zap.Any("result", result0), zap.Duration("duration", time.

			// convertResponse converte a resposta da OpenAI para o modelo padrão
			Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.convertMessages"), zap.Any("params", __logParams))
	converted := make([]map[string]string, len(messages))
	for i, msg := range messages {
		converted[i] = map[string]string{
			"role":    string(msg.Role),
			"content": msg.Content,
		}
	}
	result0 = converted
	return
}

func (p *OpenAIProvider) convertResponse(openAIResponse map[string]interface{}) (result0 *model.LLMResponse, result1 error) {
	__logParams := map[string]any{"p": p, "openAIResponse": openAIResponse}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*OpenAIProvider.convertResponse"), zap.Any("result", map[string]any{"result0": result0, "result1":

		// Converte choices
		result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*OpenAIProvider.convertResponse"), zap.Any("params", __logParams))
	response := &model.LLMResponse{
		ID:      getString(openAIResponse, "id"),
		Object:  getString(openAIResponse, "object"),
		Created: getInt64(openAIResponse, "created"),
		Model:   getString(openAIResponse, "model"),
	}

	if choicesData, ok := openAIResponse["choices"].([]interface{}); ok {
		response.Choices = make([]model.Choice, len(choicesData))
		for i, choiceData := range choicesData {
			if choice, ok := choiceData.(map[string]interface{}); ok {
				response.Choices[i] = model.Choice{
					Index:        getInt(choice, "index"),
					FinishReason: getString(choice, "finish_reason"),
				}

				// Converte message
				if msgData, ok := choice["message"].(map[string]interface{}); ok {
					response.Choices[i].Message = model.Message{
						Role:    model.MessageRole(getString(msgData, "role")),
						Content: getString(msgData, "content"),
					}
				}
			}
		}
	}

	// Converte usage
	if usageData, ok := openAIResponse["usage"].(map[string]interface{}); ok {
		response.Usage = model.Usage{
			PromptTokens:     getInt(usageData, "prompt_tokens"),
			CompletionTokens: getInt(usageData, "completion_tokens"),
			TotalTokens:      getInt(usageData, "total_tokens"),
		}
	}
	result0 = response
	result1 = nil
	return
}

// Funções auxiliares para conversão de tipos
func getString(m map[string]interface{}, key string) (result0 string) {
	__logParams := map[string]any{"m": m, "key": key}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "getString"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "getString"), zap.Any("params", __logParams))
	if val, ok := m[key].(string); ok {
		result0 = val
		return
	}
	result0 = ""
	return
}

func getInt(m map[string]interface{}, key string) (result0 int) {
	__logParams := map[string]any{"m": m, "key": key}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "getInt"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "getInt"), zap.Any("params", __logParams))
	if val, ok := m[key].(float64); ok {
		result0 = int(val)
		return
	}
	result0 = 0
	return
}

func getInt64(m map[string]interface{}, key string) (result0 int64) {
	__logParams := map[string]any{"m": m, "key": key}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "getInt64"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "getInt64"), zap.Any("params", __logParams))
	if val, ok := m[key].(float64); ok {
		result0 = int64(val)
		return
	}
	result0 = 0
	return
}
