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

	if input.PantryID != nil {
		hasAccess, err := s.pantryRepo.IsUserInPantry(ctx, *input.PantryID, userID)
		if err != nil {
			return nil, err
		}
		if !hasAccess {
			return nil, domain.ErrPantryAccessDenied
		}
	}

	var estimatedCost float64
	for _, itemDto := range input.Items {
		estimatedCost += itemDto.EstimatedPrice * itemDto.Quantity
	}
	shoppingList.EstimatedCost = estimatedCost

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
			item.Priority = 3
		}
		shoppingList.Items = append(shoppingList.Items, *item)
	}

	if err := s.shoppingListRepo.Create(ctx, shoppingList); err != nil {
		return nil, fmt.Errorf("create shopping list: %w", err)
	}

	created, err := s.shoppingListRepo.GetByID(ctx, shoppingList.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShoppingListNotFound
		}
		return nil, fmt.Errorf("reload shopping list: %w", err)
	}

	return s.convertToResponseDTO(ctx, created), nil
}

func (s *shoppingListService) GetShoppingListByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.ShoppingListResponseDTO, error) {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShoppingListNotFound
		}
		return nil, fmt.Errorf("get shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return s.convertToResponseDTO(ctx, shoppingList), nil
}

func (s *shoppingListService) GetShoppingListsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*dto.ShoppingListSummaryDTO, error) {
	shoppingLists, err := s.shoppingListRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting shopping lists: %w", err)
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
			CreatedAt:      sl.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      sl.UpdatedAt.Format(time.RFC3339),
		})
	}

	return summaries, nil
}

func (s *shoppingListService) UpdateShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID, input dto.UpdateShoppingListDTO) (*dto.ShoppingListResponseDTO, error) {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShoppingListNotFound
		}
		return nil, fmt.Errorf("get shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

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
		return nil, fmt.Errorf("update shopping list: %w", err)
	}

	updated, err := s.shoppingListRepo.GetByID(ctx, shoppingList.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShoppingListNotFound
		}
		return nil, fmt.Errorf("reload shopping list: %w", err)
	}

	return s.convertToResponseDTO(ctx, updated), nil
}

func (s *shoppingListService) DeleteShoppingList(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrShoppingListNotFound
		}
		return fmt.Errorf("get shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return domain.ErrUnauthorized
	}

	if err := s.shoppingListRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete shopping list: %w", err)
	}

	return nil
}

func (s *shoppingListService) UpdateShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID, input dto.UpdateShoppingListItemDTO) (*dto.ShoppingListItemResponseDTO, error) {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShoppingListNotFound
		}
		return nil, fmt.Errorf("get shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	var targetItem *shoppingModel.ShoppingListItem
	for idx := range shoppingList.Items {
		if shoppingList.Items[idx].ID == itemID {
			targetItem = &shoppingList.Items[idx]
			break
		}
	}

	if targetItem == nil {
		return nil, domain.ErrItemNotFound
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
		return nil, fmt.Errorf("update shopping list item: %w", err)
	}

	reloadedItems, err := s.shoppingListRepo.GetItemsByShoppingListID(ctx, shoppingListID)
	if err == nil {
		for _, item := range reloadedItems {
			if item.ID == itemID {
				return s.convertItemToResponseDTO(item), nil
			}
		}
	}

	return s.convertItemToResponseDTO(targetItem), nil
}

func (s *shoppingListService) DeleteShoppingListItem(ctx context.Context, userID uuid.UUID, shoppingListID uuid.UUID, itemID uuid.UUID) error {
	shoppingList, err := s.shoppingListRepo.GetByID(ctx, shoppingListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrShoppingListNotFound
		}
		return fmt.Errorf("get shopping list: %w", err)
	}

	if shoppingList.UserID != userID {
		return domain.ErrUnauthorized
	}

	found := false
	for _, item := range shoppingList.Items {
		if item.ID == itemID {
			found = true
			break
		}
	}

	if !found {
		return domain.ErrItemNotFound
	}

	if err := s.shoppingListRepo.DeleteItem(ctx, itemID); err != nil {
		return fmt.Errorf("delete shopping list item: %w", err)
	}

	return nil
}

func (s *shoppingListService) GenerateAIShoppingList(ctx context.Context, userID uuid.UUID, input dto.GenerateAIShoppingListDTO) (*dto.ShoppingListResponseDTO, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get user profile: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		profile = nil
	}

	pantry, err := s.pantryRepo.GetByID(ctx, input.PantryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPantryNotFound
		}
		return nil, fmt.Errorf("get pantry: %w", err)
	}

	hasAccess, err := s.pantryRepo.IsUserInPantry(ctx, input.PantryID, userID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, domain.ErrPantryAccessDenied
	}

	pantryInsights, err := s.analyzePantryHistory(ctx, []*pantryModel.Pantry{pantry})
	if err != nil {
		return nil, fmt.Errorf("analyze pantry history: %w", err)
	}

	budget := s.determineBudget(input, profile)
	includeBasics := true
	if input.IncludeBasics != nil {
		includeBasics = *input.IncludeBasics
	}

	prompt, err := s.buildShoppingListPrompt(input, profile, pantryInsights, budget, includeBasics)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrPromptBuildFailed, err)
	}

	llmResponse, err := s.llmService.GenerateText(ctx, prompt, map[string]interface{}{
		"max_tokens":      2000,
		"temperature":     0.7,
		"response_format": "json",
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrAIRequestFailed, err)
	}

	shoppingList, err := s.parseAIResponse(userID, input, budget, llmResponse.Response)
	if err != nil {
		return nil, err
	}

	if err := s.shoppingListRepo.Create(ctx, shoppingList); err != nil {
		return nil, fmt.Errorf("create ai shopping list: %w", err)
	}

	created, err := s.shoppingListRepo.GetByID(ctx, shoppingList.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShoppingListNotFound
		}
		return nil, fmt.Errorf("reload shopping list: %w", err)
	}

	return s.convertToResponseDTO(ctx, created), nil
}

