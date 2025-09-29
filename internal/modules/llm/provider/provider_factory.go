package provider

import (
	"fmt"
	"time"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	"go.uber.org/zap"
)

// ProviderFactory é responsável por criar instâncias de provedores LLM
type ProviderFactory struct {
	providers map[model.LLMProvider]func(*model.LLMConfig) domain.LLMProvider
}

// NewProviderFactory cria uma nova instância do factory
func NewProviderFactory() (result0 *ProviderFactory) {
	__logParams := map[string]any{}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewProviderFactory"),

			// Registra provedores disponíveis
			zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewProviderFactory"), zap.Any("params", __logParams))
	factory := &ProviderFactory{
		providers: make(map[model.LLMProvider]func(*model.LLMConfig) domain.LLMProvider),
	}

	factory.RegisterProvider(model.ProviderOpenAI, func(config *model.LLMConfig) domain.LLMProvider {
		return NewOpenAIProvider(config)
	})

	factory.RegisterProvider(model.ProviderGemini, func(config *model.LLMConfig) domain.LLMProvider {
		return NewGeminiProvider(config)
	})
	result0 =

		// TODO: Adicionar outros provedores aqui
		// factory.RegisterProvider(model.ProviderAnthropic, func(config *model.LLMConfig) domain.LLMProvider {
		//     return NewAnthropicProvider(config)
		// })

		factory
	return
}

// RegisterProvider registra um novo provedor
func (f *ProviderFactory) RegisterProvider(providerType model.LLMProvider, constructor func(*model.LLMConfig) domain.LLMProvider) {
	__logParams := map[string]any{"f": f, "providerType": providerType,

		// CreateProvider cria um provedor baseado na configuração
		"constructor": constructor}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProviderFactory.RegisterProvider"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProviderFactory.RegisterProvider"), zap.Any("params", __logParams))
	f.providers[providerType] = constructor
}

func (f *ProviderFactory) CreateProvider(config *model.LLMConfig) (result0 domain.LLMProvider, result1 error) {
	__logParams := map[string]any{"f": f, "config": config}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProviderFactory.CreateProvider"), zap.Any("result", map[string]any{

			// Valida a configuração
			"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProviderFactory.CreateProvider"), zap.Any("params", __logParams))
	constructor, exists := f.providers[config.Provider]
	if !exists {
		result0 = nil
		result1 = fmt.Errorf("provedor '%s' não suportado", config.Provider)
		return
	}

	provider := constructor(config)

	if err := provider.ValidateConfig(); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ProviderFactory.CreateProvider"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("configuração inválida para provedor '%s': %w", config.Provider, err)
		return
	}
	result0 = provider
	result1 = nil
	return
}

// GetSupportedProviders retorna a lista de provedores suportados
func (f *ProviderFactory) GetSupportedProviders() (result0 []string) {
	__logParams := map[string]any{"f": f}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProviderFactory.GetSupportedProviders"), zap.Any("result",

			// IsProviderSupported verifica se um provedor é suportado
			result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProviderFactory.GetSupportedProviders"), zap.Any("params", __logParams))
	providers := make([]string, 0, len(f.providers))
	for provider := range f.providers {
		providers = append(providers, string(provider))
	}
	result0 = providers
	return
}

func (f *ProviderFactory) IsProviderSupported(provider model.LLMProvider) (result0 bool) {
	__logParams := map[string]any{"f": f, "provider": provider}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProviderFactory.IsProviderSupported"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProviderFactory.IsProviderSupported"), zap.Any("params", __logParams))
	_, exists := f.providers[provider]
	result0 = exists
	return
}
