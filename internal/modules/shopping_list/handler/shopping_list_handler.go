package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/dto"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

type ShoppingListHandler struct {
	shoppingListService domain.ShoppingListService
}

func NewShoppingListHandler(shoppingListService domain.ShoppingListService) (result0 *ShoppingListHandler) {
	__logParams := map[string]any{"shoppingListService": shoppingListService}
	__logStart := time.Now(

	// CreateShoppingList godoc
	// @Summary Create shopping list
	// @Description Create a new shopping list for the authenticated user
	// @Tags shopping-list
	// @Accept json
	// @Produce json
	// @Param shopping-list body dto.CreateShoppingListDTO true "Shopping list data"
	// @Success 201 {object} response.APIResponse{data=dto.ShoppingListResponseDTO}
	// @Failure 400 {object} response.APIResponse
	// @Failure 401 {object} response.APIResponse
	// @Failure 500 {object} response.APIResponse
	// @Router /shopping-lists [post]
	// @Security BearerAuth
	)
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewShoppingListHandler"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewShoppingListHandler"), zap.Any("params", __logParams))
	result0 = &ShoppingListHandler{
		shoppingListService: shoppingListService,
	}
	return
}

func (h *ShoppingListHandler) CreateShoppingList(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.CreateShoppingList"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.CreateShoppingList"), zap.Any("params", __logParams))
	var input dto.CreateShoppingListDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.CreateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	shoppingList, err := h.shoppingListService.CreateShoppingList(c.Request.Context(), userUUID, input)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.CreateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrPantryAccessDenied):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		case errors.Is(err, domain.ErrPantryNotFound):
			response.Fail(c, http.StatusNotFound, "PANTRY_NOT_FOUND", "Pantry not found")
		default:
			response.InternalError(c, "Failed to create shopping list")
		}
		return
	}

	response.Success(c, http.StatusCreated, shoppingList)
}

// GetShoppingLists godoc
// @Summary Get shopping lists
// @Description Get all shopping lists for the authenticated user
// @Tags shopping-list
// @Produce json
// @Param limit query int false "Limit number of results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} response.APIResponse{data=[]dto.ShoppingListSummaryDTO}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /shopping-lists [get]
// @Security BearerAuth
func (h *ShoppingListHandler) GetShoppingLists(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func

	// Get pagination parameters
	() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.GetShoppingLists"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.GetShoppingLists"), zap.Any("params", __logParams))
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.GetShoppingLists"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.GetShoppingLists"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid offset parameter")
		return
	}

	shoppingLists, err := h.shoppingListService.GetShoppingListsByUserID(c.Request.Context(), userUUID, limit, offset)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.GetShoppingLists"), zap.Error(err), zap.Any("params",

			// GetShoppingList godoc
			// @Summary Get shopping list by ID
			// @Description Get a specific shopping list by ID
			// @Tags shopping-list
			// @Produce json
			// @Param id path string true "Shopping list ID"
			// @Success 200 {object} response.APIResponse{data=dto.ShoppingListResponseDTO}
			// @Failure 400 {object} response.APIResponse
			// @Failure 401 {object} response.APIResponse
			// @Failure 404 {object} response.APIResponse
			// @Failure 500 {object} response.APIResponse
			// @Router /shopping-lists/{id} [get]
			// @Security BearerAuth
			__logParams))
		response.InternalError(c, "Failed to fetch shopping lists")
		return
	}

	response.OK(c, shoppingLists)
}

func (h *ShoppingListHandler) GetShoppingList(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.GetShoppingList"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.GetShoppingList"), zap.Any("params", __logParams))
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.GetShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	shoppingList, err := h.shoppingListService.GetShoppingListByID(c.Request.Context(), userUUID, shoppingListID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.GetShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrShoppingListNotFound):
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this shopping list")
		default:
			response.InternalError(c, "Failed to fetch shopping list")
		}
		return
	}

	response.OK(c, shoppingList)
}

// UpdateShoppingList godoc
// @Summary Update shopping list
// @Description Update a shopping list
// @Tags shopping-list
// @Accept json
// @Produce json
// @Param id path string true "Shopping list ID"
// @Param shopping-list body dto.UpdateShoppingListDTO true "Shopping list update data"
// @Success 200 {object} response.APIResponse{data=dto.ShoppingListResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /shopping-lists/{id} [put]
// @Security BearerAuth
func (h *ShoppingListHandler) UpdateShoppingList(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.UpdateShoppingList"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.UpdateShoppingList"), zap.Any("params", __logParams))
	var input dto.UpdateShoppingListDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.UpdateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.UpdateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	shoppingList, err := h.shoppingListService.UpdateShoppingList(c.Request.Context(), userUUID, shoppingListID, input)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.UpdateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrShoppingListNotFound):
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this shopping list")
		default:
			response.InternalError(c, "Failed to update shopping list")
		}
		return
	}

	response.OK(c, shoppingList)
}

// DeleteShoppingList godoc
// @Summary Delete shopping list
// @Description Delete a shopping list
// @Tags shopping-list
// @Produce json
// @Param id path string true "Shopping list ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /shopping-lists/{id} [delete]
// @Security BearerAuth
func (h *ShoppingListHandler) DeleteShoppingList(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.DeleteShoppingList"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.DeleteShoppingList"), zap.Any("params", __logParams))
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.DeleteShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	err = h.shoppingListService.DeleteShoppingList(c.Request.Context(), userUUID, shoppingListID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.DeleteShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrShoppingListNotFound):
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this shopping list")
		default:
			response.InternalError(c, "Failed to delete shopping list")
		}
		return
	}

	response.OK(c, gin.H{"message": "Shopping list deleted successfully"})
}

