package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/service"
)

// Exemplo de uso do módulo LLM
func main() {
	fmt.Println("=== Exemplo de uso do módulo LLM ===")

	// 1. Criar serviço LLM
	llmService := service.NewLLMService()

	// 2. Configurar provedor OpenAI
	openAIConfig := &model.LLMConfig{
		Provider:      model.ProviderOpenAI,
		APIKey:        "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", // Em produção, usar variáveis de ambiente
		Model:         "gpt-3.5-turbo",
		MaxTokens:     2000,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	// 3. Adicionar configuração do provedor
	if err := llmService.AddProviderConfig("openai", openAIConfig); err != nil {
		log.Printf("Erro ao configurar provedor OpenAI: %v", err)
		return
	}

	fmt.Printf("✅ Provedor OpenAI configurado com sucesso\n")
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
