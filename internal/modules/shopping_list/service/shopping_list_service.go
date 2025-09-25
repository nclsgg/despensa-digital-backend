package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	itemDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/item/domain"
	llmDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/llm/domain"
	pantryDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	pantryModel "github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	profileDomain "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	profileModel "github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/dto"
	shoppingModel "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
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
	Brand          string  `json:"brand"`
	Priority       int     `json:"priority"`
	Reason         string  `json:"reason"`
}

type AIShoppingListResponse struct {
	Items          []AIShoppingItem `json:"items"`
	Reasoning      string           `json:"reasoning"`
	EstimatedTotal float64          `json:"estimated_total"`
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
) domain.ShoppingListService {
	return &shoppingListService{
		shoppingListRepo: shoppingListRepo,
		pantryRepo:       pantryRepo,
		itemRepo:         itemRepo,
		profileRepo:      profileRepo,
		llmService:       llmService,
	}
}

func (s *shoppingListService) CreateShoppingList(ctx context.Context, userID uuid.UUID, input dto.CreateShoppingListDTO) (*dto.ShoppingListResponseDTO, error) {
	shoppingList := &shoppingModel.ShoppingList{
		UserID:      userID,
		PantryID:    input.PantryID,
		Name:        input.Name,
		TotalBudget: input.TotalBudget,
		Status:      "pending",
		GeneratedBy: "manual",
	}

	// Validate pantry access if PantryID is provided
	if input.PantryID != nil {
		hasAccess, err := s.pantryRepo.IsUserInPantry(ctx, *input.PantryID, userID)
		if err != nil || !hasAccess {
			return nil, fmt.Errorf("access denied or pantry not found")
		}
	}

	// Calculate estimated cost
	var estimatedCost float64
	for _, itemDto := range input.Items {
		estimatedCost += itemDto.EstimatedPrice * itemDto.Quantity
	}
	shoppingList.EstimatedCost = estimatedCost

	// Create items
	for _, itemDto := range input.Items {
		item := &shoppingModel.ShoppingListItem{
			Name:           itemDto.Name,
			Quantity:       itemDto.Quantity,
			Unit:           itemDto.Unit,
			EstimatedPrice: itemDto.EstimatedPrice,
			Category:       itemDto.Category,
			Brand:          itemDto.Brand,
			Priority:       itemDto.Priority,
			Notes:          itemDto.Notes,
			Source:         "manual",
		}
		if item.Priority == 0 {
			item.Priority = 3 // default to low priority
		}
		shoppingList.Items = append(shoppingList.Items, *item)
	}

	if err := s.shoppingListRepo.Create(ctx, shoppingList); err != nil {
		return nil, fmt.Errorf("error creating shopping list: %w", err)
	}

	return s.convertToResponseDTO(shoppingList), nil
}

func (s *shoppingListService) GetShoppingListByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.ShoppingListResponseDTO, error) {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("shopping list not found")
		}
		return nil, fmt.Errorf("error getting shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return nil, fmt.Errorf("shopping list not found")
	}

	return s.convertToResponseDTO(shoppingList), nil
}

func (s *shoppingListService) GetShoppingListsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*dto.ShoppingListSummaryDTO, error) {
	shoppingLists, err := s.shoppingListRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting shopping lists: %w", err)
	}

	var summaries []*dto.ShoppingListSummaryDTO
	for _, sl := range shoppingLists {
		// Get item counts
		items, _ := s.shoppingListRepo.GetItemsByShoppingListID(ctx, sl.ID)
		itemCount := len(items)
		purchasedCount := 0
		for _, item := range items {
			if item.Purchased {
				purchasedCount++
			}
		}

		summary := &dto.ShoppingListSummaryDTO{
			ID:             sl.ID,
			PantryID:       sl.PantryID,
			Name:           sl.Name,
			Status:         sl.Status,
			TotalBudget:    sl.TotalBudget,
			EstimatedCost:  sl.EstimatedCost,
			ActualCost:     sl.ActualCost,
			GeneratedBy:    sl.GeneratedBy,
			ItemCount:      itemCount,
			PurchasedCount: purchasedCount,
			CreatedAt:      sl.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      sl.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// Get pantry name if PantryID is set
		if sl.PantryID != nil {
			pantry, err := s.pantryRepo.GetByID(ctx, *sl.PantryID)
			if err == nil {
				summary.PantryName = pantry.Name
			}
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (s *shoppingListService) UpdateShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID, input dto.UpdateShoppingListDTO) (*dto.ShoppingListResponseDTO, error) {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("shopping list not found")
		}
		return nil, fmt.Errorf("error getting shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return nil, fmt.Errorf("shopping list not found")
	}

	// Update fields if provided
	if input.Name != nil {
		shoppingList.Name = *input.Name
	}
	if input.Status != nil {
		shoppingList.Status = *input.Status
	}
	if input.TotalBudget != nil {
		shoppingList.TotalBudget = *input.TotalBudget
	}
	if input.ActualCost != nil {
		shoppingList.ActualCost = *input.ActualCost
	}

	if err := s.shoppingListRepo.Update(ctx, shoppingList); err != nil {
		return nil, fmt.Errorf("error updating shopping list: %w", err)
	}

	return s.convertToResponseDTO(shoppingList), nil
}

func (s *shoppingListService) DeleteShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("shopping list not found")
		}
		return fmt.Errorf("error getting shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return fmt.Errorf("shopping list not found")
	}

	if err := s.shoppingListRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting shopping list: %w", err)
	}

	return nil
}

