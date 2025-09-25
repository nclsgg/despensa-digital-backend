package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/dto"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

type ShoppingListHandler struct {
	shoppingListService domain.ShoppingListService
}

func NewShoppingListHandler(shoppingListService domain.ShoppingListService) *ShoppingListHandler {
	return &ShoppingListHandler{
		shoppingListService: shoppingListService,
	}
}

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
func (h *ShoppingListHandler) CreateShoppingList(c *gin.Context) {
	var input dto.CreateShoppingListDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	shoppingList, err := h.shoppingListService.CreateShoppingList(c.Request.Context(), userUUID, input)
	if err != nil {
		response.InternalError(c, "Error creating shopping list: "+err.Error())
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
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	// Get pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.BadRequest(c, "Invalid offset parameter")
		return
	}

	shoppingLists, err := h.shoppingListService.GetShoppingListsByUserID(c.Request.Context(), userUUID, limit, offset)
	if err != nil {
		response.InternalError(c, "Error getting shopping lists: "+err.Error())
		return
	}

	response.OK(c, shoppingLists)
}

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
func (h *ShoppingListHandler) GetShoppingList(c *gin.Context) {
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	shoppingList, err := h.shoppingListService.GetShoppingListByID(c.Request.Context(), userUUID, shoppingListID)
	if err != nil {
		if err.Error() == "shopping list not found" {
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
			return
		}
		response.InternalError(c, "Error getting shopping list: "+err.Error())
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
	var input dto.UpdateShoppingListDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	shoppingList, err := h.shoppingListService.UpdateShoppingList(c.Request.Context(), userUUID, shoppingListID, input)
	if err != nil {
		if err.Error() == "shopping list not found" {
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
			return
		}
		response.InternalError(c, "Error updating shopping list: "+err.Error())
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
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	err = h.shoppingListService.DeleteShoppingList(c.Request.Context(), userUUID, shoppingListID)
	if err != nil {
		if err.Error() == "shopping list not found" {
			response.Fail(c, http.StatusNotFound, "SHOPPING_LIST_NOT_FOUND", "Shopping list not found")
			return
		}
		response.InternalError(c, "Error deleting shopping list: "+err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Shopping list deleted successfully"})
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
	var input dto.UpdateShoppingListItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	itemIdParam := c.Param("itemId")
	itemID, err := uuid.Parse(itemIdParam)
	if err != nil {
		response.BadRequest(c, "Invalid item ID")
		return
	}

	item, err := h.shoppingListService.UpdateShoppingListItem(c.Request.Context(), userUUID, shoppingListID, itemID, input)
	if err != nil {
		if err.Error() == "shopping list not found" || err.Error() == "item not found" {
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		response.InternalError(c, "Error updating shopping list item: "+err.Error())
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
	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	idParam := c.Param("id")
	shoppingListID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid shopping list ID")
		return
	}

	itemIdParam := c.Param("itemId")
	itemID, err := uuid.Parse(itemIdParam)
	if err != nil {
		response.BadRequest(c, "Invalid item ID")
		return
	}

	err = h.shoppingListService.DeleteShoppingListItem(c.Request.Context(), userUUID, shoppingListID, itemID)
	if err != nil {
		if err.Error() == "shopping list not found" || err.Error() == "item not found" {
			response.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		response.InternalError(c, "Error deleting shopping list item: "+err.Error())
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
	var input dto.GenerateAIShoppingListDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userUUID := rawID.(uuid.UUID)

	shoppingList, err := h.shoppingListService.GenerateAIShoppingList(c.Request.Context(), userUUID, input)
	if err != nil {
		response.InternalError(c, "Error generating AI shopping list: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, shoppingList)
}
