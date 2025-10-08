package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/model"
)

type CreditService interface {
	GetWallet(ctx context.Context, userID uuid.UUID) (*dto.CreditWalletResponse, error)
	ConsumeCredit(ctx context.Context, userID uuid.UUID, description string) error
	AddCredit(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, amount int, description string) (*dto.CreditWalletResponse, error)
	ListTransactions(ctx context.Context, userID uuid.UUID, filter dto.TransactionFilter) ([]*dto.CreditTransactionResponse, error)
}

type CreditRepository interface {
	WithTx(ctx context.Context, fn func(repo CreditRepository) error) error
	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*model.CreditWallet, error)
	GetWalletByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*model.CreditWallet, error)
	CreateWallet(ctx context.Context, wallet *model.CreditWallet) error
	UpdateWallet(ctx context.Context, wallet *model.CreditWallet) error
	CreateTransaction(ctx context.Context, tx *model.CreditTransaction) error
	ListTransactions(ctx context.Context, userID uuid.UUID, filter dto.TransactionFilter) ([]*model.CreditTransaction, error)
}

type CreditHandler interface {
	GetWallet(c *gin.Context)
	ListTransactions(c *gin.Context)
	AddCredits(c *gin.Context)
}
