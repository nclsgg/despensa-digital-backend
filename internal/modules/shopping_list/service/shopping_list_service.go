package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	itemDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	itemModel "github.com/nclsgg/despensa-digital/backend/internal/modules/item/model"
	llmDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/domain"
	pantryDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	profileDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	profileModel "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/dto"
	shoppingModel "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PantryInsights struct {
	CommonItems       []ItemInsight      `json:"common_items"`
	LowStockItems     []ItemInsight      `json:"low_stock_items"`
	ExpiringSoonItems []ItemInsight      `json:"expiring_soon_items"`
	AverageItemPrice  map[string]float64 `json:"average_item_price"`
	Categories        []string           `json:"categories"`
	TotalItems        int                `json:"total_items"`
}

type ItemInsight struct {
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	Frequency       int       `json:"frequency"`
	AveragePrice    float64   `json:"average_price"`
	LastPurchased   time.Time `json:"last_purchased"`
	QuantityPattern float64   `json:"quantity_pattern"`
	Unit            string    `json:"unit"`
}

type AIShoppingItem struct {
	Name           string  `json:"name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	EstimatedPrice float64 `json:"estimated_price"`
	Category       string  `json:"category"`
	Priority       int     `json:"priority"`
	Reason         string  `json:"reason"`
}

type AIShoppingListResponse struct {
	Items          []AIShoppingItem `json:"items"`
	Reasoning      string           `json:"reasoning"`
	EstimatedTotal float64          `json:"estimated_total"`
}

type shoppingPreferences struct {
	HouseholdSize       int
	MonthlyIncome       float64
	DietaryRestrictions []string
}

type shoppingListService struct {
	shoppingListRepo domain.ShoppingListRepository
	pantryRepo       pantryDomain.PantryRepository
	itemRepo         itemDomain.ItemRepository
	profileRepo      profileDomain.ProfileRepository
	llmService       llmDomain.LLMService
}

func NewShoppingListService(
	shoppingListRepo domain.ShoppingListRepository,
	pantryRepo pantryDomain.PantryRepository,
	itemRepo itemDomain.ItemRepository,
	profileRepo profileDomain.ProfileRepository,
	llmService llmDomain.LLMService,
) (result0 domain.ShoppingListService) {
	__logParams := map[string]any{"shoppingListRepo": shoppingListRepo, "pantryRepo": pantryRepo, "itemRepo": itemRepo, "profileRepo": profileRepo, "llmService": llmService}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewShoppingListService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewShoppingListService"), zap.Any("params", __logParams))
	result0 = &shoppingListService{
		shoppingListRepo: shoppingListRepo,
		pantryRepo:       pantryRepo,
		itemRepo:         itemRepo,
		profileRepo:      profileRepo,
		llmService:       llmService,
	}
	return
}

func (s *shoppingListService) CreateShoppingList(ctx context.Context, userID uuid.UUID, input dto.CreateShoppingListDTO) (result0 *dto.ShoppingListResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.CreateShoppingList"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.CreateShoppingList"), zap.Any("params", __logParams))

	preferences, err := s.resolvePreferences(ctx, userID, input.Preferences)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	shoppingList := &shoppingModel.ShoppingList{
		UserID:              userID,
		PantryID:            input.PantryID,
		Name:                input.Name,
		TotalBudget:         input.TotalBudget,
		Status:              "pending",
		GeneratedBy:         "manual",
		HouseholdSize:       preferences.HouseholdSize,
		MonthlyIncome:       preferences.MonthlyIncome,
		DietaryRestrictions: shoppingModel.StringArray(normalizeStringSlice(preferences.DietaryRestrictions)),
	}

	if input.PantryID != nil {
		hasAccess, err := s.pantryRepo.IsUserInPantry(ctx, *input.PantryID, userID)
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
			result0 = nil
			result1 = err
			return
		}
		if !hasAccess {
			result0 = nil
			result1 = domain.ErrPantryAccessDenied
			return
		}
	}

	for _, itemDto := range input.Items {
		priceQuantity := itemDto.PriceQuantity
		if priceQuantity <= 0 {
			priceQuantity = 1
		}
		item := &shoppingModel.ShoppingListItem{
			Name:           itemDto.Name,
			Quantity:       itemDto.Quantity,
			Unit:           itemDto.Unit,
			EstimatedPrice: itemDto.EstimatedPrice,
			PriceQuantity:  priceQuantity,
			Category:       itemDto.Category,
			Priority:       itemDto.Priority,
			Notes:          itemDto.Notes,
			Source:         "manual",
		}
		if item.Priority == 0 {
			item.Priority = 3
		}
		if itemDto.PantryItemID != nil {
			item.PantryItemID = itemDto.PantryItemID
		}
		shoppingList.Items = append(shoppingList.Items, *item)
	}

	estimatedTotal, actualTotal := calculateListTotals(shoppingList.Items)
	shoppingList.EstimatedCost = estimatedTotal
	shoppingList.ActualCost = actualTotal

	if err := s.shoppingListRepo.Create(ctx, shoppingList); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("create shopping list: %w", err)
		return
	}

	created, err := s.shoppingListRepo.GetByID(ctx, shoppingList.ID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrShoppingListNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("reload shopping list: %w", err)
		return
	}
	result0 = s.convertToResponseDTO(ctx, created)
	result1 = nil
	return
}

func (s *shoppingListService) GetShoppingListByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (result0 *dto.ShoppingListResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.GetShoppingListByID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.GetShoppingListByID"), zap.Any("params", __logParams))
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GetShoppingListByID"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrShoppingListNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get shopping list: %w", err)
		return
	}

	if shoppingList.UserID != userID {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}
	result0 = s.convertToResponseDTO(ctx, shoppingList)
	result1 = nil
	return
}

func (s *shoppingListService) GetShoppingListsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) (result0 []*dto.ShoppingListSummaryDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "limit": limit, "offset": offset}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.GetShoppingListsByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.GetShoppingListsByUserID"), zap.Any("params", __logParams))
	shoppingLists, err := s.shoppingListRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GetShoppingListsByUserID"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("error getting shopping lists: %w", err)
		return
	}

	pantryNames := s.resolvePantryNames(ctx, shoppingLists)

	summaries := make([]*dto.ShoppingListSummaryDTO, 0, len(shoppingLists))
	for _, sl := range shoppingLists {
		itemCount := len(sl.Items)
		purchasedCount := 0
		for _, item := range sl.Items {
			if item.Purchased {
				purchasedCount++
			}
		}

		var pantryID *string
		pantryName := ""
		if sl.PantryID != nil {
			idStr := sl.PantryID.String()
			pantryID = &idStr
			if name, ok := pantryNames[*sl.PantryID]; ok {
				pantryName = name
			}
		}

		summaries = append(summaries, &dto.ShoppingListSummaryDTO{
			ID:             sl.ID.String(),
			PantryID:       pantryID,
			PantryName:     pantryName,
			Name:           sl.Name,
			Status:         sl.Status,
			TotalBudget:    sl.TotalBudget,
			EstimatedCost:  sl.EstimatedCost,
			ActualCost:     sl.ActualCost,
			GeneratedBy:    sl.GeneratedBy,
			ItemCount:      itemCount,
			PurchasedCount: purchasedCount,
			Preferences:    convertPreferencesToDTO(sl),
			CreatedAt:      sl.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      sl.UpdatedAt.Format(time.RFC3339),
		})
	}
	result0 = summaries
	result1 = nil
	return
}

func (s *shoppingListService) UpdateShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID, input dto.UpdateShoppingListDTO) (result0 *dto.ShoppingListResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "id": id, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.UpdateShoppingList"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.UpdateShoppingList"), zap.Any("params", __logParams))
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.UpdateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrShoppingListNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get shopping list: %w", err)
		return
	}

	if shoppingList.UserID != userID {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	if input.Name != nil {
		shoppingList.Name = *input.Name
	}
	if input.TotalBudget != nil {
		shoppingList.TotalBudget = *input.TotalBudget
	}

	if input.Preferences != nil {
		prefs := shoppingPreferences{
			HouseholdSize:       shoppingList.HouseholdSize,
			MonthlyIncome:       shoppingList.MonthlyIncome,
			DietaryRestrictions: toShoppingListStringSlice(shoppingList.DietaryRestrictions),
		}
		prefs = applyPreferenceOverrides(prefs, input.Preferences)
		shoppingList.HouseholdSize = prefs.HouseholdSize
		shoppingList.MonthlyIncome = prefs.MonthlyIncome
		shoppingList.DietaryRestrictions = shoppingModel.StringArray(normalizeStringSlice(prefs.DietaryRestrictions))
	}

	checkoutPerformed := false
	checkoutCost := 0.0
	if input.Status != nil {
		targetStatus := *input.Status
		if targetStatus == "completed" && shoppingList.Status != "completed" {
			cost, err := s.performCheckout(ctx, userID, shoppingList)
			if err != nil {
				zap.L().Error("function.error", zap.String("func", "*shoppingListService.UpdateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
				result0 = nil
				result1 = err
				return
			}
			checkoutPerformed = true
			checkoutCost = cost
		}
		shoppingList.Status = targetStatus
	}

	if input.ActualCost != nil {
		shoppingList.ActualCost = *input.ActualCost
	} else if checkoutPerformed {
		shoppingList.ActualCost = checkoutCost
	}

	if err := s.shoppingListRepo.Update(ctx, shoppingList); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.UpdateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("update shopping list: %w", err)
		return
	}

	updated, err := s.shoppingListRepo.GetByID(ctx, shoppingList.ID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.UpdateShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrShoppingListNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("reload shopping list: %w", err)
		return
	}
	result0 = s.convertToResponseDTO(ctx, updated)
	result1 = nil
	return
}

func (s *shoppingListService) DeleteShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.DeleteShoppingList"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.DeleteShoppingList"), zap.Any("params", __logParams))
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.DeleteShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = domain.ErrShoppingListNotFound
			return
		}
		result0 = fmt.Errorf("get shopping list: %w", err)
		return
	}

	if shoppingList.UserID != userID {
		result0 = domain.ErrUnauthorized
		return
	}

	if err := s.shoppingListRepo.Delete(ctx, id); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.DeleteShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = fmt.Errorf("delete shopping list: %w", err)
		return
	}
	result0 = nil
	return
}

func (s *shoppingListService) CreateShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, input dto.CreateShoppingListItemDTO) (result0 *dto.ShoppingListResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "shoppingListID": shoppingListID, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.CreateShoppingListItem"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.CreateShoppingListItem"), zap.Any("params", __logParams))

	// Get shopping list and verify ownership
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrShoppingListNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get shopping list: %w", err)
		return
	}

	if shoppingList.UserID != userID {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	// Create new item with default values
	priceQuantity := input.PriceQuantity
	if priceQuantity <= 0 {
		priceQuantity = 1
	}

	priority := input.Priority
	if priority == 0 {
		priority = 3
	}

	newItem := &shoppingModel.ShoppingListItem{
		ShoppingListID: shoppingListID,
		Name:           input.Name,
		Quantity:       input.Quantity,
		Unit:           input.Unit,
		EstimatedPrice: input.EstimatedPrice,
		PriceQuantity:  priceQuantity,
		Category:       input.Category,
		Priority:       priority,
		Notes:          input.Notes,
		Source:         "manual",
	}

	if input.PantryItemID != nil {
		newItem.PantryItemID = input.PantryItemID
	}

	// Save the new item
	if err := s.shoppingListRepo.CreateItem(ctx, newItem); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("create shopping list item: %w", err)
		return
	}

	// Reload shopping list to get all items including the new one
	shoppingList, err = s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("reload shopping list: %w", err)
		return
	}

	// Recalculate totals
	estimatedTotal, actualTotal := calculateListTotals(shoppingList.Items)
	shoppingList.EstimatedCost = estimatedTotal
	shoppingList.ActualCost = actualTotal

	// Update shopping list with new totals
	if err := s.shoppingListRepo.Update(ctx, shoppingList); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.CreateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("update shopping list totals: %w", err)
		return
	}

	// Return the full shopping list with the new item
	result0 = s.convertToResponseDTO(ctx, shoppingList)
	result1 = nil
	return
}

func (s *shoppingListService) UpdateShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID, input dto.UpdateShoppingListItemDTO) (result0 *dto.ShoppingListItemResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "shoppingListID": shoppingListID, "itemID": itemID, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.UpdateShoppingListItem"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.UpdateShoppingListItem"), zap.Any("params", __logParams))
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.UpdateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrShoppingListNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get shopping list: %w", err)
		return
	}

	if shoppingList.UserID != userID {
		result0 = nil
		result1 = domain.ErrUnauthorized
		return
	}

	var (
		targetItem  *shoppingModel.ShoppingListItem
		targetIndex = -1
	)
	for idx := range shoppingList.Items {
		if shoppingList.Items[idx].ID == itemID {
			targetItem = &shoppingList.Items[idx]
			targetIndex = idx
			break
		}
	}

	if targetItem == nil {
		result0 = nil
		result1 = domain.ErrItemNotFound
		return
	}

	if targetItem.PriceQuantity <= 0 {
		targetItem.PriceQuantity = 1
	}

	if input.Name != nil {
		targetItem.Name = *input.Name
	}
	if input.Quantity != nil {
		targetItem.Quantity = *input.Quantity
	}
	if input.Unit != nil {
		targetItem.Unit = *input.Unit
	}
	if input.EstimatedPrice != nil {
		targetItem.EstimatedPrice = *input.EstimatedPrice
	}
	if input.PriceQuantity != nil {
		priceQuantity := *input.PriceQuantity
		if priceQuantity <= 0 {
			priceQuantity = 1
		}
		targetItem.PriceQuantity = priceQuantity
	}
	if input.ActualPrice != nil {
		targetItem.ActualPrice = *input.ActualPrice
	}
	if input.Category != nil {
		targetItem.Category = *input.Category
	}
	if input.Priority != nil {
		targetItem.Priority = *input.Priority
	}
	if input.Purchased != nil {
		targetItem.Purchased = *input.Purchased
	}
	if input.Notes != nil {
		targetItem.Notes = *input.Notes
	}
	if input.PantryItemID != nil {
		targetItem.PantryItemID = input.PantryItemID
	}

	if err := s.shoppingListRepo.UpdateItem(ctx, targetItem); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.UpdateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("update shopping list item: %w", err)
		return
	}

	if targetIndex >= 0 {
		shoppingList.Items[targetIndex] = *targetItem
	}

	estimatedTotal, actualTotal := calculateListTotals(shoppingList.Items)
	shoppingList.EstimatedCost = estimatedTotal
	shoppingList.ActualCost = actualTotal

	if err := s.shoppingListRepo.Update(ctx, shoppingList); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.UpdateShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("update shopping list totals: %w", err)
		return
	}

	reloadedItems, err := s.shoppingListRepo.GetItemsByShoppingListID(ctx, shoppingListID)
	if err == nil {
		for _, item := range reloadedItems {
			if item.ID == itemID {
				result0 = s.convertItemToResponseDTO(item)
				result1 = nil
				return
			}
		}
	}
	result0 = s.convertItemToResponseDTO(targetItem)
	result1 = nil
	return
}

func (s *shoppingListService) DeleteShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "shoppingListID": shoppingListID, "itemID": itemID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.DeleteShoppingListItem"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.DeleteShoppingListItem"), zap.Any("params", __logParams))
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.DeleteShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = domain.ErrShoppingListNotFound
			return
		}
		result0 = fmt.Errorf("get shopping list: %w", err)
		return
	}

	if shoppingList.UserID != userID {
		result0 = domain.ErrUnauthorized
		return
	}

	found := false
	for _, item := range shoppingList.Items {
		if item.ID == itemID {
			found = true
			break
		}
	}

	if !found {
		result0 = domain.ErrItemNotFound
		return
	}

	if err := s.shoppingListRepo.DeleteItem(ctx, itemID); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.DeleteShoppingListItem"), zap.Error(err), zap.Any("params", __logParams))
		result0 = fmt.Errorf("delete shopping list item: %w", err)
		return
	}
	result0 = nil
	return
}

func (s *shoppingListService) GenerateAIShoppingList(ctx context.Context, userID uuid.UUID, input dto.GenerateAIShoppingListDTO) (result0 *dto.ShoppingListResponseDTO, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Any("params", __logParams))
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		result0 = nil
		result1 = fmt.Errorf("get user profile: %w", err)
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		profile = nil
	}

	pantry, err := s.pantryRepo.GetByID(ctx, input.PantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrPantryNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get pantry: %w", err)
		return
	}

	hasAccess, err := s.pantryRepo.IsUserInPantry(ctx, input.PantryID, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	if !hasAccess {
		result0 = nil
		result1 = domain.ErrPantryAccessDenied
		return
	}

	pantryInsights, err := s.analyzePantryHistory(ctx, []*pantryModel.Pantry{pantry})
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("analyze pantry history: %w", err)
		return
	}

	budget := s.determineBudget(input, profile)
	includeBasics := true
	if input.IncludeBasics != nil {
		includeBasics = *input.IncludeBasics
	}

	preferences, err := s.resolvePreferences(ctx, userID, input.Preferences)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	prompt, err := s.buildShoppingListPrompt(input, preferences, profile, pantryInsights, budget, includeBasics)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("%w: %v", domain.ErrPromptBuildFailed, err)
		return
	}

	llmResponse, err := s.llmService.GenerateText(ctx, prompt, map[string]interface{}{
		"max_tokens":      2000,
		"temperature":     0.7,
		"response_format": "json",
	})
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("%w: %v", domain.ErrAIRequestFailed, err)
		return
	}

	shoppingList, err := s.parseAIResponse(userID, input, budget, preferences, llmResponse.Response)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	if err := s.shoppingListRepo.Create(ctx, shoppingList); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("create ai shopping list: %w", err)
		return
	}

	created, err := s.shoppingListRepo.GetByID(ctx, shoppingList.ID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.GenerateAIShoppingList"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrShoppingListNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("reload shopping list: %w", err)
		return
	}
	result0 = s.convertToResponseDTO(ctx, created)
	result1 = nil
	return
}

// Helper methods

func (s *shoppingListService) resolvePreferences(ctx context.Context, userID uuid.UUID, overrides *dto.ShoppingListPreferencesOverrideDTO) (result0 shoppingPreferences, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "overrides": overrides}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.resolvePreferences"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.resolvePreferences"), zap.Any("params", __logParams))

	prefs := shoppingPreferences{
		HouseholdSize:       1,
		MonthlyIncome:       0,
		DietaryRestrictions: []string{},
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = applyPreferenceOverrides(prefs, overrides)
			result1 = nil
			return
		}
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.resolvePreferences"), zap.Error(err), zap.Any("params", __logParams))
		result0 = shoppingPreferences{}
		result1 = fmt.Errorf("get profile: %w", err)
		return
	}

	if profile != nil {
		prefs.HouseholdSize = profile.HouseholdSize
		prefs.MonthlyIncome = profile.MonthlyIncome
		prefs.DietaryRestrictions = toStringSlice(profile.DietaryRestrictions)
	}

	result0 = applyPreferenceOverrides(prefs, overrides)
	result1 = nil
	return
}

func applyPreferenceOverrides(base shoppingPreferences, overrides *dto.ShoppingListPreferencesOverrideDTO) (result0 shoppingPreferences) {
	__logParams := map[string]any{"base": base, "overrides": overrides}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "applyPreferenceOverrides"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "applyPreferenceOverrides"), zap.Any("params", __logParams))
	result0 = base
	if overrides == nil {
		return
	}
	if overrides.HouseholdSize != nil {
		result0.HouseholdSize = *overrides.HouseholdSize
	}
	if overrides.MonthlyIncome != nil {
		result0.MonthlyIncome = *overrides.MonthlyIncome
	}
	if overrides.DietaryRestrictions != nil {
		result0.DietaryRestrictions = normalizeStringSlice(overrides.DietaryRestrictions)
	}
	return
}

func convertPreferencesToDTO(sl *shoppingModel.ShoppingList) (result0 dto.ShoppingListPreferencesDTO) {
	__logParams := map[string]any{"sl": sl}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "convertPreferencesToDTO"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "convertPreferencesToDTO"), zap.Any("params", __logParams))
	result0 = dto.ShoppingListPreferencesDTO{
		HouseholdSize:       sl.HouseholdSize,
		MonthlyIncome:       sl.MonthlyIncome,
		DietaryRestrictions: toShoppingListStringSlice(sl.DietaryRestrictions),
	}
	return
}

func (s *shoppingListService) performCheckout(ctx context.Context, userID uuid.UUID, sl *shoppingModel.ShoppingList) (result0 float64, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "sl": sl}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.performCheckout"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.performCheckout"), zap.Any("params", __logParams))

	actualCost := 0.0

	normalizeName := func(name string) string {
		return strings.ToLower(strings.TrimSpace(name))
	}

	var pantryItemsByID map[uuid.UUID]*itemModel.Item
	var pantryItemsByName map[string]*itemModel.Item

	if sl.PantryID != nil && s.itemRepo != nil {
		items, err := s.itemRepo.ListByPantryID(ctx, *sl.PantryID)
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "*shoppingListService.performCheckout"), zap.Error(err), zap.Any("params", __logParams))
			result0 = 0
			result1 = fmt.Errorf("list pantry items: %w", err)
			return
		}
		pantryItemsByID = make(map[uuid.UUID]*itemModel.Item, len(items))
		pantryItemsByName = make(map[string]*itemModel.Item, len(items))
		for _, pantryItem := range items {
			pantryItemsByID[pantryItem.ID] = pantryItem
			pantryItemsByName[normalizeName(pantryItem.Name)] = pantryItem
		}
	}

	for idx := range sl.Items {
		item := &sl.Items[idx]
		if !item.Purchased {
			continue
		}

		priceQuantity := normalizePriceQuantity(item.PriceQuantity)
		basePrice := resolveUnitPrice(item.ActualPrice, item.EstimatedPrice)
		quantityFactor := item.Quantity / priceQuantity
		if quantityFactor < 0 {
			quantityFactor = 0
		}
		item.ActualPrice = basePrice
		actualCost += basePrice * quantityFactor

		perUnitPrice := 0.0
		if priceQuantity > 0 {
			perUnitPrice = basePrice / priceQuantity
		}

		var matchedPantryItem *itemModel.Item
		if sl.PantryID != nil && s.itemRepo != nil {
			if item.PantryItemID != nil {
				if pantryItemsByID != nil {
					if cached, ok := pantryItemsByID[*item.PantryItemID]; ok {
						matchedPantryItem = cached
					}
				}
				if matchedPantryItem == nil {
					found, err := s.itemRepo.FindByID(ctx, *item.PantryItemID)
					if err == nil {
						matchedPantryItem = found
						if pantryItemsByID != nil {
							pantryItemsByID[found.ID] = found
							pantryItemsByName[normalizeName(found.Name)] = found
						}
					}
				}
				if matchedPantryItem != nil && matchedPantryItem.ID == uuid.Nil {
					matchedPantryItem = nil
				}
			} else if pantryItemsByName != nil {
				if candidate, ok := pantryItemsByName[normalizeName(item.Name)]; ok {
					matchedPantryItem = candidate
					copied := candidate.ID
					item.PantryItemID = &copied
				}
			}
			if matchedPantryItem != nil && matchedPantryItem.ID == uuid.Nil {
				matchedPantryItem = nil
			}

			if matchedPantryItem != nil {
				matchedPantryItem.Quantity += item.Quantity
				if perUnitPrice > 0 {
					matchedPantryItem.PricePerUnit = perUnitPrice
					matchedPantryItem.PriceQuantity = priceQuantity
				}
				if matchedPantryItem.Unit == "" {
					matchedPantryItem.Unit = item.Unit
				}
				if err := s.itemRepo.Update(ctx, matchedPantryItem); err != nil {
					zap.L().Error("function.error", zap.String("func", "*shoppingListService.performCheckout"), zap.Error(err), zap.Any("params", __logParams))
					result0 = 0
					result1 = fmt.Errorf("update pantry item: %w", err)
					return
				}
				if pantryItemsByID != nil {
					pantryItemsByID[matchedPantryItem.ID] = matchedPantryItem
					pantryItemsByName[normalizeName(matchedPantryItem.Name)] = matchedPantryItem
				}
			} else if sl.PantryID != nil {
				newItem := &itemModel.Item{
					ID:            uuid.New(),
					PantryID:      *sl.PantryID,
					AddedBy:       userID,
					Name:          item.Name,
					Quantity:      item.Quantity,
					PricePerUnit:  perUnitPrice,
					PriceQuantity: priceQuantity,
					Unit:          item.Unit,
				}
				if err := s.itemRepo.Create(ctx, newItem); err != nil {
					zap.L().Error("function.error", zap.String("func", "*shoppingListService.performCheckout"), zap.Error(err), zap.Any("params", __logParams))
					result0 = 0
					result1 = fmt.Errorf("create pantry item: %w", err)
					return
				}
				copied := newItem.ID
				item.PantryItemID = &copied
				if pantryItemsByID != nil {
					pantryItemsByID[newItem.ID] = newItem
					pantryItemsByName[normalizeName(newItem.Name)] = newItem
				}
			}
		}

		if err := s.shoppingListRepo.UpdateItem(ctx, item); err != nil {
			zap.L().Error("function.error", zap.String("func", "*shoppingListService.performCheckout"), zap.Error(err), zap.Any("params", __logParams))
			result0 = 0
			result1 = fmt.Errorf("update shopping list item: %w", err)
			return
		}
	}

	result0 = actualCost
	result1 = nil
	return
}

func resolveUnitPrice(actualPrice, estimatedPrice float64) (result0 float64) {
	__logParams := map[string]any{"actualPrice": actualPrice, "estimatedPrice": estimatedPrice}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "resolveUnitPrice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "resolveUnitPrice"), zap.Any("params", __logParams))
	price := actualPrice
	if price <= 0 {
		price = estimatedPrice
	}
	if price < 0 {
		price = 0
	}
	result0 = price
	return
}

func normalizePriceQuantity(value float64) (result0 float64) {
	__logParams := map[string]any{"value": value}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "normalizePriceQuantity"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "normalizePriceQuantity"), zap.Any("params", __logParams))
	if value <= 0 {
		result0 = 1
		return
	}
	result0 = value
	return
}

func clampNonNegative(value float64) (result0 float64) {
	if value < 0 {
		result0 = 0
		return
	}
	result0 = value
	return
}

func sanitizeUnit(unit string) (result0 string) {
	result0 = strings.ToLower(strings.TrimSpace(unit))
	return
}

func requiresPriceQuantity(unit string) (result0 bool) {
	switch sanitizeUnit(unit) {
	case "kg", "g", "grama", "gramas", "l", "litro", "ml", "mililitro", "pacote", "pct", "pac", "cx", "caixa", "garrafa", "lata", "sache", "sachê":
		result0 = true
	default:
		result0 = false
	}
	return
}

func deriveAIPriceQuantity(quantity float64, unit string) (result0 float64) {
	normalizedQuantity := clampNonNegative(quantity)
	if !requiresPriceQuantity(unit) {
		result0 = 1
		return
	}
	if normalizedQuantity <= 0 {
		result0 = 1
		return
	}
	result0 = normalizedQuantity
	return
}

func calculateListTotals(items []shoppingModel.ShoppingListItem) (result0 float64, result1 float64) {
	__logParams := map[string]any{"items": items}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "calculateListTotals"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "calculateListTotals"), zap.Any("params", __logParams))
	for _, item := range items {
		priceQuantity := normalizePriceQuantity(item.PriceQuantity)
		quantityFactor := item.Quantity / priceQuantity
		if quantityFactor < 0 {
			quantityFactor = 0
		}
		result0 += item.EstimatedPrice * quantityFactor
		if item.Purchased {
			basePrice := resolveUnitPrice(item.ActualPrice, item.EstimatedPrice)
			result1 += basePrice * quantityFactor
		}
	}
	return
}

func normalizeStringSlice(values []string) (result0 []string) {
	__logParams := map[string]any{"values": values}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "normalizeStringSlice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "normalizeStringSlice"), zap.Any("params", __logParams))
	if values == nil {
		result0 = []string{}
		return
	}
	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, trimmed)
	}
	result0 = normalized
	return
}

func toStringSlice(values profileModel.StringArray) (result0 []string) {
	__logParams := map[string]any{"values": values}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "toStringSlice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "toStringSlice"), zap.Any("params", __logParams))
	if values == nil {
		result0 = []string{}
		return
	}
	result0 = append([]string(nil), values...)
	return
}

func toShoppingListStringSlice(values shoppingModel.StringArray) (result0 []string) {
	__logParams := map[string]any{"values": values}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "toShoppingListStringSlice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "toShoppingListStringSlice"), zap.Any("params", __logParams))
	if values == nil {
		result0 = []string{}
		return
	}
	result0 = append([]string(nil), values...)
	return
}

func normalizeShoppingListStringSlice(values []string) (result0 []string) {
	__logParams := map[string]any{"values": values}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "normalizeShoppingListStringSlice"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "normalizeShoppingListStringSlice"), zap.Any("params", __logParams))
	result0 = normalizeStringSlice(values)
	return
}

func (s *shoppingListService) determineBudget(input dto.GenerateAIShoppingListDTO, profile *profileModel.Profile) (result0 float64) {
	__logParams := map[string]any{"s": s, "input": input, "profile": profile}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.determineBudget"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.determineBudget"), zap.Any("params", __logParams))
	if input.MaxBudget != nil && *input.MaxBudget > 0 {
		result0 = *input.MaxBudget
		return
	}
	if profile != nil {
		if profile.PreferredBudget > 0 {
			result0 = profile.PreferredBudget
			return
		}
		if profile.MonthlyIncome > 0 {
			calculated := profile.MonthlyIncome * 0.15
			if calculated > 0 {
				result0 = calculated
				return
			}
		}
	}
	result0 = 300.0
	return
}

func (s *shoppingListService) resolvePantryNames(ctx context.Context, lists []*shoppingModel.ShoppingList) (result0 map[uuid.UUID]string) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "lists": lists}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.resolvePantryNames"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.resolvePantryNames"), zap.Any("params", __logParams))
	names := make(map[uuid.UUID]string)
	seen := make(map[uuid.UUID]struct{})
	for _, sl := range lists {
		if sl.PantryID == nil {
			continue
		}
		id := *sl.PantryID
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if pantry, err := s.pantryRepo.GetByID(ctx, id); err == nil {
			names[id] = pantry.Name
		}
	}
	result0 = names
	return
}

func (s *shoppingListService) convertToResponseDTO(ctx context.Context, sl *shoppingModel.ShoppingList) (result0 *dto.ShoppingListResponseDTO) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "sl": sl}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.convertToResponseDTO"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.convertToResponseDTO"), zap.Any("params", __logParams))
	items := make([]dto.ShoppingListItemResponseDTO, 0, len(sl.Items))
	for _, item := range sl.Items {
		items = append(items, *s.convertItemToResponseDTO(&item))
	}

	var pantryID *string
	pantryName := ""
	if sl.PantryID != nil {
		idStr := sl.PantryID.String()
		pantryID = &idStr
		pantryName = s.lookupPantryName(ctx, *sl.PantryID)
	}
	result0 = &dto.ShoppingListResponseDTO{
		ID:            sl.ID.String(),
		UserID:        sl.UserID.String(),
		PantryID:      pantryID,
		PantryName:    pantryName,
		Name:          sl.Name,
		Status:        sl.Status,
		TotalBudget:   sl.TotalBudget,
		EstimatedCost: sl.EstimatedCost,
		ActualCost:    sl.ActualCost,
		GeneratedBy:   sl.GeneratedBy,
		Items:         items,
		Preferences:   convertPreferencesToDTO(sl),
		CreatedAt:     sl.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     sl.UpdatedAt.Format(time.RFC3339),
	}
	return
}

func (s *shoppingListService) lookupPantryName(ctx context.Context, pantryID uuid.UUID) (result0 string) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantryID": pantryID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.lookupPantryName"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.lookupPantryName"), zap.Any("params", __logParams))
	pantry, err := s.pantryRepo.GetByID(ctx, pantryID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.lookupPantryName"), zap.Error(err), zap.Any("params", __logParams))
		result0 = ""
		return
	}
	result0 = pantry.Name
	return
}

func (s *shoppingListService) convertItemToResponseDTO(item *shoppingModel.ShoppingListItem) (result0 *dto.ShoppingListItemResponseDTO) {
	__logParams := map[string]any{"s": s, "item": item}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.convertItemToResponseDTO"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.convertItemToResponseDTO"), zap.Any("params", __logParams))
	var pantryItemID *string
	if item.PantryItemID != nil {
		id := item.PantryItemID.String()
		pantryItemID = &id
	}
	priceQuantity := normalizePriceQuantity(item.PriceQuantity)
	item.PriceQuantity = priceQuantity
	result0 = &dto.ShoppingListItemResponseDTO{
		ID:             item.ID.String(),
		ShoppingListID: item.ShoppingListID.String(),
		Name:           item.Name,
		Quantity:       item.Quantity,
		Unit:           item.Unit,
		EstimatedPrice: item.EstimatedPrice,
		PriceQuantity:  priceQuantity,
		ActualPrice:    item.ActualPrice,
		Category:       item.Category,
		Priority:       item.Priority,
		Purchased:      item.Purchased,
		Notes:          item.Notes,
		Source:         item.Source,
		PantryItemID:   pantryItemID,
		CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      item.UpdatedAt.Format(time.RFC3339),
	}
	return
}

func (s *shoppingListService) analyzePantryHistory(ctx context.Context, pantries []*pantryModel.Pantry) (result0 *PantryInsights, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "pantries": pantries}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.analyzePantryHistory"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.analyzePantryHistory"), zap.Any("params", __logParams))
	insights := &PantryInsights{
		CommonItems:       []ItemInsight{},
		LowStockItems:     []ItemInsight{},
		ExpiringSoonItems: []ItemInsight{},
		AverageItemPrice:  make(map[string]float64),
		Categories:        []string{},
		TotalItems:        0,
	}

	itemFreq := make(map[string]*ItemInsight)
	priceSum := make(map[string]float64)
	priceCount := make(map[string]int)
	categories := make(map[string]bool)

	// Analyze each pantry
	for range pantries {
		// Get items for this pantry - we'll need to query the item repository
		// For now, we'll create a basic insight structure
		insights.TotalItems += 10 // placeholder
	}

	// Convert maps to slices
	for _, insight := range itemFreq {
		if insight.Frequency >= 3 { // Items that appear in 3+ purchases
			insights.CommonItems = append(insights.CommonItems, *insight)
		}
	}

	for category := range categories {
		insights.Categories = append(insights.Categories, category)
	}

	for itemName, sum := range priceSum {
		if count := priceCount[itemName]; count > 0 {
			insights.AverageItemPrice[itemName] = sum / float64(count)
		}
	}
	result0 = insights
	result1 = nil
	return
}

func (s *shoppingListService) buildShoppingListPrompt(input dto.GenerateAIShoppingListDTO, preferences shoppingPreferences, profile *profileModel.Profile, insights *PantryInsights, budget float64, includeBasics bool) (result0 string, result1 error) {
	__logParams := map[string]any{"s": s, "input": input, "preferences": preferences, "profile": profile, "insights": insights, "budget": budget, "includeBasics": includeBasics}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.buildShoppingListPrompt"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.buildShoppingListPrompt"), zap.Any("params", __logParams))
	shoppingType := input.ShoppingType
	if shoppingType == "" {
		shoppingType = "general"
	}

	prompt := fmt.Sprintf(`Você é um assistente especializado em criar listas de compras inteligentes para brasileiros.

CONTEXTO DO USUÁRIO:
- Orçamento máximo: R$ %.2f
- Tipo de compra: %s
- Incluir itens básicos: %t
`, budget, shoppingType, includeBasics)

	if input.PeopleCount != nil {
		prompt += fmt.Sprintf("- Número de pessoas atendidas: %d\n", *input.PeopleCount)
	} else if preferences.HouseholdSize > 0 {
		prompt += fmt.Sprintf("- Número de pessoas atendidas: %d\n", preferences.HouseholdSize)
	}

	prompt += fmt.Sprintf(`
PREFERÊNCIAS DO USUÁRIO:
- Tamanho da família: %d pessoas
- Renda mensal informada: R$ %.2f
`, preferences.HouseholdSize, preferences.MonthlyIncome)

	if len(preferences.DietaryRestrictions) > 0 {
		prompt += fmt.Sprintf("- Restrições alimentares: %s\n", strings.Join(preferences.DietaryRestrictions, ", "))
	}

	if profile != nil {
		prompt += fmt.Sprintf("- Orçamento preferido do perfil: R$ %.2f\n", profile.PreferredBudget)
		prompt += fmt.Sprintf("- Frequência típica de compras: %s\n", profile.ShoppingFrequency)
	}

	if insights.TotalItems > 0 {
		prompt += fmt.Sprintf(`
ANÁLISE DA DESPENSA:
- Total de itens cadastrados: %d
- Categorias mais comuns: %s
`, insights.TotalItems, strings.Join(insights.Categories, ", "))

		if len(insights.CommonItems) > 0 {
			prompt += "- Itens frequentemente comprados:\n"
			for _, item := range insights.CommonItems {
				prompt += fmt.Sprintf("  * %s (%.2f %s) - R$ %.2f em média\n",
					item.Name, item.QuantityPattern, item.Unit, item.AveragePrice)
			}
		}
	}

	if len(input.ExcludeItems) > 0 {
		prompt += fmt.Sprintf("\nITENS PARA EXCLUIR: %s\n", strings.Join(input.ExcludeItems, ", "))
	}

	if input.Notes != "" {
		prompt += fmt.Sprintf("\nOBSERVAÇÕES ESPECIAIS: %s\n", input.Notes)
	}

	if input.Prompt != "" {
		prompt += fmt.Sprintf("\nINSTRUÇÕES PERSONALIZADAS: %s\n", input.Prompt)
	}

	prompt += `
INSTRUÇÕES:
1. Crie uma lista de compras balanceada e econômica
2. Considere preços médios de mercados brasileiros (usando dados de 2024/2025)
3. Priorize itens essenciais e de qualidade
4. Para produtos sem preço histórico, pesquise preços atuais no Brasil
5. Considere a proporção família/orçamento
6. Inclua uma breve explicação para cada item

FORMATO DE RESPOSTA (JSON):
{
  "items": [
    {
      "name": "Nome do produto",
      "quantity": 1.0,
      "unit": "unidade/kg/litro",
      "estimated_price": 0.00,
			"category": "categoria",
      "priority": 1,
      "reason": "motivo da inclusão"
    }
  ],
  "reasoning": "Explicação geral da lista",
  "estimated_total": 0.00
}

PRIORIDADES: 1=essencial, 2=importante, 3=desejável

Crie a lista agora:`
	result0 = prompt
	result1 = nil
	return
}

func (s *shoppingListService) parseAIResponse(userID uuid.UUID, input dto.GenerateAIShoppingListDTO, budget float64, preferences shoppingPreferences, aiResponse string) (result0 *shoppingModel.ShoppingList, result1 error) {
	__logParams := map[string]any{"s": s, "userID": userID, "input": input, "budget": budget, "preferences": preferences, "aiResponse": aiResponse}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*shoppingListService.parseAIResponse"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*shoppingListService.parseAIResponse"), zap.Any("params", __logParams))
	var aiList AIShoppingListResponse

	jsonStart := strings.Index(aiResponse, "{")
	jsonEnd := strings.LastIndex(aiResponse, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		result0 = nil
		result1 = domain.ErrAIResponseInvalid
		return
	}

	jsonResponse := aiResponse[jsonStart : jsonEnd+1]

	if err := json.Unmarshal([]byte(jsonResponse), &aiList); err != nil {
		zap.L().Error("function.error", zap.String("func", "*shoppingListService.parseAIResponse"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("%w: %v", domain.ErrAIResponseInvalid, err)
		return
	}

	pantryID := input.PantryID
	shoppingList := &shoppingModel.ShoppingList{
		UserID:              userID,
		PantryID:            &pantryID,
		Name:                input.Name,
		Status:              "pending",
		TotalBudget:         budget,
		EstimatedCost:       aiList.EstimatedTotal,
		GeneratedBy:         "ai",
		HouseholdSize:       preferences.HouseholdSize,
		MonthlyIncome:       preferences.MonthlyIncome,
		DietaryRestrictions: shoppingModel.StringArray(normalizeStringSlice(preferences.DietaryRestrictions)),
	}

	// Convert AI items to shopping list items
	for _, aiItem := range aiList.Items {
		quantity := clampNonNegative(aiItem.Quantity)
		priceQuantity := deriveAIPriceQuantity(quantity, aiItem.Unit)
		quantityFactor := 1.0
		if priceQuantity > 0 {
			quantityFactor = quantity / priceQuantity
		}
		if quantityFactor <= 0 {
			quantityFactor = 1
		}
		estimatedPrice := clampNonNegative(aiItem.EstimatedPrice)
		if estimatedPrice > 0 {
			estimatedPrice = estimatedPrice / quantityFactor
		}

		item := shoppingModel.ShoppingListItem{
			Name:           aiItem.Name,
			Quantity:       quantity,
			Unit:           aiItem.Unit,
			EstimatedPrice: estimatedPrice,
			PriceQuantity:  priceQuantity,
			Category:       aiItem.Category,
			Priority:       aiItem.Priority,
			Notes:          aiItem.Reason,
			Source:         "ai_suggestion",
		}

		// Validate priority
		if item.Priority < 1 || item.Priority > 3 {
			item.Priority = 2 // default to medium
		}

		shoppingList.Items = append(shoppingList.Items, item)
	}

	estimatedTotal, actualTotal := calculateListTotals(shoppingList.Items)
	if estimatedTotal > 0 {
		shoppingList.EstimatedCost = estimatedTotal
	}
	shoppingList.ActualCost = actualTotal
	result0 = shoppingList
	result1 = nil
	return
}