// Helper methods

func (s *shoppingListService) determineBudget(input dto.GenerateAIShoppingListDTO, profile *profileModel.Profile) float64 {
	if input.MaxBudget != nil && *input.MaxBudget > 0 {
		return *input.MaxBudget
	}
	if profile != nil {
		if profile.PreferredBudget > 0 {
			return profile.PreferredBudget
		}
		if profile.MonthlyIncome > 0 {
			calculated := profile.MonthlyIncome * 0.15
			if calculated > 0 {
				return calculated
			}
		}
	}
	return 300.0
}

func (s *shoppingListService) resolvePantryNames(ctx context.Context, lists []*shoppingModel.ShoppingList) map[uuid.UUID]string {
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
	return names
}

func (s *shoppingListService) convertToResponseDTO(ctx context.Context, sl *shoppingModel.ShoppingList) *dto.ShoppingListResponseDTO {
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

	return &dto.ShoppingListResponseDTO{
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
		CreatedAt:     sl.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     sl.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *shoppingListService) lookupPantryName(ctx context.Context, pantryID uuid.UUID) string {
	pantry, err := s.pantryRepo.GetByID(ctx, pantryID)
	if err != nil {
		return ""
	}
	return pantry.Name
}

func (s *shoppingListService) convertItemToResponseDTO(item *shoppingModel.ShoppingListItem) *dto.ShoppingListItemResponseDTO {
	return &dto.ShoppingListItemResponseDTO{
		ID:             item.ID.String(),
		ShoppingListID: item.ShoppingListID.String(),
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
		CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      item.UpdatedAt.Format(time.RFC3339),
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

func (s *shoppingListService) buildShoppingListPrompt(input dto.GenerateAIShoppingListDTO, profile *profileModel.Profile, insights *PantryInsights, budget float64, includeBasics bool) (string, error) {
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
	} else if profile != nil && profile.HouseholdSize > 0 {
		prompt += fmt.Sprintf("- Número de pessoas atendidas: %d\n", profile.HouseholdSize)
	}

	if profile != nil {
		prompt += fmt.Sprintf(`
PERFIL DO USUÁRIO:
- Renda mensal: R$ %.2f
- Orçamento preferido: R$ %.2f
- Frequência de compras: %s
`, profile.MonthlyIncome, profile.PreferredBudget, profile.ShoppingFrequency)

		if len(profile.DietaryRestrictions) > 0 {
			prompt += fmt.Sprintf("- Restrições alimentares: %s\n", strings.Join(profile.DietaryRestrictions, ", "))
		}
	}

	preferredBrands := make([]string, 0)
	if profile != nil {
		preferredBrands = append(preferredBrands, profile.PreferredBrands...)
	}
	if len(input.PreferredBrands) > 0 {
		preferredBrands = append(preferredBrands, input.PreferredBrands...)
	}
	if len(preferredBrands) > 0 {
		brandSet := make(map[string]struct{})
		dedup := make([]string, 0, len(preferredBrands))
		for _, brand := range preferredBrands {
			brand = strings.TrimSpace(brand)
			if brand == "" {
				continue
			}
			if _, exists := brandSet[brand]; exists {
				continue
			}
			brandSet[brand] = struct{}{}
			dedup = append(dedup, brand)
		}
		if len(dedup) > 0 {
			prompt += fmt.Sprintf("- Marcas preferidas: %s\n", strings.Join(dedup, ", "))
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

func (s *shoppingListService) parseAIResponse(userID uuid.UUID, input dto.GenerateAIShoppingListDTO, budget float64, aiResponse string) (*shoppingModel.ShoppingList, error) {
	var aiList AIShoppingListResponse

	jsonStart := strings.Index(aiResponse, "{")
	jsonEnd := strings.LastIndex(aiResponse, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, domain.ErrAIResponseInvalid
	}

	jsonResponse := aiResponse[jsonStart : jsonEnd+1]

	if err := json.Unmarshal([]byte(jsonResponse), &aiList); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrAIResponseInvalid, err)
	}

	pantryID := input.PantryID
	shoppingList := &shoppingModel.ShoppingList{
		UserID:        userID,
		PantryID:      &pantryID,
		Name:          input.Name,
		Status:        "pending",
		TotalBudget:   budget,
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
