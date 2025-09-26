package domain

import "errors"

var (
	ErrProviderNotSupported   = errors.New("auth: provider not supported")
	ErrEmailAlreadyRegistered = errors.New("auth: email already registered")
	ErrUserNotFound           = errors.New("auth: user not found")
	ErrProfileUpdateFailed    = errors.New("auth: profile update failed")
)