// CreateShoppingListItem godoc
// @Summary Create shopping list item
// @Description Add a new item to an existing shopping list
// @Tags shopping-list
// @Accept json
// @Produce json
// @Param id path string true "Shopping list ID"
// @Param item body dto.CreateShoppingListItemDTO true "Shopping list item data"
// @Success 201 {object} response.APIResponse{data=dto.ShoppingListResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /shopping-lists/{id}/items [post]
// @Security BearerAuth
func (h *ShoppingListHandler) CreateShoppingListItem(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.CreateShoppingListItem"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.CreateShoppingListItem"), zap.Any("params", __logParams))

	var input dto.CreateShoppingListItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.CreateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.CreateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	shoppingList, err := h.shoppingListService.CreateShoppingListItem(c.Request.Context(), userUUID, shoppingListID, input)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.CreateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrShoppingListNotFound):
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this shopping list")
		default:
			response.InternalError(c, "Failed to create shopping list item")
		}
		return
	}

	response.Success(c, http.StatusCreated, shoppingList)
}

// UpdateShoppingListItem godoc
// @Summary Update shopping list item
// @Description Update an item in a shopping list
// @Tags shopping-list
// @Accept json
// @Produce json
// @Param id path string true "Shopping list ID"
// @Param itemId path string true "Shopping list item ID"
// @Param item body dto.UpdateShoppingListItemDTO true "Shopping list item update data"
// @Success 200 {object} response.APIResponse{data=dto.ShoppingListItemResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /shopping-lists/{id}/items/{itemId} [put]
// @Security BearerAuth
func (h *ShoppingListHandler) UpdateShoppingListItem(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.UpdateShoppingListItem"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.UpdateShoppingListItem"), zap.Any("params", __logParams))
	var input dto.UpdateShoppingListItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.UpdateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.UpdateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	itemIdParam := c.Param("itemId")
	itemID, err := uuid.Parse(itemIdParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.UpdateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid item ID")
		return
	}

	item, err := h.shoppingListService.UpdateShoppingListItem(c.Request.Context(), userUUID, shoppingListID, itemID, input)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.UpdateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrShoppingListNotFound):
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
		case errors.Is(err, domain.ErrItemNotFound):
			response.Fail(c, http.StatusNotFound, "ITEM_NOT_FOUND", "Item not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this shopping list")
		default:
			response.InternalError(c, "Failed to update shopping list item")
		}
		return
	}

	response.OK(c, item)
}

// DeleteShoppingListItem godoc
// @Summary Delete shopping list item
// @Description Delete an item from a shopping list
// @Tags shopping-list
// @Produce json
// @Param id path string true "Shopping list ID"
// @Param itemId path string true "Shopping list item ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /shopping-lists/{id}/items/{itemId} [delete]
// @Security BearerAuth
func (h *ShoppingListHandler) DeleteShoppingListItem(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.DeleteShoppingListItem"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.DeleteShoppingListItem"), zap.Any("params", __logParams))
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.DeleteShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	itemIdParam := c.Param("itemId")
	itemID, err := uuid.Parse(itemIdParam)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.DeleteShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid item ID")
		return
	}

	err = h.shoppingListService.DeleteShoppingListItem(c.Request.Context(), userUUID, shoppingListID, itemID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.DeleteShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrShoppingListNotFound):
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
		case errors.Is(err, domain.ErrItemNotFound):
			response.Fail(c, http.StatusNotFound, "ITEM_NOT_FOUND", "Item not found")
		case errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this shopping list")
		default:
			response.InternalError(c, "Failed to delete shopping list item")
		}
		return
	}

	response.OK(c, gin.H{"message": "Shopping list item deleted successfully"})
}

// GenerateAIShoppingList godoc
// @Summary Generate AI shopping list
// @Description Generate a shopping list using AI based on user profile and specific pantry history
// @Tags shopping-list
// @Accept json
// @Produce json
// @Param shopping-list body dto.GenerateAIShoppingListDTO true "AI shopping list generation data (requires pantry_id)"
// @Success 201 {object} response.APIResponse{data=dto.ShoppingListResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /shopping-lists/generate [post]
// @Security BearerAuth
func (h *ShoppingListHandler) GenerateAIShoppingList(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListHandler.GenerateAIShoppingList"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListHandler.GenerateAIShoppingList"), zap.Any("params", __logParams))
	var input dto.GenerateAIShoppingListDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	shoppingList, err := h.shoppingListService.GenerateAIShoppingList(c.Request.Context(), userUUID, input)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ShoppingListHandler.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrPantryNotFound):
			response.Fail(c, http.StatusNotFound, "PANTRY_NOT_FOUND", "Pantry not found")
		case errors.Is(err, domain.ErrPantryAccessDenied), errors.Is(err, domain.ErrUnauthorized):
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "Access denied to this pantry")
		case errors.Is(err, domain.ErrPromptBuildFailed):
			response.Fail(c, http.StatusUnprocessableEntity, "PROMPT_BUILD_FAILED", "Unable to build AI prompt")
		case errors.Is(err, domain.ErrAIRequestFailed):
			response.Fail(c, http.StatusBadGateway, "AI_REQUEST_FAILED", "AI provider unavailable")
		case errors.Is(err, domain.ErrAIResponseInvalid):
			response.Fail(c, http.StatusBadGateway, "AI_RESPONSE_INVALID", "AI returned an invalid response")
		case errors.Is(err, domain.ErrShoppingListNotFound):
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
		default:
			response.InternalError(c, "Failed to generate AI shopping list")
		}
		return
	}

	response.Success(c, http.StatusCreated, shoppingList)
}
