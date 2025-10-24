package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GinRequestIDMiddleware adds a unique request ID to each request context
func GinRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := WithRequestID(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		c.Next()
	}
}

// GinLoggingMiddleware logs HTTP requests
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Add logger to context
		ctx := WithLogger(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Get request ID from context
		requestID := GetRequestID(c.Request.Context())

		// Build HTTP fields
		httpFields := HTTPFields{
			Method:       c.Request.Method,
			Path:         path,
			Query:        query,
			RemoteAddr:   c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			Status:       c.Writer.Status(),
			ResponseSize: int64(c.Writer.Size()),
			Duration:     duration,
		}

		fields := []zap.Field{
			zap.String(FieldRequestID, requestID),
			zap.String(FieldModule, "http"),
		}
		fields = append(fields, httpFields.ToFields()...)

		// Log based on status code
		if c.Writer.Status() >= 500 {
			logger.Error("HTTP request failed", fields...)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("HTTP request client error", fields...)
		} else {
			logger.Info("HTTP request completed", fields...)
		}
	}
}

// GinRecoveryMiddleware recovers from panics and logs them
func GinRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c.Request.Context())

				httpFields := HTTPFields{
					Method:     c.Request.Method,
					Path:       c.Request.URL.Path,
					Query:      c.Request.URL.RawQuery,
					RemoteAddr: c.ClientIP(),
					UserAgent:  c.Request.UserAgent(),
					Status:     500,
				}

				fields := []zap.Field{
					zap.String(FieldRequestID, requestID),
					zap.String(FieldModule, "http"),
					zap.String(FieldFunction, "recovery"),
					zap.Any("panic", err),
					zap.Stack("stacktrace"),
				}
				fields = append(fields, httpFields.ToFields()...)

				logger.Error("Panic recovered", fields...)

				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}

// GinLoggerFromContext returns the logger from the Gin context
func GinLoggerFromContext(c *gin.Context) *zap.Logger {
	return FromContext(c.Request.Context())
}
