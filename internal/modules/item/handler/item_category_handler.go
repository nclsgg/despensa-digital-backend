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

type itemCategoryHandler struct {
	service domain.ItemCategoryService
}

func NewItemCategoryHandler(service domain.ItemCategoryService) (result0 domain.ItemCategoryHandler) {
	__logParams := map[string]any{"service": service}
	__logStart :=

		// @Summary Create a new item category
		// @Tags Item Categories
		// @Accept json
		// @Produce json
		// @Param body body dto.CreateItemCategoryDTO true "Item Category data"
		// @Success 201 {object} dto.ItemCategoryResponse
		// @Failure 400 {object} response.APIResponse
		// @Failure 500 {object} response.APIResponse
		// @Router /item-categories [post]
		time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewItemCategoryHandler"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewItemCategoryHandler"), zap.Any("params", __logParams))
	result0 = &itemCategoryHandler{service}
	return
}

func (h *itemCategoryHandler) CreateItemCategory(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.CreateItemCategory"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.CreateItemCategory"), zap.Any("params", __logParams))
	var input dto.CreateItemCategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.CreateItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.Create(c.Request.Context(), input, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.CreateItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrInvalidPantry):
			response.BadRequest(c, "Invalid pantry ID")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to create item category")
		}
		return
	}

	response.Success(c, http.StatusCreated, category)
}

// @Summary Create a new default item category
// @Tags Item Categories
// @Accept json
// @Produce json
// @Param body body dto.CreateDefaultItemCategoryDTO true "Default Item Category data"
// @Success 201 {object} dto.ItemCategoryResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /item-categories/default [post]
func (h *itemCategoryHandler) CreateDefaultItemCategory(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.CreateDefaultItemCategory"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.CreateDefaultItemCategory"), zap.Any("params", __logParams))
	var input dto.CreateDefaultItemCategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.CreateDefaultItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.CreateDefault(c.Request.Context(), input, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.CreateDefaultItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		response.InternalError(c, "Failed to create default item category")
		return
	}

	response.Success(c, http.StatusCreated, category)
}

// @Summary Clone a default item category to a pantry
// @Tags Item Categories
// @Accept json
// @Produce json
// @Param pantry_id path string true "Pantry ID"
// @Param default_category_id path string true "Default Category ID"
// @Success 201 {object} dto.ItemCategoryResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /item-categories/from-default/{default_id}/pantry/{pantry_id} [post]
func (h *itemCategoryHandler) CloneDefaultCategoryToPantry(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.CloneDefaultCategoryToPantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.CloneDefaultCategoryToPantry"), zap.Any("params", __logParams))
	pantryIDStr := c.Param("pantry_id")
	pantryID, err := uuid.Parse(pantryIDStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.CloneDefaultCategoryToPantry"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Pantry ID")
		return
	}

	defaultCategoryIDStr := c.Param("default_id")
	defaultCategoryID, err := uuid.Parse(defaultCategoryIDStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.CloneDefaultCategoryToPantry"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Default Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.CloneDefaultCategoryToPantry(c.Request.Context(), defaultCategoryID, pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.CloneDefaultCategoryToPantry"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Default category not found")
		case errors.Is(err, domain.ErrCategoryNotDefault):
			response.BadRequest(c, "Source category is not marked as default")
		default:
			response.InternalError(c, "Failed to clone item category")
		}
		return
	}

	response.Success(c, http.StatusCreated, category)
}

// @Summary Update an item category
// @Tags Item Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param body body dto.UpdateItemCategoryDTO true "Updated fields"
// @Success 200 {object} dto.ItemCategoryResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /item-categories/{id} [put]
func (h *itemCategoryHandler) UpdateItemCategory(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.UpdateItemCategory"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.UpdateItemCategory"), zap.Any("params", __logParams))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.UpdateItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	var input dto.UpdateItemCategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.UpdateItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input")
		return
	}

	category, err := h.service.Update(c.Request.Context(), id, input, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.UpdateItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item category not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to update item category")
		}
		return
	}

	response.OK(c, category)
}

// @Summary Get an item category by ID
// @Tags Item Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.ItemCategoryResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /item-categories/{id} [get]
func (h *itemCategoryHandler) GetItemCategory(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.GetItemCategory"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.GetItemCategory"), zap.Any("params", __logParams))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.GetItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.GetItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item category not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to fetch item category")
		}
		return
	}
	response.OK(c, category)
}

// @Summary Delete an item category by ID
// @Tags Item Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /item-categories/{id} [delete]
func (h *itemCategoryHandler) DeleteItemCategory(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.DeleteItemCategory"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.DeleteItemCategory"), zap.Any("params", __logParams))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.DeleteItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.DeleteItemCategory"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item category not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Only the creator can delete this category")
		default:
			response.InternalError(c, "Failed to delete item category")
		}
		return
	}

	response.OK(c, response.MessagePayload{Message: "Item category deleted successfully"})
}

// @Summary List item categories by pantry
// @Tags Item Categories
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {array} dto.ItemCategoryResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pantries/{id}/item-categories [get]
func (h *itemCategoryHandler) ListItemCategoriesByPantry(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.ListItemCategoriesByPantry"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.ListItemCategoriesByPantry"), zap.Any("params", __logParams))
	pantryIDStr := c.Param("id")
	pantryID, err := uuid.Parse(pantryIDStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.ListItemCategoriesByPantry"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid Pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	categories, err := h.service.ListByPantryID(c.Request.Context(), pantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.ListItemCategoriesByPantry"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			response.InternalError(c, "Failed to list item categories")
		}
		return
	}

	response.OK(c, categories)
}

// @Summary List item categories created by the user
// @Tags Item Categories
// @Produce json
// @Success 200 {array} dto.ItemCategoryResponse
// @Failure 500 {object} response.APIResponse
// @Router /item-categories/user [get]
func (h *itemCategoryHandler) ListItemCategoriesByUser(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*itemCategoryHandler.ListItemCategoriesByUser"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*itemCategoryHandler.ListItemCategoriesByUser"), zap.Any("params", __logParams))
	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	categories, err := h.service.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*itemCategoryHandler.ListItemCategoriesByUser"), zap.Error(err), zap.Any("params", __logParams))
		response.InternalError(c, "Failed to list item categories by user")
		return
	}

	response.OK(c, categories)
}
