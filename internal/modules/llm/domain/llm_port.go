package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
	recipeDTO "github.com/nclsgg/despensa-digital/backend/internal/modules/recipe/dto"
)

// LLMProvider interface para abstrair diferentes provedores de LLM
type LLMProvider interface {
	// Chat realiza uma conversa com o LLM
	Chat(ctx context.Context, request *model.LLMRequest) (*model.LLMResponse, error)

	// GetModel retorna o modelo atual sendo usado
	GetModel() string

	// GetProviderName retorna o nome do provedor
	GetProviderName() string

	// ValidateConfig valida se a configuração está correta
	ValidateConfig() error

	// EstimateTokens estima o número de tokens para um texto
	EstimateTokens(text string) int
}

// PromptBuilder interface para construção de prompts
type PromptBuilder interface {
	// BuildSystemPrompt constrói o prompt do sistema
	BuildSystemPrompt(template string, variables map[string]string) (string, error)

	// BuildUserPrompt constrói o prompt do usuário
	BuildUserPrompt(template string, variables map[string]string) (string, error)

	// BuildMessages constrói uma lista de mensagens
	BuildMessages(systemPrompt, userPrompt string) []model.Message

	// AddContext adiciona contexto ao prompt
	AddContext(prompt string, context map[string]string) string

	// ValidateTemplate valida um template de prompt
	ValidateTemplate(template string, requiredVariables []string) error
}

// LLMService interface para serviços de LLM
type LLMService interface {
	// ProcessRequest processa uma requisição genérica de LLM
	ProcessRequest(ctx context.Context, request *dto.LLMRequestDTO) (*dto.LLMResponseDTO, error)

	// GenerateText gera texto baseado em um prompt
	GenerateText(ctx context.Context, prompt string, options map[string]interface{}) (*dto.LLMResponseDTO, error)

	// BuildPrompt constrói um prompt usando um template
	BuildPrompt(ctx context.Context, templateID string, variables map[string]string) (string, error)

	// GetAvailableProviders retorna os provedores disponíveis
	GetAvailableProviders() []string

	// SetProvider define o provedor ativo
	SetProvider(providerName string) error

	// GetCurrentProvider retorna o provedor atual
	GetCurrentProvider() string
}

// RecipeService interface específica para geração de receitas
type RecipeService interface {
	// GenerateRecipe gera uma receita baseada nos parâmetros
	GenerateRecipe(ctx context.Context, request *dto.RecipeRequestDTO, userID uuid.UUID) (*dto.RecipeResponseDTO, error)

	// GetAvailableIngredients obtém ingredientes disponíveis na despensa
	GetAvailableIngredients(ctx context.Context, pantryID uuid.UUID, userID uuid.UUID) ([]recipeDTO.AvailableIngredientDTO, error)

	// SearchRecipesByIngredients busca receitas por ingredientes
	SearchRecipesByIngredients(ctx context.Context, ingredients []string, filters map[string]string) ([]dto.RecipeResponseDTO, error)

	// ValidateRecipeRequest valida uma requisição de receita
	ValidateRecipeRequest(request *dto.RecipeRequestDTO) (uuid.UUID, error)

	// EnrichRecipeWithNutrition adiciona informações nutricionais
	EnrichRecipeWithNutrition(ctx context.Context, recipe *dto.RecipeResponseDTO) error
}

// PromptTemplateRepository interface para persistência de templates
type PromptTemplateRepository interface {
	// Create cria um novo template
	Create(ctx context.Context, template *model.PromptTemplate) error

	// GetByID obtém um template por ID
	GetByID(ctx context.Context, id string) (*model.PromptTemplate, error)

	// GetByName obtém um template por nome
	GetByName(ctx context.Context, name string) (*model.PromptTemplate, error)

	// List lista templates com paginação
	List(ctx context.Context, offset, limit int) ([]*model.PromptTemplate, error)

	// Update atualiza um template
	Update(ctx context.Context, template *model.PromptTemplate) error

	// Delete remove um template
	Delete(ctx context.Context, id string) error
}

// LLMSessionRepository interface para persistência de sessões
type LLMSessionRepository interface {
	// Create cria uma nova sessão
	Create(ctx context.Context, session *model.LLMSession) error

	// GetByID obtém uma sessão por ID
	GetByID(ctx context.Context, id string) (*model.LLMSession, error)

	// GetByUserID obtém sessões de um usuário
	GetByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.LLMSession, error)

	// Update atualiza uma sessão
	Update(ctx context.Context, session *model.LLMSession) error

	// Delete remove uma sessão
	Delete(ctx context.Context, id string) error

	// Cleanup remove sessões expiradas
	Cleanup(ctx context.Context) error
}
