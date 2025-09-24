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

// GeminiProvider implementa o provedor Google Gemini
type GeminiProvider struct {
	config     *model.LLMConfig
	httpClient *http.Client
}

// NewGeminiProvider cria uma nova instância do provedor Gemini
func NewGeminiProvider(config *model.LLMConfig) *GeminiProvider {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &GeminiProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Chat realiza uma conversa com o Gemini
func (p *GeminiProvider) Chat(ctx context.Context, request *model.LLMRequest) (*model.LLMResponse, error) {
	// Prepara a requisição para a API do Gemini
	geminiRequest := map[string]interface{}{
		"contents": p.convertMessages(request.Messages),
	}

	// Adiciona configurações de geração
	generationConfig := make(map[string]interface{})

	if request.MaxTokens > 0 {
		generationConfig["maxOutputTokens"] = request.MaxTokens
	}
	if request.Temperature > 0 {
		generationConfig["temperature"] = request.Temperature
	}
	if request.TopP > 0 {
		generationConfig["topP"] = request.TopP
	}
	if len(request.Stop) > 0 {
		generationConfig["stopSequences"] = request.Stop
	}

	if len(generationConfig) > 0 {
		geminiRequest["generationConfig"] = generationConfig
	}

	// Serializa a requisição
	jsonData, err := json.Marshal(geminiRequest)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Cria a requisição HTTP
	baseURL := p.config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	model := p.getModel(request.Model)
	endpoint := fmt.Sprintf("%s/models/%s:generateContent", baseURL, model)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Adiciona headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", p.config.APIKey)

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

	retryDelay := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay):
			}
		}

		response, lastErr = p.httpClient.Do(req)
		if lastErr == nil && response.StatusCode < 500 {
			break
		}

		if response != nil {
			response.Body.Close()
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("erro na requisição HTTP após %d tentativas: %w", maxRetries, lastErr)
	}
	defer response.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if response.StatusCode != 200 {
		return nil, p.parseError(body, response.StatusCode)
	}

	// Parse da resposta do Gemini
	var geminiResponse GeminiChatResponse
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	return p.convertResponse(&geminiResponse), nil
}

// convertMessages converte mensagens para o formato do Gemini
func (p *GeminiProvider) convertMessages(messages []model.Message) []map[string]interface{} {
	var geminiMessages []map[string]interface{}

	for _, msg := range messages {
		role := "user"
		if msg.Role == model.RoleAssistant {
			role = "model"
		} else if msg.Role == model.RoleSystem {
			// Gemini não tem role system, então convertemos para user
			role = "user"
		}

		geminiMessages = append(geminiMessages, map[string]interface{}{
			"role": role,
			"parts": []map[string]interface{}{
				{
					"text": msg.Content,
				},
			},
		})
	}

	return geminiMessages
}

// convertResponse converte a resposta do Gemini para o formato padrão
func (p *GeminiProvider) convertResponse(geminiResponse *GeminiChatResponse) *model.LLMResponse {
	response := &model.LLMResponse{
		ID:      fmt.Sprintf("gemini-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   p.config.Model,
		Choices: []model.Choice{},
	}

	if len(geminiResponse.Candidates) > 0 {
		candidate := geminiResponse.Candidates[0]
		if len(candidate.Content.Parts) > 0 {
			content := candidate.Content.Parts[0].Text

			choice := model.Choice{
				Index: 0,
				Message: model.Message{
					Role:    model.RoleAssistant,
					Content: content,
				},
				FinishReason: p.convertFinishReason(candidate.FinishReason),
			}
			response.Choices = append(response.Choices, choice)
		}
	}

	// Adiciona informações de uso se disponíveis
	if geminiResponse.UsageMetadata.TotalTokenCount > 0 {
		response.Usage = model.Usage{
			PromptTokens:     geminiResponse.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResponse.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResponse.UsageMetadata.TotalTokenCount,
		}
	}

	return response
}

// convertFinishReason converte o motivo de finalização do Gemini
func (p *GeminiProvider) convertFinishReason(reason string) string {
	switch reason {
	case "STOP":
		return "stop"
	case "MAX_TOKENS":
		return "length"
	case "SAFETY":
		return "content_filter"
	case "RECITATION":
		return "content_filter"
	default:
		return "stop"
	}
}

// parseError faz o parse de erros da API do Gemini
func (p *GeminiProvider) parseError(body []byte, statusCode int) error {
	var errorResponse map[string]interface{}
	if err := json.Unmarshal(body, &errorResponse); err == nil {
		if errorMsg, ok := errorResponse["error"].(map[string]interface{}); ok {
			if message, ok := errorMsg["message"].(string); ok {
				return fmt.Errorf("erro da API Gemini (%d): %s", statusCode, message)
			}
		}
	}
	return fmt.Errorf("erro da API Gemini (%d): %s", statusCode, string(body))
}

// getModel retorna o modelo a ser usado, com fallback para o padrão
func (p *GeminiProvider) getModel(requestModel string) string {
	if requestModel != "" {
		return requestModel
	}
	if p.config.Model != "" {
		return p.config.Model
	}
	return "gemini-1.5-flash"
}

// GetProviderName retorna o nome do provedor
func (p *GeminiProvider) GetProviderName() string {
	return "gemini"
}

// GetModel retorna o modelo configurado
func (p *GeminiProvider) GetModel() string {
	if p.config.Model != "" {
		return p.config.Model
	}
	return "gemini-1.5-flash"
}

// EstimateTokens estima a quantidade de tokens
func (p *GeminiProvider) EstimateTokens(text string) int {
	// Estimativa aproximada para Gemini (similar ao GPT)
	words := len(strings.Fields(text))
	return int(float64(words) * 1.3)
}

// ValidateConfig valida a configuração do provedor
func (p *GeminiProvider) ValidateConfig() error {
	if p.config.APIKey == "" {
		return fmt.Errorf("API key é obrigatória para o provedor Gemini")
	}
	return nil
}

// IsHealthy verifica se o provedor está funcionando
func (p *GeminiProvider) IsHealthy(ctx context.Context) error {
	// Faz uma requisição simples para verificar a conectividade
	testRequest := &model.LLMRequest{
		Messages: []model.Message{
			{
				Role:    model.RoleUser,
				Content: "Hello",
			},
		},
		MaxTokens: 1,
	}

	_, err := p.Chat(ctx, testRequest)
	return err
}

// GeminiChatResponse representa a resposta da API do Gemini
type GeminiChatResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason  string `json:"finishReason"`
		Index         int    `json:"index"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	PromptFeedback struct {
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}
