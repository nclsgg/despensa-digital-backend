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
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
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
	var req dto.CreatePantryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Name is required")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	pantry, err := h.service.CreatePantry(c.Request.Context(), req.Name, userID)
	if err != nil {
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

	response.Success(c, 201, summary)
}

// @Summary List all pantries from the current user
// @Tags Pantry
// @Produce json
// @Success 200 {array} dto.PantrySummaryResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries [get]
func (h *pantryHandler) ListPantries(c *gin.Context) {
	rawID, _ := c.Get("userID")
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	pantries, err := h.service.ListPantriesWithItemCount(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "Failed to list pantries")
		return
	}

	var res []dto.PantrySummaryResponse
	for _, pantryWithCount := range pantries {
		pantry := pantryWithCount.Pantry
		res = append(res, dto.PantrySummaryResponse{
			ID:        pantry.ID.String(),
			Name:      pantry.Name,
			OwnerID:   pantry.OwnerID.String(),
			ItemCount: pantryWithCount.ItemCount,
			CreatedAt: pantry.CreatedAt.Format(time.RFC3339),
			UpdatedAt: pantry.UpdatedAt.Format(time.RFC3339),
		})
	}

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
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	pantryWithCount, err := h.service.GetPantryWithItemCount(c.Request.Context(), pantryID, userID)
	if err != nil {
		response.Fail(c, 404, "NOT_FOUND", "Pantry not found or user has no access")
		return
	}

	items, err := h.itemService.ListByPantryID(c.Request.Context(), pantryID, userID)
	if err != nil {
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
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	var req dto.UpdatePantryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Name is required")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	err = h.service.UpdatePantry(c.Request.Context(), pantryID, userID, req.Name)
	if err != nil {
		response.InternalError(c, "Failed to update pantry")
		return
	}

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
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	err = h.service.DeletePantry(c.Request.Context(), pantryID, userID)
	if err != nil {
		response.InternalError(c, "Failed to delete pantry")
		return
	}

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
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	var req dto.ModifyPantryUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "User ID is required")
		return
	}

	rawID, _ := c.Get("userID")
	ownerID := rawID.(uuid.UUID)

	err = h.service.AddUserToPantry(c.Request.Context(), pantryID, ownerID, req.Email)
	if err != nil {
		response.InternalError(c, "Failed to add user to pantry")
		return
	}

	response.OK(c, response.MessagePayload{Message: "User added to pantry"})
}

// @Summary Remove a user from the pantry
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
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	var req dto.ModifyPantryUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "User ID is required")
		return
	}

	rawID, _ := c.Get("userID")
	ownerID := rawID.(uuid.UUID)

	err = h.service.RemoveUserFromPantry(c.Request.Context(), pantryID, ownerID, req.Email)
	if err != nil {
		response.InternalError(c, "Failed to remove user from pantry")
		return
	}

	response.OK(c, response.MessagePayload{Message: "User removed from pantry"})
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
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	users, err := h.service.ListUsersInPantry(c.Request.Context(), pantryID, userID)
	if err != nil {
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

	response.OK(c, responses)
}
