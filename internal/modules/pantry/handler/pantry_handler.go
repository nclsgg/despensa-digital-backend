package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	itemDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/dto"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

type pantryHandler struct {
	service     domain.PantryService
	itemService itemDomain.ItemService
}

func NewPantryHandler(service domain.PantryService, itemService itemDomain.ItemService) domain.PantryHandler {
	return &pantryHandler{service: service, itemService: itemService}
}

// @Summary Create a new pantry
// @Tags Pantry
// @Accept json
// @Produce json
// @Param body body dto.CreatePantryRequest true "Pantry data"
// @Success 201 {object} dto.PantrySummaryResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries [post]
func (h *pantryHandler) CreatePantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	var req dto.CreatePantryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid pantry creation request",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "CreatePantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Name is required")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	pantry, err := h.service.CreatePantry(c.Request.Context(), req.Name, userID)
	if err != nil {
		logger.Error("Failed to create pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "CreatePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to create pantry")
		return
	}

	summary := dto.PantrySummaryResponse{
		ID:        pantry.ID.String(),
		Name:      pantry.Name,
		OwnerID:   pantry.OwnerID.String(),
		ItemCount: 0,
		CreatedAt: pantry.CreatedAt.Format(time.RFC3339),
		UpdatedAt: pantry.UpdatedAt.Format(time.RFC3339),
	}

	logger.Info("Pantry created successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "CreatePantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantry.ID.String()),
	)

	response.Success(c, 201, summary)
}

// @Summary List all pantries from the current user
// @Tags Pantry
// @Produce json
// @Success 200 {array} dto.PantrySummaryResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries [get]
func (h *pantryHandler) ListPantries(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	rawID, _ := c.Get("userID")
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	pantries, err := h.service.ListPantriesWithItemCount(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to list pantries",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListPantries"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to list pantries")
		return
	}

	logger.Info("Pantries listed successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "ListPantries"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(pantries)),
	)

	response.Success(c, 200, pantries)
}

