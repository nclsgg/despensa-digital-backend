package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

type itemHandler struct {
	service domain.ItemService
}

func NewItemHandler(service domain.ItemService) domain.ItemHandler {
	return &itemHandler{service}
}

// @Summary Create a new item
// @Tags Items
// @Accept json
// @Produce json
// @Param body body dto.CreateItemDTO true "Item data"
// @Success 201 {object} dto.ItemResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /items [post]
func (h *itemHandler) CreateItem(c *gin.Context) {
	var input dto.CreateItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input")
		return
	}

	rawID, _ := c.Get("user_id")
	userID := rawID.(uuid.UUID)

	item, err := h.service.Create(c.Request.Context(), input, userID)
	if err != nil {
		response.InternalError(c, "Failed to create item")
		return
	}

	response.Success(c, http.StatusCreated, item)
}

// @Summary Update an item
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param body body dto.UpdateItemDTO true "Updated fields"
// @Success 200 {object} dto.ItemResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /items/{id} [put]
func (h *itemHandler) UpdateItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid Item ID")
		return
	}

	rawID, _ := c.Get("user_id")
	userID := rawID.(uuid.UUID)

	var input dto.UpdateItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input")
		return
	}

	item, err := h.service.Update(c.Request.Context(), id, input, userID)
	if err != nil {
		response.InternalError(c, "Failed to update item")
		return
	}

	response.OK(c, item)
}

// @Summary Get an item by ID
// @Tags Items
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} dto.ItemResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /items/{id} [get]
func (h *itemHandler) GetItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid Item ID")
		return
	}

	rawID, _ := c.Get("user_id")
	userID := rawID.(uuid.UUID)

	item, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item not found")
		return
	}
	response.OK(c, item)
}

// @Summary Delete an item by ID
// @Tags Items
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /items/{id} [delete]
func (h *itemHandler) DeleteItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid Item ID")
		return
	}

	rawID, _ := c.Get("user_id")
	userID := rawID.(uuid.UUID)

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		response.InternalError(c, "Failed to delete item")
		return
	}
	response.OK(c, response.MessagePayload{Message: "Item deleted successfully"})
}

// @Summary List all items by pantry ID
// @Tags Items
// @Produce json
// @Param pantry_id query string true "Pantry ID"
// @Success 200 {array} []dto.ItemResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /items [get]
func (h *itemHandler) ListItems(c *gin.Context) {
	pantryIDStr := c.Param("id")
	rawID, _ := c.Get("user_id")
	userID := rawID.(uuid.UUID)
	if pantryIDStr == "" {
		response.BadRequest(c, "Pantry ID is required")
		return
	}

	pantryID, err := uuid.Parse(pantryIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid Pantry ID")
		return
	}

	items, err := h.service.ListByPantryID(c.Request.Context(), pantryID, userID)
	if err != nil {
		response.InternalError(c, "Failed to list items")
		return
	}

	response.OK(c, items)
}
