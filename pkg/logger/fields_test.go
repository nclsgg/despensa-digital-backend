package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHTTPFields_ToFields(t *testing.T) {
	fields := HTTPFields{
		Method:       "GET",
		Path:         "/api/users",
		Status:       200,
		RequestSize:  1024,
		ResponseSize: 2048,
		UserAgent:    "Mozilla/5.0",
		RemoteAddr:   "192.168.1.1",
		Host:         "example.com",
		Scheme:       "https",
		Query:        "page=1&limit=10",
		Duration:     150 * time.Millisecond,
	}

	zapFields := fields.ToFields()

	assert.NotEmpty(t, zapFields)
	assert.GreaterOrEqual(t, len(zapFields), 3) // Pelo menos method, path, status
}

func TestErrorFields_ToFields(t *testing.T) {
	err := assert.AnError
	fields := ErrorFields{
		Error:      err,
		ErrorType:  "ValidationError",
		ErrorCode:  "ERR001",
		StackTrace: "stack trace here",
		Operation:  "CreateUser",
	}

	zapFields := fields.ToFields()

	assert.NotEmpty(t, zapFields)
	assert.Len(t, zapFields, 5)
}

func TestWithHTTP(t *testing.T) {
	fields := HTTPFields{
		Method: "POST",
		Path:   "/api/orders",
		Status: 201,
	}

	zapFields := WithHTTP(fields)
	assert.NotEmpty(t, zapFields)
}

func TestWithError(t *testing.T) {
	fields := ErrorFields{
		Error:     assert.AnError,
		ErrorType: "DatabaseError",
	}

	zapFields := WithError(fields)
	assert.NotEmpty(t, zapFields)
}

func TestWithUser(t *testing.T) {
	field := WithUser("user-123")
	assert.NotNil(t, field)
}

func TestWithDuration(t *testing.T) {
	duration := 250 * time.Millisecond
	field := WithDuration(duration)
	assert.NotNil(t, field)
}

func TestWithOperation(t *testing.T) {
	field := WithOperation("database_query")
	assert.NotNil(t, field)
}

func TestWithEntity(t *testing.T) {
	fields := WithEntity("order", "order-123")
	assert.Len(t, fields, 2)
}

func TestWithPagination(t *testing.T) {
	fields := WithPagination(1, 20, 100)
	assert.Len(t, fields, 3)
}

func TestWithCount(t *testing.T) {
	field := WithCount(42)
	assert.NotNil(t, field)
}

