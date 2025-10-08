package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type creditService struct {
	repo domain.CreditRepository
}

func NewCreditService(repo domain.CreditRepository) (result0 domain.CreditService) {
	__logParams := map[string]any{"repo": repo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewCreditService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewCreditService"), zap.Any("params", __logParams))
	result0 = &creditService{repo: repo}
	return
}

func (s *creditService) GetWallet(ctx context.Context, userID uuid.UUID) (result0 *dto.CreditWalletResponse, result1 error) {
	__logParams := map[string]any{"ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditService.GetWallet"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditService.GetWallet"), zap.Any("params", __logParams))

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet, err = s.createWalletWithDefaults(ctx, userID)
			if err != nil {
				result0 = nil
				result1 = err
				return
			}
		} else {
			result0 = nil
			result1 = err
			return
		}
	}

	result0 = toWalletResponse(wallet)
	result1 = nil
	return
}

func (s *creditService) ConsumeCredit(ctx context.Context, userID uuid.UUID, description string) (result0 error) {
	__logParams := map[string]any{"ctx": ctx, "userID": userID, "description": description}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditService.ConsumeCredit"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditService.ConsumeCredit"), zap.Any("params", __logParams))

	desc := strings.TrimSpace(description)
	if desc == "" {
		desc = "LLM request"
	}

	err := s.repo.WithTx(ctx, func(repo domain.CreditRepository) error {
		wallet, err := repo.GetWalletByUserIDForUpdate(ctx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrInsufficientCredits
			}
			return err
		}

		if wallet.Balance <= 0 {
			return domain.ErrInsufficientCredits
		}

		wallet.Balance -= 1
		if err := repo.UpdateWallet(ctx, wallet); err != nil {
			return err
		}

		transaction := &model.CreditTransaction{
			WalletID:    wallet.ID,
			UserID:      userID,
			Amount:      -1,
			Type:        "consume",
			Description: desc,
		}

		if err := repo.CreateTransaction(ctx, transaction); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, domain.ErrInsufficientCredits) {
			result0 = domain.ErrInsufficientCredits
			return
		}
		result0 = err
		return
	}

	result0 = nil
	return
}

func (s *creditService) AddCredit(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, amount int, description string) (result0 *dto.CreditWalletResponse, result1 error) {
	__logParams := map[string]any{"ctx": ctx, "actorID": actorID, "targetUserID": targetUserID, "amount": amount, "description": description}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditService.AddCredit"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditService.AddCredit"), zap.Any("params", __logParams))

	if amount <= 0 {
		result0 = nil
		result1 = domain.ErrInvalidCreditAmount
		return
	}

	desc := strings.TrimSpace(description)
	if desc == "" {
		desc = "Manual credit adjustment"
	}

	var updatedWallet *model.CreditWallet
	if err := s.repo.WithTx(ctx, func(repo domain.CreditRepository) error {
		wallet, err := repo.GetWalletByUserIDForUpdate(ctx, targetUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet, err = s.createWalletWithinTx(ctx, repo, targetUserID)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		wallet.Balance += amount
		if err := repo.UpdateWallet(ctx, wallet); err != nil {
			return err
		}

		creditTx := &model.CreditTransaction{
			WalletID:    wallet.ID,
			UserID:      targetUserID,
			Amount:      amount,
			Type:        "add",
			Description: desc,
		}

		if err := repo.CreateTransaction(ctx, creditTx); err != nil {
			return err
		}

		updatedWallet = wallet
		return nil
	}); err != nil {
		result0 = nil
		result1 = err
		return
	}

	result0 = toWalletResponse(updatedWallet)
	result1 = nil
	return
}

func (s *creditService) ListTransactions(ctx context.Context, userID uuid.UUID, filter dto.TransactionFilter) (result0 []*dto.CreditTransactionResponse, result1 error) {
	__logParams := map[string]any{"ctx": ctx, "userID": userID, "filter": filter}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditService.ListTransactions"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditService.ListTransactions"), zap.Any("params", __logParams))

	transactions, err := s.repo.ListTransactions(ctx, userID, filter)
	if err != nil {
		result0 = nil
		result1 = err
		return
	}

	responses := make([]*dto.CreditTransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		responses = append(responses, toTransactionResponse(tx))
	}

	result0 = responses
	result1 = nil
	return
}

func (s *creditService) createWalletWithDefaults(ctx context.Context, userID uuid.UUID) (*model.CreditWallet, error) {
	var wallet *model.CreditWallet
	err := s.repo.WithTx(ctx, func(repo domain.CreditRepository) error {
		var err error
		wallet, err = s.createWalletWithinTx(ctx, repo, userID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *creditService) createWalletWithinTx(ctx context.Context, repo domain.CreditRepository, userID uuid.UUID) (*model.CreditWallet, error) {
	wallet := &model.CreditWallet{
		UserID:  userID,
		Balance: 10,
	}

	if err := repo.CreateWallet(ctx, wallet); err != nil {
		return nil, err
	}

	initialTransaction := &model.CreditTransaction{
		WalletID:    wallet.ID,
		UserID:      userID,
		Amount:      10,
		Type:        "add",
		Description: "Initial credit allocation",
	}

	if err := repo.CreateTransaction(ctx, initialTransaction); err != nil {
		return nil, err
	}

	return wallet, nil
}

func toWalletResponse(wallet *model.CreditWallet) *dto.CreditWalletResponse {
	if wallet == nil {
		return nil
	}
	return &dto.CreditWalletResponse{
		WalletID:  wallet.ID.String(),
		UserID:    wallet.UserID.String(),
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: wallet.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toTransactionResponse(tx *model.CreditTransaction) *dto.CreditTransactionResponse {
	if tx == nil {
		return nil
	}
	return &dto.CreditTransactionResponse{
		TransactionID: tx.ID.String(),
		WalletID:      tx.WalletID.String(),
		UserID:        tx.UserID.String(),
		Amount:        tx.Amount,
		Type:          tx.Type,
		Description:   tx.Description,
		CreatedAt:     tx.CreatedAt.UTC().Format(time.RFC3339),
	}
}
