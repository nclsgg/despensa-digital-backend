package dto

import "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"

// ChatRequestDTO representa uma solicitação de chat simples
type ChatRequestDTO struct {
	Message  string `json:"message" validate:"required,min=1,max=2000"`
	Provider string `json:"provider,omitempty" validate:"omitempty,oneof=openai gemini anthropic ollama"`
	Context  string `json:"context,omitempty"`
}

// ChatResponseDTO representa a resposta de um chat
type ChatResponseDTO struct {
	Response string   `json:"response"`
	Provider string   `json:"provider"`
	Model    string   `json:"model"`
	Usage    UsageDTO `json:"usage"`
}

// LLMRequestDTO representa uma solicitação de LLM via API
type LLMRequestDTO struct {
	Messages         []MessageDTO      `json:"messages" validate:"required,min=1"`
	Provider         string            `json:"provider,omitempty" validate:"omitempty,oneof=openai gemini anthropic ollama"`
	MaxTokens        int               `json:"max_tokens,omitempty" validate:"min=1,max=4096"`
	Temperature      float64           `json:"temperature,omitempty" validate:"min=0,max=2"`
	TopP             float64           `json:"top_p,omitempty" validate:"min=0,max=1"`
	FrequencyPenalty float64           `json:"frequency_penalty,omitempty" validate:"min=-2,max=2"`
	PresencePenalty  float64           `json:"presence_penalty,omitempty" validate:"min=-2,max=2"`
	Stop             []string          `json:"stop,omitempty"`
	Stream           bool              `json:"stream,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// MessageDTO representa uma mensagem individual
type MessageDTO struct {
	Role    string `json:"role" validate:"required,oneof=system user assistant"`
	Content string `json:"content" validate:"required,min=1"`
}

// LLMResponseDTO representa a resposta padronizada do LLM
type LLMResponseDTO struct {
	ID       string            `json:"id"`
	Response string            `json:"response"`
	Model    string            `json:"model"`
	Usage    UsageDTO          `json:"usage"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UsageDTO representa informações de uso de tokens
type UsageDTO struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// PromptBuilderDTO representa dados para construção de prompts
type PromptBuilderDTO struct {
	Template  string            `json:"template" validate:"required"`
	Variables map[string]string `json:"variables" validate:"required"`
	Context   map[string]string `json:"context,omitempty"`
}

// RecipeRequestDTO representa uma solicitação de receita
type RecipeRequestDTO struct {
	PantryID            string   `json:"pantry_id" validate:"required,uuid"`
	Provider            string   `json:"provider,omitempty" validate:"omitempty,oneof=openai gemini anthropic ollama"`
	CookingTime         int      `json:"cooking_time,omitempty" validate:"min=5,max=480"` // 5 minutos a 8 horas
	MealType            string   `json:"meal_type,omitempty" validate:"oneof=breakfast lunch dinner snack dessert"`
	Difficulty          string   `json:"difficulty,omitempty" validate:"oneof=easy medium hard"`
	Cuisine             string   `json:"cuisine,omitempty" validate:"max=50"`
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
	ServingSize         int      `json:"serving_size,omitempty" validate:"min=1,max=20"`
	Purpose             string   `json:"purpose,omitempty" validate:"max=200"`
	AdditionalNotes     string   `json:"additional_notes,omitempty" validate:"max=500"`
}

// SetDefaults preenche campos opcionais com valores padrão se não enviados
func (dto *RecipeRequestDTO) SetDefaults() {
	if dto.DietaryRestrictions == nil {
		dto.DietaryRestrictions = []string{}
	}
	if dto.Purpose == "" {
		dto.Purpose = ""
	}
	if dto.AdditionalNotes == "" {
		dto.AdditionalNotes = ""
	}
}

// RecipeResponseDTO representa a resposta de uma receita gerada
type RecipeResponseDTO struct {
	ID                  string                 `json:"id"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	Ingredients         []RecipeIngredientDTO  `json:"ingredients"`
	Instructions        []RecipeInstructionDTO `json:"instructions"`
	CookingTime         int                    `json:"cooking_time"`
	PreparationTime     int                    `json:"preparation_time"`
	TotalTime           int                    `json:"total_time"`
	ServingSize         int                    `json:"serving_size"`
	Difficulty          string                 `json:"difficulty"`
	MealType            string                 `json:"meal_type"`
	Cuisine             string                 `json:"cuisine"`
	DietaryRestrictions []string               `json:"dietary_restrictions"`
	NutritionInfo       RecipeNutritionDTO     `json:"nutrition_info,omitempty"`
	Tips                []string               `json:"tips,omitempty"`
	SourceURL           string                 `json:"source_url,omitempty"`
	GeneratedAt         string                 `json:"generated_at"`
}

// RecipeIngredientDTO representa um ingrediente na receita
type RecipeIngredientDTO struct {
	Name        string   `json:"name"`
	Amount      *float64 `json:"amount"`
	Unit        string   `json:"unit"`
	Available   bool     `json:"available"`
	Alternative string   `json:"alternative,omitempty"`
}

// RecipeInstructionDTO representa uma instrução da receita
type RecipeInstructionDTO struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
	Time        int    `json:"time,omitempty"`
	Temperature string `json:"temperature,omitempty"`
}

// RecipeNutritionDTO representa informações nutricionais
type RecipeNutritionDTO struct {
	Calories      int     `json:"calories,omitempty"`
	Protein       float64 `json:"protein,omitempty"`
	Carbohydrates float64 `json:"carbohydrates,omitempty"`
	Fat           float64 `json:"fat,omitempty"`
	Fiber         float64 `json:"fiber,omitempty"`
	Sugar         float64 `json:"sugar,omitempty"`
	Sodium        float64 `json:"sodium,omitempty"`
}

// ToLLMRequest converte LLMRequestDTO para model.LLMRequest
func (dto *LLMRequestDTO) ToLLMRequest() *model.LLMRequest {
	messages := make([]model.Message, len(dto.Messages))
	for i, msg := range dto.Messages {
		messages[i] = model.Message{
			Role:    model.MessageRole(msg.Role),
			Content: msg.Content,
		}
	}

	return &model.LLMRequest{
		Messages:         messages,
		MaxTokens:        dto.MaxTokens,
		Temperature:      dto.Temperature,
		TopP:             dto.TopP,
		FrequencyPenalty: dto.FrequencyPenalty,
		PresencePenalty:  dto.PresencePenalty,
		Stop:             dto.Stop,
		Stream:           dto.Stream,
		Metadata:         dto.Metadata,
	}
}

// FromLLMResponse converte model.LLMResponse para LLMResponseDTO
func FromLLMResponse(response *model.LLMResponse) *LLMResponseDTO {
	var responseText string
	if len(response.Choices) > 0 {
		responseText = response.Choices[0].Message.Content
	}

	return &LLMResponseDTO{
		ID:       response.ID,
		Response: responseText,
		Model:    response.Model,
		Usage: UsageDTO{
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		},
		Metadata: response.Metadata,
	}
}
