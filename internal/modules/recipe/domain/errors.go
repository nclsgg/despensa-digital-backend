package domain

import "errors"

var (
	ErrInvalidRequest     = errors.New("recipe: invalid request")
	ErrUnauthorized       = errors.New("recipe: user not authorized")
	ErrPantryNotFound     = errors.New("recipe: pantry not found")
	ErrNoIngredients      = errors.New("recipe: no ingredients available")
	ErrLLMRequest         = errors.New("recipe: llm request failed")
	ErrInvalidLLMResponse = errors.New("recipe: invalid llm response")
)
