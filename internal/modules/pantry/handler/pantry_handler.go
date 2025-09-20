package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/dto"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

type pantryHandler struct {
	service domain.PantryService
}

func NewPantryHandler(service domain.PantryService) domain.PantryHandler {
	return &pantryHandler{service}
}

// @Summary Create a new pantry
// @Tags Pantry
// @Accept json
// @Produce json
// @Param body body dto.CreatePantryRequest true "Pantry data"
// @Success 201 {object} dto.PantryResponse
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

	res := dto.PantryResponse{
		ID:      pantry.ID,
		Name:    pantry.Name,
		OwnerID: pantry.OwnerID,
	}
	response.Success(c, 201, res)
}

// @Summary List all pantries from the current user
// @Tags Pantry
// @Produce json
// @Success 200 {array} dto.PantryResponse
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

	var res []dto.PantryResponse
	for _, pantryWithCount := range pantries {
		res = append(res, dto.PantryResponse{
			ID:        pantryWithCount.Pantry.ID,
			Name:      pantryWithCount.Pantry.Name,
			OwnerID:   pantryWithCount.Pantry.OwnerID,
			ItemCount: pantryWithCount.ItemCount,
		})
	}

	response.OK(c, res)
}

// @Summary Get a specific pantry
// @Tags Pantry
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {object} dto.PantryResponse
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

	res := dto.PantryResponse{
		ID:        pantryWithCount.Pantry.ID,
		Name:      pantryWithCount.Pantry.Name,
		OwnerID:   pantryWithCount.Pantry.OwnerID,
		ItemCount: pantryWithCount.ItemCount,
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

	var res []dto.PantryUserResponse
	for _, user := range users {
		res = append(res, dto.PantryUserResponse{
			UserID: user.UserID,
			Email:  user.Email,
			Role:   user.Role,
		})
	}

	response.OK(c, res)
}
