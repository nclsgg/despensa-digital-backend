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
)

// OpenAIProvider implementa o provedor OpenAI
type OpenAIProvider struct {
	config     *model.LLMConfig
	httpClient *http.Client
}

// NewOpenAIProvider cria uma nova instância do provedor OpenAI
func NewOpenAIProvider(config *model.LLMConfig) *OpenAIProvider {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &OpenAIProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Chat realiza uma conversa com o OpenAI
func (p *OpenAIProvider) Chat(ctx context.Context, request *model.LLMRequest) (*model.LLMResponse, error) {
	// Prepara a requisição para a API da OpenAI
	openAIRequest := map[string]interface{}{
		"model":    p.getModel(request.Model),
		"messages": p.convertMessages(request.Messages),
	}

	// Adiciona parâmetros opcionais
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

	// Serializa a requisição
	jsonData, err := json.Marshal(openAIRequest)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Cria a requisição HTTP
	baseURL := p.config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Adiciona headers
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
		return nil, fmt.Errorf("erro na requisição HTTP após %d tentativas: %w", maxRetries+1, lastErr)
	}
	defer response.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou erro %d: %s", response.StatusCode, string(body))
	}

	// Parseia a resposta
	var openAIResponse map[string]interface{}
	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	// Converte para o modelo padrão
	return p.convertResponse(openAIResponse)
}

// GetModel retorna o modelo atual
func (p *OpenAIProvider) GetModel() string {
	return p.config.Model
}

// GetProviderName retorna o nome do provedor
func (p *OpenAIProvider) GetProviderName() string {
	return string(model.ProviderOpenAI)
}

// ValidateConfig valida a configuração
func (p *OpenAIProvider) ValidateConfig() error {
	if p.config.APIKey == "" {
		return fmt.Errorf("API key é obrigatória para OpenAI")
	}
	if p.config.Model == "" {
		return fmt.Errorf("modelo é obrigatório para OpenAI")
	}
	return nil
}

// EstimateTokens estima o número de tokens
func (p *OpenAIProvider) EstimateTokens(text string) int {
	// Estimativa aproximada: ~4 caracteres por token em inglês
	// Para português, pode ser ligeiramente diferente
	return len(strings.Fields(text)) + len(text)/4
}

// getModel retorna o modelo a ser usado
func (p *OpenAIProvider) getModel(requestModel string) string {
	if requestModel != "" {
		return requestModel
	}
	if p.config.Model != "" {
		return p.config.Model
	}
	return "gpt-3.5-turbo" // modelo padrão
}

// convertMessages converte mensagens para o formato da OpenAI
func (p *OpenAIProvider) convertMessages(messages []model.Message) []map[string]string {
	converted := make([]map[string]string, len(messages))
	for i, msg := range messages {
		converted[i] = map[string]string{
			"role":    string(msg.Role),
			"content": msg.Content,
		}
	}
	return converted
}

// convertResponse converte a resposta da OpenAI para o modelo padrão
func (p *OpenAIProvider) convertResponse(openAIResponse map[string]interface{}) (*model.LLMResponse, error) {
	response := &model.LLMResponse{
		ID:      getString(openAIResponse, "id"),
		Object:  getString(openAIResponse, "object"),
		Created: getInt64(openAIResponse, "created"),
		Model:   getString(openAIResponse, "model"),
	}

	// Converte choices
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

	return response, nil
}

// Funções auxiliares para conversão de tipos
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	return 0
}

func getInt64(m map[string]interface{}, key string) int64 {
	if val, ok := m[key].(float64); ok {
		return int64(val)
	}
	return 0
}
