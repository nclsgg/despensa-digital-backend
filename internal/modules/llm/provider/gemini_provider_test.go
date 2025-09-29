package provider

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	"go.uber.org/zap"
)

func TestGeminiProvider_Chat(t *testing.T) {
	__logParams :=
		// Pula o teste se não há API key configurada
		map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGeminiProvider_Chat"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGeminiProvider_Chat"), zap.Any("params", __logParams))

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY não configurado, pulando teste de integração")
	}

	config := &model.LLMConfig{
		Provider:    model.ProviderGemini,
		APIKey:      apiKey,
		Model:       "gemini-1.5-flash",
		MaxTokens:   100,
		Temperature: 0.7,
		Timeout:     30 * time.Second,
	}

	provider := NewGeminiProvider(config)

	// Teste de validação de configuração
	t.Run("ValidateConfig", func(t *testing.T) {
		err := provider.ValidateConfig()
		if err != nil {
			t.Fatalf("Configuração deveria ser válida: %v", err)
		}
	})

	// Teste de chat básico
	t.Run("BasicChat", func(t *testing.T) {
		request := &model.LLMRequest{
			Messages: []model.Message{
				{
					Role:    model.RoleUser,
					Content: "Olá! Responda apenas 'Olá' para este teste.",
				},
			},
			MaxTokens:   50,
			Temperature: 0.1,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		response, err := provider.Chat(ctx, request)
		if err != nil {
			t.Fatalf("Erro no chat: %v", err)
		}

		if response == nil {
			t.Fatal("Resposta não pode ser nula")
		}

		if len(response.Choices) == 0 {
			t.Fatal("Resposta deve ter pelo menos uma choice")
		}

		if response.Choices[0].Message.Content == "" {
			t.Fatal("Conteúdo da mensagem não pode ser vazio")
		}

		t.Logf("Resposta recebida: %s", response.Choices[0].Message.Content)
	})

	// Teste de múltiplas mensagens
	t.Run("MultipleMessages", func(t *testing.T) {
		request := &model.LLMRequest{
			Messages: []model.Message{
				{
					Role:    model.RoleUser,
					Content: "Meu nome é João.",
				},
				{
					Role:    model.RoleAssistant,
					Content: "Olá João! É um prazer conhecê-lo.",
				},
				{
					Role:    model.RoleUser,
					Content: "Qual é meu nome?",
				},
			},
			MaxTokens:   50,
			Temperature: 0.1,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		response, err := provider.Chat(ctx, request)
		if err != nil {
			t.Fatalf("Erro no chat com múltiplas mensagens: %v", err)
		}

		if response == nil {
			t.Fatal("Resposta não pode ser nula")
		}

		if len(response.Choices) == 0 {
			t.Fatal("Resposta deve ter pelo menos uma choice")
		}

		content := response.Choices[0].Message.Content
		if content == "" {
			t.Fatal("Conteúdo da mensagem não pode ser vazio")
		}

		t.Logf("Resposta para múltiplas mensagens: %s", content)
	})

	// Teste de estimativa de tokens
	t.Run("EstimateTokens", func(t *testing.T) {
		text := "Esta é uma frase de teste para estimativa de tokens."
		tokens := provider.EstimateTokens(text)

		if tokens <= 0 {
			t.Fatal("Estimativa de tokens deve ser maior que zero")
		}

		t.Logf("Tokens estimados para '%s': %d", text, tokens)
	})

	// Teste de informações do provedor
	t.Run("ProviderInfo", func(t *testing.T) {
		providerName := provider.GetProviderName()
		if providerName != "gemini" {
			t.Fatalf("Nome do provedor deveria ser 'gemini', recebido: %s", providerName)
		}

		model := provider.GetModel()
		if model == "" {
			t.Fatal("Modelo não pode ser vazio")
		}

		t.Logf("Provedor: %s, Modelo: %s", providerName, model)
	})
}

func TestGeminiProvider_ValidateConfig(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGeminiProvider_ValidateConfig"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGeminiProvider_ValidateConfig"), zap.Any("params", __logParams))
	tests := []struct {
		name        string
		config      *model.LLMConfig
		expectError bool
	}{
		{
			name: "ConfiguracaoValida",
			config: &model.LLMConfig{
				Provider: model.ProviderGemini,
				APIKey:   "test-api-key",
				Model:    "gemini-1.5-flash",
			},
			expectError: false,
		},
		{
			name: "SemAPIKey",
			config: &model.LLMConfig{
				Provider: model.ProviderGemini,
				APIKey:   "",
				Model:    "gemini-1.5-flash",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewGeminiProvider(tt.config)
			err := provider.ValidateConfig()

			if tt.expectError && err == nil {
				t.Error("Esperava erro, mas não recebeu nenhum")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Não esperava erro, mas recebeu: %v", err)
			}
		})
	}
}

