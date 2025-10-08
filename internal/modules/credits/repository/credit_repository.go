package repository

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/credits/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type creditRepository struct {
	db *gorm.DB
}

func NewCreditRepository(db *gorm.DB) (result0 domain.CreditRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewCreditRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewCreditRepository"), zap.Any("params", __logParams))
	result0 = &creditRepository{db: db}
	return
}

func (r *creditRepository) WithTx(ctx context.Context, fn func(repo domain.CreditRepository) error) (result0 error) {
	__logParams := map[string]any{"ctx": ctx, "fn": fn}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditRepository.WithTx"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditRepository.WithTx"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &creditRepository{db: tx}
		return fn(txRepo)
	})
	return
}

func (r *creditRepository) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (result0 *model.CreditWallet, result1 error) {
	__logParams := map[string]any{"ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditRepository.GetWalletByUserID"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditRepository.GetWalletByUserID"), zap.Any("params", __logParams))
	var wallet model.CreditWallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			zap.L().Error("function.error", zap.String("func", "*creditRepository.GetWalletByUserID"), zap.Error(err), zap.Any("params", __logParams))
		}
		result0 = nil
		result1 = err
		return
	}
	result0 = &wallet
	result1 = nil
	return
}

func (r *creditRepository) GetWalletByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (result0 *model.CreditWallet, result1 error) {
	__logParams := map[string]any{"ctx": ctx, "userID": userID}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditRepository.GetWalletByUserIDForUpdate"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditRepository.GetWalletByUserIDForUpdate"), zap.Any("params", __logParams))
	var wallet model.CreditWallet
	if err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			zap.L().Error("function.error", zap.String("func", "*creditRepository.GetWalletByUserIDForUpdate"), zap.Error(err), zap.Any("params", __logParams))
		}
		result0 = nil
		result1 = err
		return
	}
	result0 = &wallet
	result1 = nil
	return
}

func (r *creditRepository) CreateWallet(ctx context.Context, wallet *model.CreditWallet) (result0 error) {
	__logParams := map[string]any{"ctx": ctx, "wallet": wallet}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditRepository.CreateWallet"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditRepository.CreateWallet"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(wallet).Error
	return
}

func (r *creditRepository) UpdateWallet(ctx context.Context, wallet *model.CreditWallet) (result0 error) {
	__logParams := map[string]any{"ctx": ctx, "wallet": wallet}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditRepository.UpdateWallet"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditRepository.UpdateWallet"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(wallet).Error
	return
}

func (r *creditRepository) CreateTransaction(ctx context.Context, creditTx *model.CreditTransaction) (result0 error) {
	__logParams := map[string]any{"ctx": ctx, "creditTx": creditTx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditRepository.CreateTransaction"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditRepository.CreateTransaction"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(creditTx).Error
	return
}

func (r *creditRepository) ListTransactions(ctx context.Context, userID uuid.UUID, filter dto.TransactionFilter) (result0 []*model.CreditTransaction, result1 error) {
	__logParams := map[string]any{"ctx": ctx, "userID": userID, "filter": filter}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*creditRepository.ListTransactions"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*creditRepository.ListTransactions"), zap.Any("params", __logParams))

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filter.Type != nil && *filter.Type != "" {
		typeValue := strings.ToLower(strings.TrimSpace(*filter.Type))
		query = query.Where("LOWER(type) = ?", typeValue)
	}

	if filter.From != nil {
		query = query.Where("created_at >= ?", filter.From)
	}

	if filter.To != nil {
		query = query.Where("created_at <= ?", filter.To)
	}

	var transactions []*model.CreditTransaction
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		zap.L().Error("function.error", zap.String("func", "*creditRepository.ListTransactions"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}

	result0 = transactions
	result1 = nil
	return
}
