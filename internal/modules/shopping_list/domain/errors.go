package domain

import "errors"

var (
	ErrShoppingListNotFound = errors.New("shopping_list: not found")
	ErrItemNotFound         = errors.New("shopping_list: item not found")
	ErrUnauthorized         = errors.New("shopping_list: unauthorized")
	ErrPantryNotFound       = errors.New("shopping_list: pantry not found")
	ErrPantryAccessDenied   = errors.New("shopping_list: pantry access denied")
	ErrPromptBuildFailed    = errors.New("shopping_list: prompt build failed")
	ErrAIResponseInvalid    = errors.New("shopping_list: ai response invalid")
	ErrAIRequestFailed      = errors.New("shopping_list: ai request failed")
)