// @Summary Get the user's main pantry
// @Description Returns the first pantry of the current user (assuming single pantry usage)
// @Tags Pantry
// @Produce json
// @Success 200 {object} dto.PantryDetailResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/my-pantry [get]
func (h *pantryHandler) GetMyPantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	rawID, _ := c.Get("userID")
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	pantryWithCount, err := h.service.GetMyPantry(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user's pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetMyPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.Fail(c, 404, "NOT_FOUND", "Pantry not found")
		return
	}

	items, err := h.itemService.ListByPantryID(c.Request.Context(), pantryWithCount.Pantry.ID, userID)
	if err != nil {
		logger.Error("Failed to list items in user's pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetMyPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryWithCount.Pantry.ID.String()),
			zap.Error(err),
		)
		switch {
		case errors.Is(err, itemDomain.ErrUnauthorized):
			response.Fail(c, 403, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to list pantry items")
		}
		return
	}

	itemResponses := make([]dto.PantryItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, dto.PantryItemResponse{
			ID:             item.ID,
			PantryID:       item.PantryID,
			Name:           item.Name,
			Quantity:       item.Quantity,
			Unit:           item.Unit,
			PricePerUnit:   item.PricePerUnit,
			TotalPrice:     item.TotalPrice,
			AddedBy:        item.AddedBy,
			CategoryID:     item.CategoryID,
			ExpirationDate: item.ExpiresAt,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	res := dto.PantryDetailResponse{
		PantrySummaryResponse: dto.PantrySummaryResponse{
			ID:        pantryWithCount.Pantry.ID.String(),
			Name:      pantryWithCount.Pantry.Name,
			OwnerID:   pantryWithCount.Pantry.OwnerID.String(),
			ItemCount: pantryWithCount.ItemCount,
			CreatedAt: pantryWithCount.Pantry.CreatedAt.Format(time.RFC3339),
			UpdatedAt: pantryWithCount.Pantry.UpdatedAt.Format(time.RFC3339),
		},
		Items: itemResponses,
	}

	logger.Info("User's pantry retrieved successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "GetMyPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryWithCount.Pantry.ID.String()),
		zap.Int(appLogger.FieldCount, len(items)),
	)

	response.OK(c, res)
}

// @Summary Get a specific pantry
// @Tags Pantry
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {object} dto.PantryDetailResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /pantries/{id} [get]
func (h *pantryHandler) GetPantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	pantryWithCount, err := h.service.GetPantryWithItemCount(c.Request.Context(), pantryID, userID)
	if err != nil {
		logger.Error("Failed to get pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		response.Fail(c, 404, "NOT_FOUND", "Pantry not found or user has no access")
		return
	}

	items, err := h.itemService.ListByPantryID(c.Request.Context(), pantryID, userID)
	if err != nil {
		logger.Error("Failed to list items in pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "GetPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		switch {
		case errors.Is(err, itemDomain.ErrUnauthorized):
			response.Fail(c, 403, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to list pantry items")
		}
		return
	}

	itemResponses := make([]dto.PantryItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, dto.PantryItemResponse{
			ID:             item.ID,
			PantryID:       item.PantryID,
			Name:           item.Name,
			Quantity:       item.Quantity,
			Unit:           item.Unit,
			PricePerUnit:   item.PricePerUnit,
			TotalPrice:     item.TotalPrice,
			AddedBy:        item.AddedBy,
			CategoryID:     item.CategoryID,
			ExpirationDate: item.ExpiresAt,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	res := dto.PantryDetailResponse{
		PantrySummaryResponse: dto.PantrySummaryResponse{
			ID:        pantryWithCount.Pantry.ID.String(),
			Name:      pantryWithCount.Pantry.Name,
			OwnerID:   pantryWithCount.Pantry.OwnerID.String(),
			ItemCount: pantryWithCount.ItemCount,
			CreatedAt: pantryWithCount.Pantry.CreatedAt.Format(time.RFC3339),
			UpdatedAt: pantryWithCount.Pantry.UpdatedAt.Format(time.RFC3339),
		},
		Items: itemResponses,
	}

	logger.Info("Pantry retrieved successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "GetPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.Int(appLogger.FieldCount, len(items)),
	)

	response.OK(c, res)
}

// @Summary Update pantry name
// @Tags Pantry
// @Accept json
// @Produce json
// @Param id path string true "Pantry ID"
// @Param body body dto.UpdatePantryRequest true "New name"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id} [put]
func (h *pantryHandler) UpdatePantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "UpdatePantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	var req dto.UpdatePantryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid pantry update request",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "UpdatePantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Name is required")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	err = h.service.UpdatePantry(c.Request.Context(), pantryID, userID, req.Name)
	if err != nil {
		logger.Error("Failed to update pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "UpdatePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to update pantry")
		return
	}

	logger.Info("Pantry updated successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "UpdatePantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
	)

	response.OK(c, response.MessagePayload{Message: "Pantry updated successfully"})
}

// @Summary Soft delete a pantry
// @Tags Pantry
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id} [delete]
func (h *pantryHandler) DeletePantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "DeletePantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	err = h.service.DeletePantry(c.Request.Context(), pantryID, userID)
	if err != nil {
		logger.Error("Failed to delete pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "DeletePantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to delete pantry")
		return
	}

	logger.Info("Pantry deleted successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "DeletePantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
	)

	response.OK(c, response.MessagePayload{Message: "Pantry deleted successfully"})
}

// @Summary Add a user to the pantry
// @Tags Pantry
// @Accept json
// @Produce json
// @Param id path string true "Pantry ID"
// @Param body body dto.ModifyPantryUserRequest true "User to add"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id}/users [post]
func (h *pantryHandler) AddUserToPantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	var req dto.ModifyPantryUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid add user request",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "User ID is required")
		return
	}

	rawID, _ := c.Get("userID")
	ownerID := rawID.(uuid.UUID)

	err = h.service.AddUserToPantry(c.Request.Context(), pantryID, ownerID, req.Email)
	if err != nil {
		logger.Error("Failed to add user to pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "AddUserToPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(req.Email)),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to add user to pantry")
		return
	}

	logger.Info("User added to pantry successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "AddUserToPantry"),
		zap.String(appLogger.FieldUserID, ownerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(req.Email)),
	)

	response.OK(c, response.MessagePayload{Message: "User added to pantry successfully"})
}

// @Summary Remove a user from the pantry by email
// @Tags Pantry
// @Accept json
// @Produce json
// @Param id path string true "Pantry ID"
// @Param body body dto.ModifyPantryUserRequest true "User to remove"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id}/users [delete]
// @Summary Remove a user from the pantry by email
// @Tags Pantry
// @Accept json
// @Produce json
// @Param id path string true "Pantry ID"
// @Param body body dto.ModifyPantryUserRequest true "User to remove"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id}/users [delete]
func (h *pantryHandler) RemoveUserFromPantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	var req dto.ModifyPantryUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid remove user request",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "User email is required")
		return
	}

	rawID, _ := c.Get("userID")
	ownerID := rawID.(uuid.UUID)

	err = h.service.RemoveUserFromPantry(c.Request.Context(), pantryID, ownerID, req.Email)
	if err != nil {
		logger.Error("Failed to remove user from pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(req.Email)),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to remove user from pantry")
		return
	}

	logger.Info("User removed from pantry successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "RemoveUserFromPantry"),
		zap.String(appLogger.FieldUserID, ownerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(req.Email)),
	)

	response.OK(c, response.MessagePayload{Message: "User removed from pantry successfully"})
}

