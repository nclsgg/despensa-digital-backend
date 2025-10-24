package logger

import (
	"regexp"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Constantes para nomes de campos padronizados.
const (
	// Campos de identificação
	FieldRequestID     = "request_id"
	FieldTraceID       = "trace_id"
	FieldCorrelationID = "correlation_id"
	FieldUserID        = "user_id"
	FieldSessionID     = "session_id"
	FieldEmail         = "email"
	FieldUsername      = "username"
	FieldRole          = "role"

	// Campos de contexto
	FieldModule      = "module"
	FieldFunction    = "function"
	FieldOperation   = "operation"
	FieldEnvironment = "environment"

	// Campos HTTP
	FieldHTTPMethod       = "http_method"
	FieldHTTPPath         = "http_path"
	FieldHTTPStatus       = "http_status"
	FieldHTTPRequestSize  = "http_request_size"
	FieldHTTPResponseSize = "http_response_size"
	FieldHTTPUserAgent    = "http_user_agent"
	FieldHTTPRemoteAddr   = "http_remote_addr"
	FieldHTTPHost         = "http_host"
	FieldHTTPScheme       = "http_scheme"
	FieldHTTPQuery        = "http_query"

	// Campos de performance
	FieldDuration   = "duration_ms"
	FieldLatency    = "latency_ms"
	FieldStartTime  = "start_time"
	FieldEndTime    = "end_time"
	FieldMemoryUsed = "memory_used_bytes"

	// Campos de erro
	FieldError      = "error"
	FieldErrorType  = "error_type"
	FieldErrorCode  = "error_code"
	FieldStackTrace = "stack_trace"

	// Campos de negócio
	FieldEntityID   = "entity_id"
	FieldEntityType = "entity_type"
	FieldAction     = "action"
	FieldStatus     = "status"
	FieldReason     = "reason"

	// Campos de dados
	FieldCount      = "count"
	FieldTotal      = "total"
	FieldPage       = "page"
	FieldPageSize   = "page_size"
	FieldItemsCount = "items_count"
)

// HTTPFields representa os campos de uma requisição HTTP.
type HTTPFields struct {
	Method       string
	Path         string
	Status       int
	RequestSize  int64
	ResponseSize int64
	UserAgent    string
	RemoteAddr   string
	Host         string
	Scheme       string
	Query        string
	Duration     time.Duration
}

// ToFields converte HTTPFields para zap.Field slice.
func (h HTTPFields) ToFields() []zap.Field {
	fields := []zap.Field{
		zap.String(FieldHTTPMethod, h.Method),
		zap.String(FieldHTTPPath, h.Path),
		zap.Int(FieldHTTPStatus, h.Status),
	}

	if h.RequestSize > 0 {
		fields = append(fields, zap.Int64(FieldHTTPRequestSize, h.RequestSize))
	}

	if h.ResponseSize > 0 {
		fields = append(fields, zap.Int64(FieldHTTPResponseSize, h.ResponseSize))
	}

	if h.UserAgent != "" {
		fields = append(fields, zap.String(FieldHTTPUserAgent, h.UserAgent))
	}

	if h.RemoteAddr != "" {
		fields = append(fields, zap.String(FieldHTTPRemoteAddr, h.RemoteAddr))
	}

	if h.Host != "" {
		fields = append(fields, zap.String(FieldHTTPHost, h.Host))
	}

	if h.Scheme != "" {
		fields = append(fields, zap.String(FieldHTTPScheme, h.Scheme))
	}

	if h.Query != "" {
		// Sanitizar query string de dados sensíveis
		sanitizedQuery := SanitizeQueryString(h.Query)
		fields = append(fields, zap.String(FieldHTTPQuery, sanitizedQuery))
	}

	if h.Duration > 0 {
		fields = append(fields, zap.Int64(FieldDuration, h.Duration.Milliseconds()))
	}

	return fields
}

// ErrorFields representa os campos de um erro.
type ErrorFields struct {
	Error      error
	ErrorType  string
	ErrorCode  string
	StackTrace string
	Operation  string
}

// ToFields converte ErrorFields para zap.Field slice.
func (e ErrorFields) ToFields() []zap.Field {
	fields := make([]zap.Field, 0, 5)

	if e.Error != nil {
		fields = append(fields, zap.Error(e.Error))
	}

	if e.ErrorType != "" {
		fields = append(fields, zap.String(FieldErrorType, e.ErrorType))
	}

	if e.ErrorCode != "" {
		fields = append(fields, zap.String(FieldErrorCode, e.ErrorCode))
	}

	if e.StackTrace != "" {
		fields = append(fields, zap.String(FieldStackTrace, e.StackTrace))
	}

	if e.Operation != "" {
		fields = append(fields, zap.String(FieldOperation, e.Operation))
	}

	return fields
}

// WithHTTP retorna campos para logging de requisições HTTP.
func WithHTTP(fields HTTPFields) []zap.Field {
	return fields.ToFields()
}

// WithError retorna campos para logging de erros.
func WithError(fields ErrorFields) []zap.Field {
	return fields.ToFields()
}

// WithUser retorna campos para logging relacionados a usuário.
func WithUser(userID string) zap.Field {
	return zap.String(FieldUserID, userID)
}

// WithDuration retorna campo de duração em milissegundos.
func WithDuration(duration time.Duration) zap.Field {
	return zap.Int64(FieldDuration, duration.Milliseconds())
}

// WithOperation retorna campo de operação.
func WithOperation(operation string) zap.Field {
	return zap.String(FieldOperation, operation)
}

// WithEntity retorna campos para uma entidade de negócio.
func WithEntity(entityType, entityID string) []zap.Field {
	return []zap.Field{
		zap.String(FieldEntityType, entityType),
		zap.String(FieldEntityID, entityID),
	}
}

// WithPagination retorna campos para paginação.
func WithPagination(page, pageSize, total int) []zap.Field {
	return []zap.Field{
		zap.Int(FieldPage, page),
		zap.Int(FieldPageSize, pageSize),
		zap.Int(FieldTotal, total),
	}
}

// WithCount retorna campo de contagem.
func WithCount(count int) zap.Field {
	return zap.Int(FieldCount, count)
}

// Padrões de regex para sanitização de dados sensíveis.
var (
	// Padrões para query strings
	sensitiveQueryParams = regexp.MustCompile(`(?i)(password|token|secret|key|auth|api[-_]?key|access[-_]?token|refresh[-_]?token)=([^&]*)`)

	// Padrões para emails (apenas parcialmente ocultar)
	emailPattern = regexp.MustCompile(`([a-zA-Z0-9._%+-]+)@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)

	// Padrões para números de cartão de crédito
	creditCardPattern = regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`)

	// Padrões para CPF/CNPJ brasileiros
	cpfPattern  = regexp.MustCompile(`\b\d{3}\.?\d{3}\.?\d{3}-?\d{2}\b`)
	cnpjPattern = regexp.MustCompile(`\b\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}\b`)
)

// SanitizeQueryString remove valores de parâmetros sensíveis da query string.
func SanitizeQueryString(query string) string {
	return sensitiveQueryParams.ReplaceAllString(query, `$1=***`)
}

// SanitizeEmail oculta parte do email mantendo domínio visível.
func SanitizeEmail(email string) string {
	return emailPattern.ReplaceAllStringFunc(email, func(match string) string {
		parts := emailPattern.FindStringSubmatch(match)
		if len(parts) >= 3 {
			username := parts[1]
			domain := parts[2]

			// Mostrar apenas primeiros 2 caracteres do username
			if len(username) > 2 {
				username = username[:2] + "***"
			} else {
				username = "***"
			}

			return username + "@" + domain
		}
		return "***@***"
	})
}

// SanitizeCreditCard mascara número de cartão mostrando apenas últimos 4 dígitos.
func SanitizeCreditCard(cardNumber string) string {
	return creditCardPattern.ReplaceAllStringFunc(cardNumber, func(match string) string {
		// Remover espaços e hífens
		clean := regexp.MustCompile(`[\s-]`).ReplaceAllString(match, "")
		if len(clean) >= 4 {
			return "****-****-****-" + clean[len(clean)-4:]
		}
		return "****-****-****-****"
	})
}

// SanitizeCPF mascara CPF mostrando apenas últimos 2 dígitos.
func SanitizeCPF(cpf string) string {
	return cpfPattern.ReplaceAllString(cpf, "***.***.**-**")
}

// SanitizeCNPJ mascara CNPJ mostrando apenas últimos 2 dígitos.
func SanitizeCNPJ(cnpj string) string {
	return cnpjPattern.ReplaceAllString(cnpj, "**.***.***/****-**")
}

// SanitizeString aplica todas as sanitizações em uma string.
func SanitizeString(s string) string {
	s = SanitizeEmail(s)
	s = SanitizeCreditCard(s)
	s = SanitizeCPF(s)
	s = SanitizeCNPJ(s)
	return s
}

// SanitizedString retorna um zap.Field com string sanitizada.
func SanitizedString(key, value string) zap.Field {
	return zap.String(key, SanitizeString(value))
}

// SensitiveString retorna um zap.Field que oculta completamente o valor.
// Use para senhas, tokens, etc.
func SensitiveString(key string) zap.Field {
	return zap.String(key, "***REDACTED***")
}

// ObjectEncoder customizado para objetos complexos que precisam de sanitização.
type SanitizedObjectEncoder struct {
	fields map[string]interface{}
}

// NewSanitizedObject cria um encoder para objetos com sanitização automática.
func NewSanitizedObject(fields map[string]interface{}) zapcore.ObjectMarshalerFunc {
	return func(enc zapcore.ObjectEncoder) error {
		for key, value := range fields {
			switch v := value.(type) {
			case string:
				// Sanitizar strings automaticamente
				enc.AddString(key, SanitizeString(v))
			case int:
				enc.AddInt(key, v)
			case int64:
				enc.AddInt64(key, v)
			case float64:
				enc.AddFloat64(key, v)
			case bool:
				enc.AddBool(key, v)
			default:
				// Para outros tipos, usar reflexão zap
				enc.AddReflected(key, v)
			}
		}
		return nil
	}
}
