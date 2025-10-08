package dto

import "time"

type CreditWalletResponse struct {
	WalletID  string `json:"wallet_id"`
	UserID    string `json:"user_id"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreditTransactionResponse struct {
	TransactionID string `json:"transaction_id"`
	WalletID      string `json:"wallet_id"`
	UserID        string `json:"user_id"`
	Amount        int    `json:"amount"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
}

type TransactionFilter struct {
	Type   *string
	Limit  int
	Offset int
	From   *time.Time
	To     *time.Time
}

type AddCreditsRequest struct {
	UserID      *string `json:"user_id,omitempty"`
	Amount      int     `json:"amount"`
	Description string  `json:"description"`
}