func (s *shoppingListService) UpdateShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID, input dto.UpdateShoppingListItemDTO) (*dto.ShoppingListItemResponseDTO, error) {
	// First verify the shopping list belongs to the user
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("shopping list not found")
		}
		return nil, fmt.Errorf("error getting shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return nil, fmt.Errorf("shopping list not found")
	}

	// Find the item
	var targetItem *shoppingModel.ShoppingListItem
	for _, item := range shoppingList.Items {
		if item.ID == itemID {
			targetItem = &item
			break
		}
	}

	if targetItem == nil {
		return nil, fmt.Errorf("item not found")
	}

	// Update fields if provided
	if input.Name != nil {
		targetItem.Name = *input.Name
	}
	if input.Quantity != nil {
		targetItem.Quantity = *input.Quantity
	}
	if input.Unit != nil {
		targetItem.Unit = *input.Unit
	}
	if input.ActualPrice != nil {
		targetItem.ActualPrice = *input.ActualPrice
	}
	if input.Category != nil {
		targetItem.Category = *input.Category
	}
	if input.Brand != nil {
		targetItem.Brand = *input.Brand
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

	if err := s.shoppingListRepo.UpdateItem(ctx, targetItem); err != nil {
		return nil, fmt.Errorf("error updating shopping list item: %w", err)
	}

	return s.convertItemToResponseDTO(targetItem), nil
}

func (s *shoppingListService) DeleteShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID) error {
	// First verify the shopping list belongs to the user
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("shopping list not found")
		}
		return fmt.Errorf("error getting shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return fmt.Errorf("shopping list not found")
	}

	// Verify the item exists in the shopping list
	found := false
	for _, item := range shoppingList.Items {
		if item.ID == itemID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found")
	}

	if err := s.shoppingListRepo.DeleteItem(ctx, itemID); err != nil {
		return fmt.Errorf("error deleting shopping list item: %w", err)
	}

	return nil
}

func (s *shoppingListService) GenerateAIShoppingList(ctx context.Context, userID uuid.UUID, input dto.GenerateAIShoppingListDTO) (*dto.ShoppingListResponseDTO, error) {
	// Get user profile
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error getting user profile: %w", err)
	}

	// Get the specific pantry and verify ownership
	pantry, err := s.pantryRepo.GetByID(ctx, input.PantryID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("pantry not found")
		}
		return nil, fmt.Errorf("error getting pantry: %w", err)
	}

	// Verify pantry ownership
	if pantry.OwnerID != userID {
		return nil, fmt.Errorf("pantry does not belong to user")
	}

	// Analyze pantry data and get insights for the selected pantry
	pantries := []*pantryModel.Pantry{pantry}
	pantryInsights, err := s.analyzePantryHistory(ctx, pantries)
	if err != nil {
		return nil, fmt.Errorf("error analyzing pantry history: %w", err)
	}

	// Build LLM prompt for shopping list generation
	prompt, err := s.buildShoppingListPrompt(input, profile, pantryInsights)
	if err != nil {
		return nil, fmt.Errorf("error building prompt: %w", err)
	}

	// Call LLM to generate shopping list
	llmResponse, err := s.llmService.GenerateText(ctx, prompt, map[string]interface{}{
		"max_tokens":      2000,
		"temperature":     0.7,
		"response_format": "json",
	})
	if err != nil {
		return nil, fmt.Errorf("error generating AI shopping list: %w", err)
	}

	// Parse LLM response and create shopping list
	shoppingList, err := s.parseAIResponse(userID, input, llmResponse.Response)
	if err != nil {
		return nil, fmt.Errorf("error parsing AI response: %w", err)
	}

	// Create the shopping list in database
	if err := s.shoppingListRepo.Create(ctx, shoppingList); err != nil {
		return nil, fmt.Errorf("error creating AI shopping list: %w", err)
	}

	return s.convertToResponseDTO(shoppingList), nil
}

// Helper methods

