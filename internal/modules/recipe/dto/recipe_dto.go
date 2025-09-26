package dto

// AvailableIngredientDTO representa um ingrediente disponível em uma despensa
// exposto para o frontend.
type AvailableIngredientDTO struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}
