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

// GeminiProvider implementa o provedor Google Gemini
type GeminiProvider struct {
	config     *model.LLMConfig
	httpClient *http.Client
}

// NewGeminiProvider cria uma nova instância do provedor Gemini
func NewGeminiProvider(config *model.LLMConfig) (result0 *GeminiProvider) {
	__logParams := map[string]any{"config": config}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewGeminiProvider"), zap.Any("result", result0), zap.Duration("duration", time.

			// Chat realiza uma conversa com o Gemini
			Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewGeminiProvider"), zap.Any("params", __logParams))
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	result0 = &GeminiProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return
}

func (p *GeminiProvider) Chat(ctx context.Context, request *model.LLMRequest) (result0 *model.LLMResponse, result1 error) {
	__logParams :=
		// Prepara a requisição para a API do Gemini
		map[string]any{"p": p, "ctx": ctx, "request": request}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.Chat"),

			// Adiciona configurações de geração
			zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.Chat"), zap.Any("params", __logParams))

	geminiRequest := map[string]interface{}{
		"contents": p.convertMessages(request.Messages),
	}

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
		zap.L().Error("function.error", zap.String("func", "*GeminiProvider.Chat"), zap.Error(

			// Cria a requisição HTTP
			err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao serializar requisição: %w", err)
		return
	}

	baseURL := p.config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	model := p.getModel(request.Model)
	endpoint := fmt.Sprintf("%s/models/%s:generateContent", baseURL, model)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*GeminiProvider.Chat"), zap.Error(

			// Adiciona headers
			err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao criar requisição HTTP: %w", err)
		return
	}

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
				result0 = nil
				result1 = ctx.Err()
				return
			case <-time.After(retryDelay):
			}
		}

		response, lastErr = p.httpClient.Do(req)
		if lastErr == nil && response.StatusCode < 500 {
			break
		}

		if response != nil {
			zap.L().Error(
				"function.error",
				zap.String("func", "*GeminiProvider.Chat"),
				zap.Int("status_code", response.StatusCode),
				zap.String("status", response.Status),
				zap.Any("params", __logParams),
			)
			response.Body.Close()
		}
	}

	if lastErr != nil {
		zap.L().Error("function.error", zap.String("func", "*GeminiProvider.Chat"), zap.Error(lastErr), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro na requisição HTTP após %d tentativas: %w", maxRetries, lastErr)
		return
	}
	defer response.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*GeminiProvider.Chat"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao ler resposta: %w", err)
		return
	}

	if response.StatusCode != 200 {
		result0 = nil
		result1 = p.parseError(body, response.StatusCode)
		return
	}

	// Parse da resposta do Gemini
	var geminiResponse GeminiChatResponse
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		zap.L().Error("function.error", zap.String("func", "*GeminiProvider.Chat"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("erro ao fazer parse da resposta: %w", err)
		return
	}
	result0 = p.convertResponse(&geminiResponse)
	result1 = nil
	return
}

// convertMessages converte mensagens para o formato do Gemini
func (p *GeminiProvider) convertMessages(messages []model.Message) (result0 []map[string]interface{}) {
	__logParams := map[string]any{"p": p, "messages": messages}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.convertMessages"), zap.Any("result", result0), zap.Duration(

			// Gemini não tem role system, então convertemos para user
			"duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.convertMessages"), zap.Any("params", __logParams))
	var geminiMessages []map[string]interface{}

	for _, msg := range messages {
		role := "user"
		if msg.Role == model.RoleAssistant {
			role = "model"
		} else if msg.Role == model.RoleSystem {

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
	result0 = geminiMessages
	return
}

// convertResponse converte a resposta do Gemini para o formato padrão
func (p *GeminiProvider) convertResponse(geminiResponse *GeminiChatResponse) (result0 *model.LLMResponse) {
	__logParams := map[string]any{"p": p, "geminiResponse": geminiResponse}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.convertResponse"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.convertResponse"), zap.Any("params", __logParams))
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
	result0 = response
	return
}

