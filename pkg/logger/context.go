package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	loggerKey    contextKey = "logger"
	requestIDKey contextKey = "request_id"
	traceIDKey   contextKey = "trace_id"
	userIDKey    contextKey = "user_id"
)

// WithLogger adiciona um logger ao contexto.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext recupera o logger do contexto.
// Se não houver logger no contexto, retorna o logger global.
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}

	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok && logger != nil {
		return logger
	}

	return zap.L()
}

// WithRequestID adiciona um request ID ao contexto e ao logger.
// Se já existir um request ID no contexto, ele é reutilizado.
func WithRequestID(ctx context.Context) context.Context {
	// Verificar se já existe um request ID
	if existingID := GetRequestID(ctx); existingID != "" {
		return ctx
	}

	// Gerar novo request ID
	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	// Adicionar ao logger se existir no contexto
	if logger := FromContext(ctx); logger != nil {
		logger = logger.With(zap.String(FieldRequestID, requestID))
		ctx = WithLogger(ctx, logger)
	}

	return ctx
}

// WithRequestIDValue adiciona um request ID específico ao contexto e ao logger.
func WithRequestIDValue(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return WithRequestID(ctx)
	}

	ctx = context.WithValue(ctx, requestIDKey, requestID)

	// Adicionar ao logger se existir no contexto
	if logger := FromContext(ctx); logger != nil {
		logger = logger.With(zap.String(FieldRequestID, requestID))
		ctx = WithLogger(ctx, logger)
	}

	return ctx
}

// GetRequestID recupera o request ID do contexto.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}

	return ""
}

// WithTraceID adiciona um trace ID ao contexto e ao logger.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}

	ctx = context.WithValue(ctx, traceIDKey, traceID)

	// Adicionar ao logger se existir no contexto
	if logger := FromContext(ctx); logger != nil {
		logger = logger.With(zap.String(FieldTraceID, traceID))
		ctx = WithLogger(ctx, logger)
	}

	return ctx
}

// GetTraceID recupera o trace ID do contexto.
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}

	return ""
}

// WithUserID adiciona um user ID ao contexto e ao logger.
func WithUserID(ctx context.Context, userID string) context.Context {
	if userID == "" {
		return ctx
	}

	ctx = context.WithValue(ctx, userIDKey, userID)

	// Adicionar ao logger se existir no contexto
	if logger := FromContext(ctx); logger != nil {
		logger = logger.With(zap.String(FieldUserID, userID))
		ctx = WithLogger(ctx, logger)
	}

	return ctx
}

// GetUserID recupera o user ID do contexto.
func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}

	return ""
}

// WithModule adiciona o nome do módulo ao logger no contexto.
func WithModule(ctx context.Context, module string) context.Context {
	logger := FromContext(ctx).With(zap.String(FieldModule, module))
	return WithLogger(ctx, logger)
}

// WithFunction adiciona o nome da função ao logger no contexto.
func WithFunction(ctx context.Context, function string) context.Context {
	logger := FromContext(ctx).With(zap.String(FieldFunction, function))
	return WithLogger(ctx, logger)
}

// WithModuleAndFunction adiciona módulo e função ao logger no contexto.
func WithModuleAndFunction(ctx context.Context, module, function string) context.Context {
	logger := FromContext(ctx).With(
		zap.String(FieldModule, module),
		zap.String(FieldFunction, function),
	)
	return WithLogger(ctx, logger)
}

// WithFields adiciona campos customizados ao logger no contexto.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	logger := FromContext(ctx).With(fields...)
	return WithLogger(ctx, logger)
}
