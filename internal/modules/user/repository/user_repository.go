package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) (result0 domain.UserRepository) {
	__logParams := map[string]any{"db": db}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewUserRepository"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewUserRepository"), zap.Any("params", __logParams))
	result0 = &userRepository{db}
	return
}

func (r *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*userRepository.GetUserById"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*userRepository.GetUserById"), zap.Any("params", __logParams))
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	result0 = &user
	result1 = err
	return
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "email": email}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*userRepository.GetUserByEmail"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*userRepository.GetUserByEmail"), zap.Any("params", __logParams))
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*userRepository.GetUserByEmail"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = &user
	result1 = nil
	return
}

func (r *userRepository) GetAllUsers(ctx context.Context) (result0 []model.User, result1 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*userRepository.GetAllUsers"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*userRepository.GetAllUsers"), zap.Any("params", __logParams))
	var users []model.User
	err := r.db.WithContext(ctx).Find(&users).Error
	result0 = users
	result1 = err
	return
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) (result0 error) {
	__logParams := map[string]any{"r": r, "ctx": ctx, "user": user}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*userRepository.UpdateUser"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*userRepository.UpdateUser"), zap.Any("params", __logParams))
	result0 = r.db.WithContext(ctx).Save(user).Error
	return
}
