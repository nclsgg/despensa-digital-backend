package domain

import "errors"

var (
	ErrUnauthorized       = errors.New("item: user not authorized for this operation")
	ErrItemNotFound       = errors.New("item: not found")
	ErrInvalidPantry      = errors.New("item: invalid pantry id")
	ErrCategoryNotFound   = errors.New("item category: not found")
	ErrCategoryNotDefault = errors.New("item category: not default")
)
