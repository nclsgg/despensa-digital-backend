package domain

import "errors"

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInsufficientCredits = errors.New("insufficient credits")
	ErrInvalidCreditAmount = errors.New("invalid credit amount")
	ErrTransactionFailed   = errors.New("credit transaction failed")
)
