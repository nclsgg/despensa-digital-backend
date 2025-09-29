package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) (result0 domain.UserService) {
	__logParams := map[string]any{"repo": repo}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewUserService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewUserService"), zap.Any("params", __logParams))
	result0 = &userService{repo}
	return
}

func (s *userService) GetUserById(ctx context.Context, id uuid.UUID) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*userService.GetUserById"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*userService.GetUserById"), zap.Any("params", __logParams))
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*userService.GetUserById"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrUserNotFound
			return
		}
		result0 = nil
		result1 = fmt.Errorf("get user: %w", err)
		return
	}
	result0 = user
	result1 = nil
	return
}

func (s *userService) GetAllUsers(ctx context.Context) (result0 []model.User, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*userService.GetAllUsers"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*userService.GetAllUsers"), zap.Any("params", __logParams))
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*userService.GetAllUsers"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = fmt.Errorf("list users: %w", err)
		return
	}
	result0 = users
	result1 = nil
	return
}

func (s *userService) CompleteProfile(ctx context.Context, id uuid.UUID, firstName, lastName string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "id": id, "firstName": firstName, "lastName": lastName}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*userService.CompleteProfile"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*userService.CompleteProfile"), zap.Any("params", __logParams))
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*userService.CompleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = domain.ErrUserNotFound
			return
		}
		result0 = fmt.Errorf("get user: %w", err)
		return
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfileCompleted = true

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		zap.L().Error("function.error", zap.String("func", "*userService.CompleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		result0 = fmt.Errorf("update user: %w", err)
		return
	}
	result0 = nil
	return
}
