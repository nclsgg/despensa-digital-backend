package provider

import (
	"fmt"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
)

// ProviderFactory é responsável por criar instâncias de provedores LLM
type ProviderFactory struct {
	providers map[model.LLMProvider]func(*model.LLMConfig) domain.LLMProvider
}

// NewProviderFactory cria uma nova instância do factory
func NewProviderFactory() *ProviderFactory {
	factory := &ProviderFactory{
		providers: make(map[model.LLMProvider]func(*model.LLMConfig) domain.LLMProvider),
	}

	// Registra provedores disponíveis
	factory.RegisterProvider(model.ProviderOpenAI, func(config *model.LLMConfig) domain.LLMProvider {
		return NewOpenAIProvider(config)
	})

	factory.RegisterProvider(model.ProviderGemini, func(config *model.LLMConfig) domain.LLMProvider {
		return NewGeminiProvider(config)
	})

	// TODO: Adicionar outros provedores aqui
	// factory.RegisterProvider(model.ProviderAnthropic, func(config *model.LLMConfig) domain.LLMProvider {
	//     return NewAnthropicProvider(config)
	// })

	return factory
}

// RegisterProvider registra um novo provedor
func (f *ProviderFactory) RegisterProvider(providerType model.LLMProvider, constructor func(*model.LLMConfig) domain.LLMProvider) {
	f.providers[providerType] = constructor
}

// CreateProvider cria um provedor baseado na configuração
func (f *ProviderFactory) CreateProvider(config *model.LLMConfig) (domain.LLMProvider, error) {
	constructor, exists := f.providers[config.Provider]
	if !exists {
		return nil, fmt.Errorf("provedor '%s' não suportado", config.Provider)
	}

	provider := constructor(config)

	// Valida a configuração
	if err := provider.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("configuração inválida para provedor '%s': %w", config.Provider, err)
	}

	return provider, nil
}

// GetSupportedProviders retorna a lista de provedores suportados
func (f *ProviderFactory) GetSupportedProviders() []string {
	providers := make([]string, 0, len(f.providers))
	for provider := range f.providers {
		providers = append(providers, string(provider))
	}
	return providers
}

// IsProviderSupported verifica se um provedor é suportado
func (f *ProviderFactory) IsProviderSupported(provider model.LLMProvider) bool {
	_, exists := f.providers[provider]
	return exists
}
