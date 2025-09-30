package model

import "time"

// Message representa uma mensagem no contexto de chat
type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// MessageRole define os tipos de papéis possíveis em uma conversa
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// LLMRequest representa uma solicitação para o LLM
type LLMRequest struct {
	Messages         []Message         `json:"messages"`
	MaxTokens        int               `json:"max_tokens,omitempty"`
	Temperature      float64           `json:"temperature,omitempty"`
	TopP             float64           `json:"top_p,omitempty"`
	FrequencyPenalty float64           `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64           `json:"presence_penalty,omitempty"`
	Stop             []string          `json:"stop,omitempty"`
	Model            string            `json:"model,omitempty"`
	Stream           bool              `json:"stream,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	ResponseFormat   string            `json:"response_format,omitempty"`
}

// LLMResponse representa a resposta do LLM
type LLMResponse struct {
	ID                string            `json:"id"`
	Object            string            `json:"object"`
	Created           int64             `json:"created"`
	Model             string            `json:"model"`
	Choices           []Choice          `json:"choices"`
	Usage             Usage             `json:"usage"`
	SystemFingerprint string            `json:"system_fingerprint,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// Choice representa uma opção de resposta do LLM
type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Logprobs     interface{} `json:"logprobs,omitempty"`
}

// Usage representa informações de uso de tokens
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// LLMProvider representa diferentes provedores de LLM
type LLMProvider string

const (
	ProviderOpenAI    LLMProvider = "openai"
	ProviderAnthropic LLMProvider = "anthropic"
	ProviderOllama    LLMProvider = "ollama"
	ProviderGemini    LLMProvider = "gemini"
)

// LLMConfig representa a configuração para um provedor de LLM
type LLMConfig struct {
	Provider       LLMProvider       `json:"provider"`
	APIKey         string            `json:"api_key"`
	BaseURL        string            `json:"base_url,omitempty"`
	Model          string            `json:"model"`
	MaxTokens      int               `json:"max_tokens,omitempty"`
	Temperature    float64           `json:"temperature,omitempty"`
	Timeout        time.Duration     `json:"timeout,omitempty"`
	RetryAttempts  int               `json:"retry_attempts,omitempty"`
	DefaultHeaders map[string]string `json:"default_headers,omitempty"`
}

// PromptTemplate representa um template de prompt reutilizável
type PromptTemplate struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	SystemPrompt string            `json:"system_prompt"`
	UserPrompt   string            `json:"user_prompt"`
	Variables    map[string]string `json:"variables"`
	Metadata     map[string]string `json:"metadata"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// LLMSession representa uma sessão de conversa com contexto
type LLMSession struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Messages  []Message         `json:"messages"`
	Context   map[string]string `json:"context"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	ExpiresAt time.Time         `json:"expires_at"`
}