func (s *shoppingListService) convertToResponseDTO(sl *shoppingModel.ShoppingList) *dto.ShoppingListResponseDTO {
	var items []dto.ShoppingListItemResponseDTO
	for _, item := range sl.Items {
		items = append(items, *s.convertItemToResponseDTO(&item))
	}

	responseDTO := &dto.ShoppingListResponseDTO{
		ID:            sl.ID,
		UserID:        sl.UserID,
		PantryID:      sl.PantryID,
		Name:          sl.Name,
		Status:        sl.Status,
		TotalBudget:   sl.TotalBudget,
		EstimatedCost: sl.EstimatedCost,
		ActualCost:    sl.ActualCost,
		GeneratedBy:   sl.GeneratedBy,
		Items:         items,
		CreatedAt:     sl.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     sl.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Get pantry name if PantryID is set
	if sl.PantryID != nil {
		pantry, err := s.pantryRepo.GetByID(context.Background(), *sl.PantryID)
		if err == nil {
			responseDTO.PantryName = pantry.Name
		}
	}

	return responseDTO
}

func (s *shoppingListService) convertItemToResponseDTO(item *shoppingModel.ShoppingListItem) *dto.ShoppingListItemResponseDTO {
	return &dto.ShoppingListItemResponseDTO{
		ID:             item.ID,
		ShoppingListID: item.ShoppingListID,
		Name:           item.Name,
		Quantity:       item.Quantity,
		Unit:           item.Unit,
		EstimatedPrice: item.EstimatedPrice,
		ActualPrice:    item.ActualPrice,
		Category:       item.Category,
		Brand:          item.Brand,
		Priority:       item.Priority,
		Purchased:      item.Purchased,
		Notes:          item.Notes,
		Source:         item.Source,
		CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      item.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (s *shoppingListService) analyzePantryHistory(ctx context.Context, pantries []*pantryModel.Pantry) (*PantryInsights, error) {
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

	return insights, nil
}

func (s *shoppingListService) buildShoppingListPrompt(input dto.GenerateAIShoppingListDTO, profile *profileModel.Profile, insights *PantryInsights) (string, error) {
	prompt := fmt.Sprintf(`Você é um assistente especializado em criar listas de compras inteligentes para brasileiros.

CONTEXTO DO USUÁRIO:
- Orçamento para compras: R$ %.2f
- Tipo de compra: %s
- Incluir itens básicos: %t

`, input.TotalBudget, input.ShoppingType, input.IncludeBasics)

	if profile != nil {
		prompt += fmt.Sprintf(`PERFIL DO USUÁRIO:
- Renda mensal: R$ %.2f
- Orçamento preferido: R$ %.2f
- Tamanho da família: %d pessoas
- Frequência de compras: %s
`, profile.MonthlyIncome, profile.PreferredBudget, profile.HouseholdSize, profile.ShoppingFrequency)

		if len(profile.DietaryRestrictions) > 0 {
			prompt += fmt.Sprintf("- Restrições alimentares: %s\n", strings.Join(profile.DietaryRestrictions, ", "))
		}

		if len(profile.PreferredBrands) > 0 {
			prompt += fmt.Sprintf("- Marcas preferidas: %s\n", strings.Join(profile.PreferredBrands, ", "))
		}
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
      "brand": "marca sugerida ou genérico",
      "priority": 1,
      "reason": "motivo da inclusão"
    }
  ],
  "reasoning": "Explicação geral da lista",
  "estimated_total": 0.00
}

PRIORIDADES: 1=essencial, 2=importante, 3=desejável

Crie a lista agora:`

	return prompt, nil
}

func (s *shoppingListService) parseAIResponse(userID uuid.UUID, input dto.GenerateAIShoppingListDTO, aiResponse string) (*shoppingModel.ShoppingList, error) {
	var aiList AIShoppingListResponse

	// Clean the response - sometimes AI adds extra text
	jsonStart := strings.Index(aiResponse, "{")
	jsonEnd := strings.LastIndex(aiResponse, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("invalid JSON response from AI")
	}

	jsonResponse := aiResponse[jsonStart : jsonEnd+1]

	if err := json.Unmarshal([]byte(jsonResponse), &aiList); err != nil {
		return nil, fmt.Errorf("error parsing AI response: %w", err)
	}

	// Create shopping list
	shoppingList := &shoppingModel.ShoppingList{
		UserID:        userID,
		Name:          input.Name,
		Status:        "pending",
		TotalBudget:   input.TotalBudget,
		EstimatedCost: aiList.EstimatedTotal,
		GeneratedBy:   "ai",
	}

	// Convert AI items to shopping list items
	for _, aiItem := range aiList.Items {
		item := shoppingModel.ShoppingListItem{
			Name:           aiItem.Name,
			Quantity:       aiItem.Quantity,
			Unit:           aiItem.Unit,
			EstimatedPrice: aiItem.EstimatedPrice,
			Category:       aiItem.Category,
			Brand:          aiItem.Brand,
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

	return shoppingList, nil
}
