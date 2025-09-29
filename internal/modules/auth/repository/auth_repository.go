package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) (result0 domain.AuthRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewAuthRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewAuthRepository"), zap.Any("params", __logParams))
	result0 = &authRepository{db}
	return
}

func (r *authRepository) CreateUser(ctx context.Context, user *model.User) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "user": user}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authRepository.CreateUser"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authRepository.CreateUser"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Create(user).Error
	return
}

func (r *authRepository) GetUserById(ctx context.Context, id uuid.UUID) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authRepository.GetUserById"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authRepository.GetUserById"), zap.Any("params", __logParams))
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*authRepository.GetUserById"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &user
	result1 = nil
	return
}

func (r *authRepository) GetUser(ctx context.Context, email string) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "email": email}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authRepository.GetUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authRepository.GetUser"), zap.Any("params", __logParams))
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*authRepository.GetUser"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &user
	result1 = nil
	return
}

func (r *authRepository) UpdateUser(ctx context.Context, user *model.User) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "user": user}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authRepository.UpdateUser"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authRepository.UpdateUser"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(user).Error
	return
}
