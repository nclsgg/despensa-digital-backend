package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWithLogger(t *testing.T) {
	logger := NewDevelopment()
	ctx := context.Background()

	newCtx := WithLogger(ctx, logger)
	assert.NotNil(t, newCtx)

	// Verificar que o logger está no contexto
	retrievedLogger := FromContext(newCtx)
	assert.NotNil(t, retrievedLogger)
}

func TestFromContext(t *testing.T) {
	t.Run("with logger in context", func(t *testing.T) {
		logger := NewDevelopment()
		ctx := WithLogger(context.Background(), logger)

		retrieved := FromContext(ctx)
		assert.NotNil(t, retrieved)
	})

	t.Run("without logger in context", func(t *testing.T) {
		ctx := context.Background()

		// Deve retornar o logger global
		retrieved := FromContext(ctx)
		assert.NotNil(t, retrieved)
	})

	t.Run("nil context", func(t *testing.T) {
		retrieved := FromContext(nil)
		assert.NotNil(t, retrieved)
	})
}

func TestWithRequestID(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	newCtx := WithRequestID(ctx)
	assert.NotNil(t, newCtx)

	requestID := GetRequestID(newCtx)
	assert.NotEmpty(t, requestID)

	// Verificar que chamadas subsequentes mantêm o mesmo ID
	newCtx2 := WithRequestID(newCtx)
	requestID2 := GetRequestID(newCtx2)
	assert.Equal(t, requestID, requestID2)
}

func TestWithRequestIDValue(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	customID := "custom-request-id-123"
	newCtx := WithRequestIDValue(ctx, customID)

	requestID := GetRequestID(newCtx)
	assert.Equal(t, customID, requestID)
}

func TestGetRequestID(t *testing.T) {
	t.Run("with request ID in context", func(t *testing.T) {
		ctx := WithRequestID(context.Background())
		requestID := GetRequestID(ctx)
		assert.NotEmpty(t, requestID)
	})

	t.Run("without request ID in context", func(t *testing.T) {
		ctx := context.Background()
		requestID := GetRequestID(ctx)
		assert.Empty(t, requestID)
	})

	t.Run("nil context", func(t *testing.T) {
		requestID := GetRequestID(nil)
		assert.Empty(t, requestID)
	})
}

func TestWithTraceID(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	traceID := "trace-xyz-789"
	newCtx := WithTraceID(ctx, traceID)

	retrieved := GetTraceID(newCtx)
	assert.Equal(t, traceID, retrieved)
}

func TestGetTraceID(t *testing.T) {
	t.Run("with trace ID in context", func(t *testing.T) {
		ctx := WithTraceID(context.Background(), "trace-123")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "trace-123", traceID)
	})

	t.Run("without trace ID in context", func(t *testing.T) {
		ctx := context.Background()
		traceID := GetTraceID(ctx)
		assert.Empty(t, traceID)
	})

	t.Run("empty trace ID", func(t *testing.T) {
		ctx := WithTraceID(context.Background(), "")
		traceID := GetTraceID(ctx)
		assert.Empty(t, traceID)
	})
}

func TestWithUserID(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	userID := "user-456"
	newCtx := WithUserID(ctx, userID)

	retrieved := GetUserID(newCtx)
	assert.Equal(t, userID, retrieved)
}

func TestGetUserID(t *testing.T) {
	t.Run("with user ID in context", func(t *testing.T) {
		ctx := WithUserID(context.Background(), "user-789")
		userID := GetUserID(ctx)
		assert.Equal(t, "user-789", userID)
	})

	t.Run("without user ID in context", func(t *testing.T) {
		ctx := context.Background()
		userID := GetUserID(ctx)
		assert.Empty(t, userID)
	})

	t.Run("empty user ID", func(t *testing.T) {
		ctx := WithUserID(context.Background(), "")
		userID := GetUserID(ctx)
		assert.Empty(t, userID)
	})
}

func TestWithModule(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	moduleName := "users"
	newCtx := WithModule(ctx, moduleName)

	// Verificar que o logger foi atualizado
	updatedLogger := FromContext(newCtx)
	assert.NotNil(t, updatedLogger)
}

func TestWithFunction(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	funcName := "CreateUser"
	newCtx := WithFunction(ctx, funcName)

	// Verificar que o logger foi atualizado
	updatedLogger := FromContext(newCtx)
	assert.NotNil(t, updatedLogger)
}

func TestWithModuleAndFunction(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	newCtx := WithModuleAndFunction(ctx, "users", "CreateUser")

	// Verificar que o logger foi atualizado
	updatedLogger := FromContext(newCtx)
	assert.NotNil(t, updatedLogger)
}

func TestWithFields(t *testing.T) {
	logger := NewDevelopment()
	ctx := WithLogger(context.Background(), logger)

	fields := []zap.Field{
		zap.String("custom_field", "value"),
		zap.Int("count", 42),
	}

	newCtx := WithFields(ctx, fields...)

	// Verificar que o logger foi atualizado
	updatedLogger := FromContext(newCtx)
	assert.NotNil(t, updatedLogger)
}

func TestContextPropagation(t *testing.T) {
	logger := NewDevelopment()
	ctx := context.Background()

	// Adicionar logger
	ctx = WithLogger(ctx, logger)

	// Adicionar request ID
	ctx = WithRequestID(ctx)
	requestID := GetRequestID(ctx)
	require.NotEmpty(t, requestID)

	// Adicionar trace ID
	ctx = WithTraceID(ctx, "trace-123")
	traceID := GetTraceID(ctx)
	assert.Equal(t, "trace-123", traceID)

	// Adicionar user ID
	ctx = WithUserID(ctx, "user-456")
	userID := GetUserID(ctx)
	assert.Equal(t, "user-456", userID)

	// Adicionar módulo e função
	ctx = WithModuleAndFunction(ctx, "orders", "ProcessOrder")

	// Verificar que o logger tem todos os campos
	finalLogger := FromContext(ctx)
	assert.NotNil(t, finalLogger)

	// Fazer um log para verificar
	finalLogger.Info("test propagation complete")
}

func BenchmarkContextOperations(b *testing.B) {
	logger := NewDevelopment()

	b.Run("WithLogger", func(b *testing.B) {
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			_ = WithLogger(ctx, logger)
		}
	})

	b.Run("FromContext", func(b *testing.B) {
		ctx := WithLogger(context.Background(), logger)
		for i := 0; i < b.N; i++ {
			_ = FromContext(ctx)
		}
	})

	b.Run("WithRequestID", func(b *testing.B) {
		ctx := WithLogger(context.Background(), logger)
		for i := 0; i < b.N; i++ {
			_ = WithRequestID(ctx)
		}
	})

	b.Run("Full propagation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx := context.Background()
			ctx = WithLogger(ctx, logger)
			ctx = WithRequestID(ctx)
			ctx = WithTraceID(ctx, "trace-123")
			ctx = WithUserID(ctx, "user-456")
			ctx = WithModuleAndFunction(ctx, "test", "BenchmarkFunc")
			_ = FromContext(ctx)
		}
	})
}
