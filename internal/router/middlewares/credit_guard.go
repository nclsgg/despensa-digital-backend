package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/domain"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

func CreditGuardMiddleware(creditService domain.CreditService) gin.HandlerFunc {
	__logParams := map[string]any{"creditService": creditService}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "CreditGuardMiddleware"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "CreditGuardMiddleware"), zap.Any("params", __logParams))

	return func(c *gin.Context) {
		userID, ok := extractUserID(c)
		if !ok {
			response.Unauthorized(c, "user not found in context")
			c.Abort()
			return
		}

		wallet, err := creditService.GetWallet(c.Request.Context(), userID)
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "CreditGuardMiddleware"), zap.Error(err), zap.Any("context", map[string]any{"user_id": userID.String()}))
			response.InternalError(c, "failed to validate credits")
			c.Abort()
			return
		}

		if wallet == nil || wallet.Balance <= 0 {
			response.Fail(c, http.StatusPaymentRequired, "INSUFFICIENT_CREDITS", "You don't have enough credits to perform this action")
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}

	switch v := value.(type) {
	case uuid.UUID:
		return v, true
	case string:
		parsed, err := uuid.Parse(strings.TrimSpace(v))
		if err != nil {
			return uuid.Nil, false
		}
		return parsed, true
	default:
		return uuid.Nil, false
	}
}
