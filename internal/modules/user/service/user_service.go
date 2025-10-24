package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo}
}

func (s *userService) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	logger := appLogger.FromContext(ctx)

	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("User not found",
				zap.String(appLogger.FieldModule, "user"),
				zap.String(appLogger.FieldFunction, "GetUserById"),
				zap.String(appLogger.FieldUserID, id.String()),
			)
			return nil, domain.ErrUserNotFound
		}
		logger.Error("Failed to get user",
			zap.String(appLogger.FieldModule, "user"),
			zap.String(appLogger.FieldFunction, "GetUserById"),
			zap.String(appLogger.FieldUserID, id.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	logger := appLogger.FromContext(ctx)

	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		logger.Error("Failed to list users",
			zap.String(appLogger.FieldModule, "user"),
			zap.String(appLogger.FieldFunction, "GetAllUsers"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("list users: %w", err)
	}

	logger.Info("Users listed successfully",
		zap.String(appLogger.FieldModule, "user"),
		zap.String(appLogger.FieldFunction, "GetAllUsers"),
		zap.Int(appLogger.FieldCount, len(users)),
	)
	return users, nil
}

func (s *userService) CompleteProfile(ctx context.Context, id uuid.UUID, firstName, lastName string) error {
	logger := appLogger.FromContext(ctx)

	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("User not found for profile completion",
				zap.String(appLogger.FieldModule, "user"),
				zap.String(appLogger.FieldFunction, "CompleteProfile"),
				zap.String(appLogger.FieldUserID, id.String()),
			)
			return domain.ErrUserNotFound
		}
		logger.Error("Failed to get user for profile completion",
			zap.String(appLogger.FieldModule, "user"),
			zap.String(appLogger.FieldFunction, "CompleteProfile"),
			zap.String(appLogger.FieldUserID, id.String()),
			zap.Error(err),
		)
		return fmt.Errorf("get user: %w", err)
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfileCompleted = true

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		logger.Error("Failed to update user profile",
			zap.String(appLogger.FieldModule, "user"),
			zap.String(appLogger.FieldFunction, "CompleteProfile"),
			zap.String(appLogger.FieldUserID, id.String()),
			zap.Error(err),
		)
		return fmt.Errorf("update user: %w", err)
	}

	logger.Info("User profile completed successfully",
		zap.String(appLogger.FieldModule, "user"),
		zap.String(appLogger.FieldFunction, "CompleteProfile"),
		zap.String(appLogger.FieldUserID, id.String()),
	)
	return nil
}
