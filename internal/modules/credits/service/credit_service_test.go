package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCreditService(t *testing.T) domain.CreditService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&model.CreditWallet{}, &model.CreditTransaction{}))

	repo := repository.NewCreditRepository(db)
	return NewCreditService(repo)
}

func TestCreditService_GetWalletCreatesWalletWithDefaultBalance(t *testing.T) {
	t.Parallel()

	svc := setupCreditService(t)
	ctx := context.Background()
	userID := uuid.New()

	wallet, err := svc.GetWallet(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, wallet)
	require.Equal(t, 10, wallet.Balance)

	transactions, err := svc.ListTransactions(ctx, userID, dto.TransactionFilter{})
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	require.Equal(t, 10, transactions[0].Amount)
	require.Equal(t, "add", transactions[0].Type)
}

func TestCreditService_ConsumeCreditDecrementsBalance(t *testing.T) {
	t.Parallel()

	svc := setupCreditService(t)
	ctx := context.Background()
	userID := uuid.New()

	wallet, err := svc.GetWallet(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, wallet)

	require.NoError(t, svc.ConsumeCredit(ctx, userID, "test consumption"))

	wallet, err = svc.GetWallet(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, 9, wallet.Balance)

	transactions, err := svc.ListTransactions(ctx, userID, dto.TransactionFilter{})
	require.NoError(t, err)
	require.Len(t, transactions, 2)
	require.Equal(t, -1, transactions[0].Amount)
	require.Equal(t, "consume", transactions[0].Type)
}