func TestSanitizeQueryString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "sanitize password",
			input: "username=john&password=secret123",
			want:  "username=john&password=***",
		},
		{
			name:  "sanitize token",
			input: "user=john&token=abc123xyz",
			want:  "user=john&token=***",
		},
		{
			name:  "sanitize api key",
			input: "service=api&api_key=sk_live_123456",
			want:  "service=api&api_key=***",
		},
		{
			name:  "sanitize access token",
			input: "user=john&access_token=eyJhbGc",
			want:  "user=john&access_token=***",
		},
		{
			name:  "no sensitive data",
			input: "page=1&limit=10&sort=name",
			want:  "page=1&limit=10&sort=name",
		},
		{
			name:  "multiple sensitive params",
			input: "password=pass123&secret=mysecret&key=apikey",
			want:  "password=***&secret=***&key=***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeQueryString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal email",
			input: "john.doe@example.com",
			want:  "jo***@example.com",
		},
		{
			name:  "short email",
			input: "ab@test.com",
			want:  "***@test.com",
		},
		{
			name:  "very short username",
			input: "a@test.com",
			want:  "***@test.com",
		},
		{
			name:  "email in text",
			input: "Contact me at john.doe@example.com for details",
			want:  "Contact me at jo***@example.com for details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeEmail(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeCreditCard(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "card with spaces",
			input: "1234 5678 9012 3456",
			want:  "****-****-****-3456",
		},
		{
			name:  "card with dashes",
			input: "1234-5678-9012-3456",
			want:  "****-****-****-3456",
		},
		{
			name:  "card no separator",
			input: "1234567890123456",
			want:  "****-****-****-3456",
		},
		{
			name:  "card in text",
			input: "Use card 1234-5678-9012-3456 for payment",
			want:  "Use card ****-****-****-3456 for payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeCreditCard(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeCPF(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "CPF with dots and dash",
			input: "123.456.789-00",
			want:  "***.***.**-**",
		},
		{
			name:  "CPF without separator",
			input: "12345678900",
			want:  "***.***.**-**",
		},
		{
			name:  "CPF in text",
			input: "My CPF is 123.456.789-00",
			want:  "My CPF is ***.***.**-**",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeCPF(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeCNPJ(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "CNPJ with dots and dash",
			input: "12.345.678/0001-00",
			want:  "**.***.***/****-**",
		},
		{
			name:  "CNPJ without separator",
			input: "12345678000100",
			want:  "**.***.***/****-**",
		},
		{
			name:  "CNPJ in text",
			input: "Company CNPJ: 12.345.678/0001-00",
			want:  "Company CNPJ: **.***.***/****-**",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeCNPJ(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeString(t *testing.T) {
	input := "Email: john@example.com, CPF: 123.456.789-00, Card: 1234-5678-9012-3456"
	sanitized := SanitizeString(input)

	assert.Contains(t, sanitized, "***")
	assert.NotContains(t, sanitized, "john@example.com")
	assert.NotContains(t, sanitized, "123.456.789-00")
}

func TestSanitizedString(t *testing.T) {
	field := SanitizedString("email", "john@example.com")
	assert.NotNil(t, field)
}

func TestSensitiveString(t *testing.T) {
	field := SensitiveString("password")
	assert.NotNil(t, field)
}

func TestNewSanitizedObject(t *testing.T) {
	fields := map[string]interface{}{
		"name":   "John Doe",
		"email":  "john@example.com",
		"age":    30,
		"active": true,
	}

	encoder := NewSanitizedObject(fields)
	assert.NotNil(t, encoder)
}

func BenchmarkSanitization(b *testing.B) {
	b.Run("SanitizeEmail", func(b *testing.B) {
		email := "john.doe@example.com"
		for i := 0; i < b.N; i++ {
			_ = SanitizeEmail(email)
		}
	})

	b.Run("SanitizeCreditCard", func(b *testing.B) {
		card := "1234-5678-9012-3456"
		for i := 0; i < b.N; i++ {
			_ = SanitizeCreditCard(card)
		}
	})

	b.Run("SanitizeQueryString", func(b *testing.B) {
		query := "username=john&password=secret&token=abc123"
		for i := 0; i < b.N; i++ {
			_ = SanitizeQueryString(query)
		}
	})

	b.Run("SanitizeString", func(b *testing.B) {
		str := "Email: john@example.com, CPF: 123.456.789-00"
		for i := 0; i < b.N; i++ {
			_ = SanitizeString(str)
		}
	})
}

func TestFieldConstants(t *testing.T) {
	// Verificar que as constantes estão definidas corretamente
	assert.Equal(t, "request_id", FieldRequestID)
	assert.Equal(t, "trace_id", FieldTraceID)
	assert.Equal(t, "user_id", FieldUserID)
	assert.Equal(t, "module", FieldModule)
	assert.Equal(t, "function", FieldFunction)
	assert.Equal(t, "http_method", FieldHTTPMethod)
	assert.Equal(t, "http_path", FieldHTTPPath)
	assert.Equal(t, "http_status", FieldHTTPStatus)
	assert.Equal(t, "duration_ms", FieldDuration)
	assert.Equal(t, "error", FieldError)
}

func TestHTTPFieldsWithEmptyValues(t *testing.T) {
	fields := HTTPFields{
		Method: "GET",
		Path:   "/test",
		Status: 200,
		// Outros campos vazios
	}

	zapFields := fields.ToFields()

	// Deve ter pelo menos os 3 campos obrigatórios
	assert.GreaterOrEqual(t, len(zapFields), 3)
}

func TestErrorFieldsWithPartialData(t *testing.T) {
	fields := ErrorFields{
		Error:     assert.AnError,
		ErrorType: "TestError",
		// Outros campos vazios
	}

	zapFields := fields.ToFields()

	// Deve ter os campos fornecidos
	assert.GreaterOrEqual(t, len(zapFields), 2)
}

func ExampleWithHTTP() {
	logger := NewDevelopment()

	fields := HTTPFields{
		Method:   "GET",
		Path:     "/api/users",
		Status:   200,
		Duration: 150 * time.Millisecond,
	}

	logger.Info("request completed", WithHTTP(fields)...)
}

func ExampleWithError() {
	logger := NewDevelopment()

	fields := ErrorFields{
		Error:     assert.AnError,
		ErrorType: "ValidationError",
		ErrorCode: "ERR001",
		Operation: "CreateUser",
	}

	logger.Error("operation failed", WithError(fields)...)
}

func ExampleSanitizeString() {
	sensitive := "User email: john.doe@example.com, CPF: 123.456.789-00"
	safe := SanitizeString(sensitive)

	logger := NewDevelopment()
	logger.Info("user data", zap.String("data", safe))
}
