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
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type creditService struct {
	repo domain.CreditRepository
}

func NewCreditService(repo domain.CreditRepository) domain.CreditService {
	return &creditService{repo: repo}
}

func (s *creditService) GetWallet(ctx context.Context, userID uuid.UUID) (*dto.CreditWalletResponse, error) {
	logger := appLogger.FromContext(ctx)

	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet, err = s.createWalletWithDefaults(ctx, userID)
			if err != nil {
				logger.Error("failed to create wallet",
					zap.String(appLogger.FieldModule, "credits"),
					zap.String(appLogger.FieldFunction, "GetWallet"),
					zap.String(appLogger.FieldUserID, userID.String()),
					zap.Error(err),
				)
				return nil, err
			}
		} else {
			logger.Error("failed to get wallet",
				zap.String(appLogger.FieldModule, "credits"),
				zap.String(appLogger.FieldFunction, "GetWallet"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			return nil, err
		}
	}

	return toWalletResponse(wallet), nil
}

func (s *creditService) ConsumeCredit(ctx context.Context, userID uuid.UUID, description string) error {
	logger := appLogger.FromContext(ctx)

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
			logger.Warn("insufficient credits",
				zap.String(appLogger.FieldModule, "credits"),
				zap.String(appLogger.FieldFunction, "ConsumeCredit"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			return domain.ErrInsufficientCredits
		}
		logger.Error("failed to consume credit",
			zap.String(appLogger.FieldModule, "credits"),
			zap.String(appLogger.FieldFunction, "ConsumeCredit"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info("credit consumed",
		zap.String(appLogger.FieldModule, "credits"),
		zap.String(appLogger.FieldFunction, "ConsumeCredit"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)
	return nil
}

func (s *creditService) AddCredit(ctx context.Context, actorID uuid.UUID, targetUserID uuid.UUID, amount int, description string) (*dto.CreditWalletResponse, error) {
	logger := appLogger.FromContext(ctx)

	if amount <= 0 {
		logger.Warn("invalid credit amount",
			zap.String(appLogger.FieldModule, "credits"),
			zap.String(appLogger.FieldFunction, "AddCredit"),
			zap.Int(appLogger.FieldCount, amount),
		)
		return nil, domain.ErrInvalidCreditAmount
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
		logger.Error("failed to add credit",
			zap.String(appLogger.FieldModule, "credits"),
			zap.String(appLogger.FieldFunction, "AddCredit"),
			zap.String(appLogger.FieldUserID, targetUserID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("credit added",
		zap.String(appLogger.FieldModule, "credits"),
		zap.String(appLogger.FieldFunction, "AddCredit"),
		zap.String(appLogger.FieldUserID, targetUserID.String()),
		zap.Int(appLogger.FieldCount, amount),
	)

	return toWalletResponse(updatedWallet), nil
}

func (s *creditService) ListTransactions(ctx context.Context, userID uuid.UUID, filter dto.TransactionFilter) ([]*dto.CreditTransactionResponse, error) {
	logger := appLogger.FromContext(ctx)

	transactions, err := s.repo.ListTransactions(ctx, userID, filter)
	if err != nil {
		logger.Error("failed to list transactions",
			zap.String(appLogger.FieldModule, "credits"),
			zap.String(appLogger.FieldFunction, "ListTransactions"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	responses := make([]*dto.CreditTransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		responses = append(responses, toTransactionResponse(tx))
	}

	return responses, nil
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