// convertFinishReason converte o motivo de finalização do Gemini
func (p *GeminiProvider) convertFinishReason(reason string) (result0 string) {
	__logParams := map[string]any{"p": p, "reason": reason}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.convertFinishReason"), zap.Any("result", result0), zap.Duration("duration", time.

			// parseError faz o parse de erros da API do Gemini
			Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.convertFinishReason"), zap.Any("params", __logParams))
	switch reason {
	case "STOP":
		result0 = "stop"
		return
	case "MAX_TOKENS":
		result0 = "length"
		return
	case "SAFETY":
		result0 = "content_filter"
		return
	case "RECITATION":
		result0 = "content_filter"
		return
	default:
		result0 = "stop"
		return
	}
}

func (p *GeminiProvider) parseError(body []byte, statusCode int) (result0 error) {
	__logParams := map[string]any{"p": p, "body": body, "statusCode": statusCode}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.parseError"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.parseError"), zap.Any("params", __logParams))
	var errorResponse map[string]interface{}
	if err := json.Unmarshal(body, &errorResponse); err == nil {
		if errorMsg, ok := errorResponse["error"].(map[string]interface{}); ok {
			if message, ok := errorMsg["message"].(string); ok {
				result0 = fmt.Errorf("erro da API Gemini (%d): %s", statusCode, message)
				return
			}
		}
	}
	result0 = fmt.Errorf("erro da API Gemini (%d): %s", statusCode, string(body))
	return
}

// getModel retorna o modelo a ser usado, com fallback para o padrão
func (p *GeminiProvider) getModel(requestModel string) (result0 string) {
	__logParams := map[string]any{"p": p, "requestModel": requestModel}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.getModel"),

			// GetProviderName retorna o nome do provedor
			zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry",

		// GetModel retorna o modelo configurado
		zap.String("func", "*GeminiProvider.getModel"), zap.Any("params", __logParams))
	if requestModel != "" {
		result0 = requestModel
		return
	}
	if p.config.Model != "" {
		result0 = p.config.Model
		return
	}
	result0 = "gemini-2.0-flash"
	return
}

func (p *GeminiProvider) GetProviderName() (result0 string) {
	__logParams := map[string]any{"p": p}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.GetProviderName"), zap.Any("result", result0), zap.Duration("duration", time.

			// EstimateTokens estima a quantidade de tokens
			Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.GetProviderName"), zap.

		// Estimativa aproximada para Gemini (similar ao GPT)
		Any("params", __logParams))
	result0 = "gemini"
	return
}

func (p *GeminiProvider) GetModel() (result0 string) {
	__logParams := map[string]any{"p": p}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.GetModel"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.GetModel"), zap.Any("params", __logParams))
	if p.config.Model != "" {
		result0 = p.config.Model
		return
	}
	result0 = "gemini-2.0-flash"
	return
}

func (p *GeminiProvider) EstimateTokens(text string) (result0 int) {
	__logParams := map[string]any{"p": p, "text": text}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.EstimateTokens"),

			// ValidateConfig valida a configuração do provedor
			zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.EstimateTokens"), zap.Any("params", __logParams))

	words := len(strings.Fields(text))
	result0 = int(float64(words) * 1.3)
	return
}

func (p *GeminiProvider) ValidateConfig() (result0 error) {
	__logParams := map[string]any{"p": p}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.ValidateConfig"),

			// IsHealthy verifica se o provedor está funcionando
			zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry",

		// Faz uma requisição simples para verificar a conectividade
		zap.String("func", "*GeminiProvider.ValidateConfig"), zap.Any("params", __logParams))
	if p.config.APIKey == "" {
		result0 = fmt.Errorf("API key é obrigatória para o provedor Gemini")
		return
	}
	result0 = nil
	return
}

func (p *GeminiProvider) IsHealthy(ctx context.Context) (result0 error) {
	__logParams := map[string]any{"p": p, "ctx": ctx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*GeminiProvider.IsHealthy"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*GeminiProvider.IsHealthy"),

		// GeminiChatResponse representa a resposta da API do Gemini
		zap.Any("params", __logParams))

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
	result0 = err
	return
}

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