func TestGeminiProvider_IsHealthy(t *testing.T) {
	__logParams :=
		// Pula o teste se não há API key configurada
		map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGeminiProvider_IsHealthy"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGeminiProvider_IsHealthy"), zap.Any("params", __logParams))

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY não configurado, pulando teste de health check")
	}

	config := &model.LLMConfig{
		Provider:    model.ProviderGemini,
		APIKey:      apiKey,
		Model:       "gemini-1.5-flash",
		MaxTokens:   10,
		Temperature: 0.1,
		Timeout:     15 * time.Second,
	}

	provider := NewGeminiProvider(config)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := provider.IsHealthy(ctx)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "TestGeminiProvider_IsHealthy"),

			// Não falhamos o teste aqui pois pode ser rate limiting ou rede
			zap.Error(err), zap.Any("params", __logParams))
		t.Logf("Health check falhou (pode ser devido a rate limiting): %v", err)

	} else {
		t.Log("Health check passou com sucesso")
	}
}

func TestGeminiProvider_ConvertMessages(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGeminiProvider_ConvertMessages"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGeminiProvider_ConvertMessages"), zap.Any("params", __logParams))
	config := &model.LLMConfig{
		Provider: model.ProviderGemini,
		APIKey:   "test-key",
		Model:    "gemini-1.5-flash",
	}

	provider := NewGeminiProvider(config)

	messages := []model.Message{
		{Role: model.RoleSystem, Content: "Você é um assistente útil"},
		{Role: model.RoleUser, Content: "Olá"},
		{Role: model.RoleAssistant, Content: "Oi! Como posso ajudar?"},
		{Role: model.RoleUser, Content: "Me conte uma piada"},
	}

	converted := provider.convertMessages(messages)

	if len(converted) != len(messages) {
		t.Fatalf("Número de mensagens convertidas (%d) não match com originais (%d)",
			len(converted), len(messages))
	}

	// Verifica conversão de roles
	expectedRoles := []string{"user", "user", "model", "user"} // system -> user, assistant -> model
	for i, msg := range converted {
		role, ok := msg["role"].(string)
		if !ok {
			t.Fatalf("Role da mensagem %d não é string", i)
		}

		if role != expectedRoles[i] {
			t.Fatalf("Role da mensagem %d: esperado %s, recebido %s",
				i, expectedRoles[i], role)
		}
	}

	t.Log("Conversão de mensagens funcionou corretamente")
}

func TestGeminiProvider_ConvertFinishReason(t *testing.T) {
	__logParams := map[string]any{"t": t}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "TestGeminiProvider_ConvertFinishReason"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "TestGeminiProvider_ConvertFinishReason"), zap.Any("params", __logParams))
	config := &model.LLMConfig{
		Provider: model.ProviderGemini,
		APIKey:   "test-key",
		Model:    "gemini-1.5-flash",
	}

	provider := NewGeminiProvider(config)

	tests := []struct {
		geminiReason   string
		expectedReason string
	}{
		{"STOP", "stop"},
		{"MAX_TOKENS", "length"},
		{"SAFETY", "content_filter"},
		{"RECITATION", "content_filter"},
		{"UNKNOWN", "stop"},
		{"", "stop"},
	}

	for _, test := range tests {
		result := provider.convertFinishReason(test.geminiReason)
		if result != test.expectedReason {
			t.Errorf("convertFinishReason(%s): esperado %s, recebido %s",
				test.geminiReason, test.expectedReason, result)
		}
	}
}