// @Summary Remove a specific user from the pantry by user ID
// @Tags Pantry
// @Produce json
// @Param id path string true "Pantry ID"
// @Param userId path string true "User ID to remove"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id}/users/{userId} [delete]
func (h *pantryHandler) RemoveSpecificUserFromPantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		logger.Warn("Invalid user ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid user ID")
		return
	}

	rawID, _ := c.Get("userID")
	ownerID := rawID.(uuid.UUID)

	err = h.service.RemoveSpecificUserFromPantry(c.Request.Context(), pantryID, ownerID, targetUserID)
	if err != nil {
		logger.Error("Failed to remove specific user from pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
			zap.String(appLogger.FieldUserID, ownerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("target_user_id", targetUserID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to remove user from pantry")
		return
	}

	logger.Info("Specific user removed from pantry successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "RemoveSpecificUserFromPantry"),
		zap.String(appLogger.FieldUserID, ownerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String("target_user_id", targetUserID.String()),
	)

	response.OK(c, response.MessagePayload{Message: "User removed from pantry successfully"})
}

// @Summary Transfer pantry ownership to another member
// @Tags Pantry
// @Accept json
// @Produce json
// @Param id path string true "Pantry ID"
// @Param body body dto.TransferOwnershipRequest true "New owner ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id}/transfer-ownership [post]
func (h *pantryHandler) TransferOwnership(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	var req dto.TransferOwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid transfer ownership request",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.Error(err),
		)
		response.BadRequest(c, "New owner ID is required")
		return
	}

	newOwnerID, err := uuid.Parse(req.NewOwnerID)
	if err != nil {
		logger.Warn("Invalid new owner ID format",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid new owner ID")
		return
	}

	rawID, _ := c.Get("userID")
	currentOwnerID := rawID.(uuid.UUID)

	err = h.service.TransferOwnership(c.Request.Context(), pantryID, currentOwnerID, newOwnerID)
	if err != nil {
		logger.Error("Failed to transfer pantry ownership",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "TransferOwnership"),
			zap.String(appLogger.FieldUserID, currentOwnerID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.String("new_owner_id", newOwnerID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to transfer ownership")
		return
	}

	logger.Info("Pantry ownership transferred successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "TransferOwnership"),
		zap.String(appLogger.FieldUserID, currentOwnerID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String("new_owner_id", newOwnerID.String()),
	)

	response.OK(c, response.MessagePayload{Message: "Ownership transferred successfully"})
}

// @Summary List users in a pantry
// @Tags Pantry
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {array} dto.PantryUserResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id}/users [get]
func (h *pantryHandler) ListUsersInPantry(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid pantry ID in URL parameter",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	users, err := h.service.ListUsersInPantry(c.Request.Context(), pantryID, userID)
	if err != nil {
		logger.Error("Failed to list users in pantry",
			zap.String(appLogger.FieldModule, "pantry"),
			zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.String("pantry_id", pantryID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to list users in pantry")
		return
	}

	responses := make([]dto.PantryUserResponse, 0, len(users))
	for _, user := range users {
		first := strings.TrimSpace(user.FirstName)
		last := strings.TrimSpace(user.LastName)
		name := strings.TrimSpace(strings.Join([]string{first, last}, " "))
		if name == "" {
			name = user.Email
		}

		responses = append(responses, dto.PantryUserResponse{
			ID:       user.ID.String(),
			UserID:   user.UserID.String(),
			PantryID: user.PantryID.String(),
			Email:    user.Email,
			Name:     name,
			Role:     user.Role,
		})
	}

	logger.Info("Users in pantry listed successfully",
		zap.String(appLogger.FieldModule, "pantry"),
		zap.String(appLogger.FieldFunction, "ListUsersInPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.Int(appLogger.FieldCount, len(users)),
	)

	response.OK(c, responses)
}
