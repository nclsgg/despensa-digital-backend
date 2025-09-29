package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, data interface{}) {
	__logParams := map[string]any{"c": c, "status": status, "data": data}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "Success"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "Success"), zap.Any("params", __logParams))
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

func Fail(c *gin.Context, status int, code, message string) {
	__logParams := map[string]any{"c": c, "status": status, "code": code, "message": message}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "Fail"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "Fail"), zap.Any("params", __logParams))
	c.JSON(status, APIResponse{
		Success: false,
		Data:    nil,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func OK(c *gin.Context, data interface{}) {
	__logParams := map[string]any{"c": c, "data": data}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "OK"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "OK"), zap.Any("params", __logParams))
	Success(c, http.StatusOK, data)
}

func BadRequest(c *gin.Context, message string) {
	__logParams := map[string]any{"c": c, "message": message}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "BadRequest"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "BadRequest"), zap.Any("params", __logParams))
	Fail(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func InternalError(c *gin.Context, message string) {
	__logParams := map[string]any{"c": c, "message": message}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "InternalError"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "InternalError"), zap.Any("params", __logParams))
	Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func Unauthorized(c *gin.Context, message string) {
	__logParams := map[string]any{"c": c, "message": message}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "Unauthorized"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "Unauthorized"), zap.Any("params", __logParams))
	Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}
