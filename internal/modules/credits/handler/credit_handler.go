package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/dto"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

type CreditHandler struct {
	service domain.CreditService
}

func NewCreditHandler(service domain.CreditService) *CreditHandler {
	return &CreditHandler{service: service}
}

func (h *CreditHandler) GetWallet(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	userID, ok := extractUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found in context")
		return
	}

	wallet, err := h.service.GetWallet(c.Request.Context(), userID)
	if err != nil {
		logger.Error("failed to get wallet",
			zap.String(appLogger.FieldModule, "credits"),
			zap.String(appLogger.FieldFunction, "GetWallet"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "failed to retrieve wallet")
		return
	}

	response.OK(c, wallet)
}

func (h *CreditHandler) ListTransactions(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	requesterID, ok := extractUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found in context")
		return
	}

	targetUserID := requesterID
	if targetParam := strings.TrimSpace(c.Query("user_id")); targetParam != "" {
		targetUUID, err := uuid.Parse(targetParam)
		if err != nil {
			response.BadRequest(c, "invalid user_id")
			return
		}
		if !isAdmin(c) && targetUUID != requesterID {
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "insufficient permissions to view other user's transactions")
			return
		}
		targetUserID = targetUUID
	}

	filter := dto.TransactionFilter{}

	if typeParam := strings.TrimSpace(c.Query("type")); typeParam != "" {
		lower := strings.ToLower(typeParam)
		filter.Type = &lower
	}

	if limitParam := strings.TrimSpace(c.Query("limit")); limitParam != "" {
		if limitVal, err := strconv.Atoi(limitParam); err == nil {
			filter.Limit = limitVal
		}
	}

	if offsetParam := strings.TrimSpace(c.Query("offset")); offsetParam != "" {
		if offsetVal, err := strconv.Atoi(offsetParam); err == nil {
			filter.Offset = offsetVal
		}
	}

	if fromParam := strings.TrimSpace(c.Query("from")); fromParam != "" {
		if parsed, err := time.Parse(time.RFC3339, fromParam); err == nil {
			filter.From = &parsed
		}
	}

	if toParam := strings.TrimSpace(c.Query("to")); toParam != "" {
		if parsed, err := time.Parse(time.RFC3339, toParam); err == nil {
			filter.To = &parsed
		}
	}

	transactions, err := h.service.ListTransactions(c.Request.Context(), targetUserID, filter)
	if err != nil {
		logger.Error("failed to list transactions",
			zap.String(appLogger.FieldModule, "credits"),
			zap.String(appLogger.FieldFunction, "ListTransactions"),
			zap.String(appLogger.FieldUserID, targetUserID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "failed to list transactions")
		return
	}

	response.OK(c, gin.H{
		"transactions": transactions,
	})
}

func (h *CreditHandler) AddCredits(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	if !isAdmin(c) {
		response.Fail(c, http.StatusForbidden, "FORBIDDEN", "only administrators can grant credits")
		return
	}

	actorID, ok := extractUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found in context")
		return
	}

	var payload dto.AddCreditsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, "invalid payload: "+err.Error())
		return
	}

	targetUserID := actorID
	if payload.UserID != nil && strings.TrimSpace(*payload.UserID) != "" {
		targetUUID, err := uuid.Parse(strings.TrimSpace(*payload.UserID))
		if err != nil {
			response.BadRequest(c, "invalid user_id")
			return
		}
		targetUserID = targetUUID
	}

	wallet, err := h.service.AddCredit(c.Request.Context(), actorID, targetUserID, payload.Amount, payload.Description)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCreditAmount):
			response.Fail(c, http.StatusBadRequest, "INVALID_AMOUNT", "amount must be positive")
		default:
			logger.Error("failed to add credits",
				zap.String(appLogger.FieldModule, "credits"),
				zap.String(appLogger.FieldFunction, "AddCredits"),
				zap.String(appLogger.FieldUserID, targetUserID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "failed to add credits")
		}
		return
	}

	response.OK(c, wallet)
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
		parsed, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, false
		}
		return parsed, true
	default:
		return uuid.Nil, false
	}
}

func isAdmin(c *gin.Context) bool {
	role := strings.ToLower(strings.TrimSpace(c.GetString("role")))
	return role == "admin"
}
