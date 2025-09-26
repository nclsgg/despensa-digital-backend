package domain

import "errors"

var (
	ErrProfileAlreadyExists = errors.New("profile: already exists")
	ErrProfileNotFound      = errors.New("profile: not found")
)
