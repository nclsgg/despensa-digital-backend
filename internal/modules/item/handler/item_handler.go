package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

type itemHandler struct {
	service domain.ItemService
}

func NewItemHandler(service domain.ItemService) (result0 domain.ItemHandler) {
	__logParams := map[string]any{"service": service}

	// @Summary Create a new item
	// @Tags Items
	// @Accept json
	// @Produce json
	// @Param body body dto.CreateItemDTO true "Item data"
	// @Success 201 {object} dto.ItemResponse
	// @Failure 400 {object} response.APIResponse
	// @Failure 500 {object} response.APIResponse
	// @Router /items [post]
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewItemHandler"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewItemHandler"), zap.Any("params", __logParams))
	result0 = &itemHandler{service}
	return
}

func (h *itemHandler) CreateItem(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemHandler.CreateItem"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemHandler.CreateItem"), zap.Any("params", __logParams))
	var input dto.CreateItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.CreateItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	item, err := h.service.Create(c.Request.Context(), input, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.CreateItem"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrInvalidPantry):
			response.BadRequest(c, "Invalid pantry ID")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to create item")
		}
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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemHandler.UpdateItem"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemHandler.UpdateItem"), zap.Any("params", __logParams))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.UpdateItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Item ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	var input dto.UpdateItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.UpdateItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input")
		return
	}

	item, err := h.service.Update(c.Request.Context(), id, input, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.UpdateItem"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrItemNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to update item")
		}
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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemHandler.GetItem"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemHandler.GetItem"), zap.Any("params", __logParams))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.GetItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Item ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	item, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.GetItem"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrItemNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to fetch item")
		}
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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemHandler.DeleteItem"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemHandler.DeleteItem"), zap.Any("params", __logParams))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.DeleteItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Item ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.DeleteItem"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrItemNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to delete item")
		}
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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemHandler.ListItems"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemHandler.ListItems"), zap.Any("params", __logParams))
	pantryIDStr := c.Param("id")
	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)
	if pantryIDStr == "" {
		response.BadRequest(c, "Pantry ID is required")
		return
	}

	pantryID, err := uuid.Parse(pantryIDStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.ListItems"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Pantry ID")
		return
	}

	items, err := h.service.ListByPantryID(c.Request.Context(), pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.ListItems"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to list items")
		}
		return
	}

	response.OK(c, items)
}

// @Summary Filter items by pantry ID with filters
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Pantry ID"
// @Param body body dto.ItemFilterDTO true "Filter criteria"
// @Success 200 {array} []dto.ItemResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /items/pantry/{id}/filter [post]
func (h *itemHandler) FilterItems(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemHandler.FilterItems"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemHandler.FilterItems"), zap.Any("params", __logParams))
	pantryIDStr := c.Param("id")
	if pantryIDStr == "" {
		response.BadRequest(c, "Pantry ID is required")
		return
	}

	pantryID, err := uuid.Parse(pantryIDStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.FilterItems"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Pantry ID")
		return
	}

	var filters dto.ItemFilterDTO
	if err := c.ShouldBindJSON(&filters); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.FilterItems"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid filter parameters")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	items, err := h.service.FilterByPantryID(c.Request.Context(), pantryID, filters, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemHandler.FilterItems"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to filter items")
		}
		return
	}

	response.OK(c, items)
}
