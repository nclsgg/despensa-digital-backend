package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

type itemCategoryHandler struct {
	service domain.ItemCategoryService
}

func NewItemCategoryHandler(service domain.ItemCategoryService) domain.ItemCategoryHandler {
	return &itemCategoryHandler{service}
}

func (h *itemCategoryHandler) CreateItemCategory(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	var input dto.CreateItemCategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.Create(c.Request.Context(), input, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidPantry):
			response.BadRequest(c, "Invalid pantry ID")
		case errors.Is(err, domain.ErrUnauthorized):
			logger.Warn("Access denied to pantry",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "CreateItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			logger.Error("Failed to create item category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "CreateItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to create item category")
		}
		return
	}

	logger.Info("Item category created successfully",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "CreateItemCategory"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", category.ID),
	)
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
	logger := appLogger.FromContext(c.Request.Context())

	var input dto.CreateDefaultItemCategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.CreateDefault(c.Request.Context(), input, userID)
	if err != nil {
		logger.Error("Failed to create default item category",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "CreateDefaultItemCategory"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to create default item category")
		return
	}

	logger.Info("Default item category created successfully",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "CreateDefaultItemCategory"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", category.ID),
	)
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
	logger := appLogger.FromContext(c.Request.Context())

	pantryIDStr := c.Param("pantry_id")
	pantryID, err := uuid.Parse(pantryIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid Pantry ID")
		return
	}

	defaultCategoryIDStr := c.Param("default_id")
	defaultCategoryID, err := uuid.Parse(defaultCategoryIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid Default Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.CloneDefaultCategoryToPantry(c.Request.Context(), defaultCategoryID, pantryID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			logger.Warn("Access denied to pantry",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
			)
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Default category not found")
		case errors.Is(err, domain.ErrCategoryNotDefault):
			response.BadRequest(c, "Source category is not marked as default")
		default:
			logger.Error("Failed to clone item category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
				zap.String("default_category_id", defaultCategoryID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to clone item category")
		}
		return
	}

	logger.Info("Default category cloned to pantry successfully",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "CloneDefaultCategoryToPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.String("category_id", category.ID),
	)
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
	logger := appLogger.FromContext(c.Request.Context())

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	var input dto.UpdateItemCategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input")
		return
	}

	category, err := h.service.Update(c.Request.Context(), id, input, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item category not found")
		case errors.Is(err, domain.ErrUnauthorized):
			logger.Warn("Access denied to update category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "UpdateItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("category_id", id.String()),
			)
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			logger.Error("Failed to update item category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "UpdateItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("category_id", id.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to update item category")
		}
		return
	}

	logger.Info("Item category updated successfully",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "UpdateItemCategory"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", id.String()),
	)
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
	logger := appLogger.FromContext(c.Request.Context())

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	category, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item category not found")
		case errors.Is(err, domain.ErrUnauthorized):
			logger.Warn("Access denied to view category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "GetItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("category_id", id.String()),
			)
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			logger.Error("Failed to fetch item category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "GetItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("category_id", id.String()),
				zap.Error(err),
			)
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
	logger := appLogger.FromContext(c.Request.Context())

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid Category ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		switch {
		case errors.Is(err, domain.ErrCategoryNotFound):
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Item category not found")
		case errors.Is(err, domain.ErrUnauthorized):
			logger.Warn("Unauthorized attempt to delete category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "DeleteItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("category_id", id.String()),
			)
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Only the creator can delete this category")
		default:
			logger.Error("Failed to delete item category",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "DeleteItemCategory"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("category_id", id.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to delete item category")
		}
		return
	}

	logger.Info("Item category deleted successfully",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "DeleteItemCategory"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("category_id", id.String()),
	)
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
	logger := appLogger.FromContext(c.Request.Context())

	pantryIDStr := c.Param("id")
	pantryID, err := uuid.Parse(pantryIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid Pantry ID")
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	categories, err := h.service.ListByPantryID(c.Request.Context(), pantryID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			logger.Warn("Access denied to list pantry categories",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "ListItemCategoriesByPantry"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
			)
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		default:
			logger.Error("Failed to list item categories",
				zap.String(appLogger.FieldModule, "item_category"),
				zap.String(appLogger.FieldFunction, "ListItemCategoriesByPantry"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.String("pantry_id", pantryID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to list item categories")
		}
		return
	}

	logger.Info("Listed item categories by pantry",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "ListItemCategoriesByPantry"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.String("pantry_id", pantryID.String()),
		zap.Int(appLogger.FieldCount, len(categories)),
	)
	response.OK(c, categories)
}

// @Summary List item categories created by the user
// @Tags Item Categories
// @Produce json
// @Success 200 {array} dto.ItemCategoryResponse
// @Failure 500 {object} response.APIResponse
// @Router /item-categories/user [get]
func (h *itemCategoryHandler) ListItemCategoriesByUser(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	categories, err := h.service.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to list item categories by user",
			zap.String(appLogger.FieldModule, "item_category"),
			zap.String(appLogger.FieldFunction, "ListItemCategoriesByUser"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to list item categories by user")
		return
	}

	logger.Info("Listed item categories by user",
		zap.String(appLogger.FieldModule, "item_category"),
		zap.String(appLogger.FieldFunction, "ListItemCategoriesByUser"),
		zap.String(appLogger.FieldUserID, userID.String()),
		zap.Int(appLogger.FieldCount, len(categories)),
	)
	response.OK(c, categories)
}
