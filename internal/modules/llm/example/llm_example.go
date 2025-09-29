package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"
	"go.uber.org/zap"
)

// Exemplo de uso do módulo LLM
func main() {
	__logParams := map[string]any{}
	__logStart := time.Now()
	defer func(

	// 1. Criar serviço LLM
	) {
		zap.L().Info("function.exit", zap.String("func", "main"), zap.Any("result",

			// 2. Configurar provedor Gemini como padrão
			nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "main"), zap.Any("params", __logParams))
	fmt.Println("=== Exemplo de uso do módulo LLM ===")

	llmService := service.NewLLMService()

	geminiConfig := &model.LLMConfig{
		Provider:      model.ProviderGemini,
		APIKey:        "AIzaSyC7eWrQc4jNKoFRkxWN2bD3Zq1GlHo8i4M", // Em produção, usar variáveis de ambiente
		Model:         "gemini-1.5-flash",
		MaxTokens:     2000,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	// 3. Adicionar configuração do provedor
	if err := llmService.AddProviderConfig("gemini", geminiConfig); err != nil {
		zap.L().Error("function.error", zap.String("func", "main"), zap.Error(err), zap.Any("params", __logParams))
		log.Printf("Erro ao configurar provedor Gemini: %v", err)
		return
	}

	fmt.Printf("✅ Provedor Gemini configurado com sucesso\n")
	fmt.Printf("📋 Provedores disponíveis: %v\n", llmService.GetAvailableProviders())
	fmt.Printf("🎯 Provedor ativo: %s\n", llmService.GetCurrentProvider())

	// 4. Demonstrar construção de prompts
	promptBuilder := service.NewPromptBuilder()

	systemTemplate := `Você é um chef especialista em {{cuisine}}. 
Ajude a criar receitas com ingredientes disponíveis: {{ingredients}}.
Tempo disponível: {{time}} minutos.`

	variables := map[string]string{
		"cuisine":     "culinária brasileira",
		"ingredients": "arroz, feijão, carne, tomate, cebola",
		"time":        "45",
	}

	systemPrompt, err := promptBuilder.BuildSystemPrompt(systemTemplate, variables)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "main"), zap.Error(err), zap.Any("params", __logParams))
		log.Printf("Erro ao construir prompt: %v", err)
		return
	}

	fmt.Printf("\n🤖 System Prompt construído:\n%s\n", systemPrompt)

	// 5. Demonstrar templates de receita
	recipeTemplates := service.GetRecipePromptTemplates()
	fmt.Printf("\n📝 Template de sistema para receitas disponível (tamanho: %d caracteres)\n",
		len(recipeTemplates.SystemPrompt))
	fmt.Printf("📝 Template de usuário para receitas disponível (tamanho: %d caracteres)\n",
		len(recipeTemplates.UserPrompt))

	// 6. Estimar tokens (estimativa simples)
	tokens := len(systemPrompt) / 4 // ~4 caracteres por token
	fmt.Printf("📊 Tokens estimados para o prompt: %d\n", tokens)

	// 7. Mostrar informações do provedor
	providerInfo, err := llmService.GetProviderInfo()
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "main"), zap.Error(err), zap.Any("params", __logParams))
		log.Printf("Erro ao obter informações do provedor: %v", err)
	} else {
		fmt.Printf("\n🔧 Informações do provedor:\n")
		for key, value := range providerInfo {
			fmt.Printf("  - %s: %v\n", key, value)
		}
	}

	fmt.Println("\n✨ Módulo LLM configurado e pronto para uso!")
	fmt.Println("\n📚 Próximos passos:")
	fmt.Println("  1. Integrar com sistema de itens da despensa")
	fmt.Println("  2. Configurar rotas e handlers REST")
	fmt.Println("  3. Adicionar suporte a outros provedores (Anthropic, Ollama, etc.)")
	fmt.Println("  4. Implementar sistema de templates persistentes")
	fmt.Println("  5. Adicionar sistema de métricas e monitoramento")
}
