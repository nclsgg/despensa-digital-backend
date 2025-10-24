package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authService struct {
	repo domain.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo domain.AuthRepository, cfg *config.Config) domain.AuthService {
	return &authService{repo, cfg}
}

func (s *authService) GenerateAccessToken(user *model.User) (string, error) {
	jwtExpiration, err := time.ParseDuration(s.cfg.JWTExpiration)
	if err != nil {
		return "", err
	}

	claims := model.MyClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiration)),
			Issuer:    s.cfg.JWTIssuer,
			Audience:  []string{s.cfg.JWTAudience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// OAuth specific methods
func (s *authService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	logger := appLogger.FromContext(ctx)

	user, err := s.repo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("User not found by email",
				zap.String(appLogger.FieldModule, "auth"),
				zap.String(appLogger.FieldFunction, "GetUserByEmail"),
				zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(email)),
			)
			return nil, domain.ErrUserNotFound
		}

		logger.Error("Failed to get user by email",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "GetUserByEmail"),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(email)),
			zap.Error(err),
		)
		return nil, err
	}

	return user, nil
}

func (s *authService) CreateUserOAuth(ctx context.Context, user *model.User) error {
	logger := appLogger.FromContext(ctx)

	if err := s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Warn("Duplicate email on OAuth user creation",
				zap.String(appLogger.FieldModule, "auth"),
				zap.String(appLogger.FieldFunction, "CreateUserOAuth"),
				zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(user.Email)),
			)
			return domain.ErrEmailAlreadyRegistered
		}

		logger.Error("Failed to create OAuth user",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "CreateUserOAuth"),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(user.Email)),
			zap.Error(err),
		)
		return err
	}

	logger.Info("OAuth user created successfully",
		zap.String(appLogger.FieldModule, "auth"),
		zap.String(appLogger.FieldFunction, "CreateUserOAuth"),
		zap.String(appLogger.FieldUserID, user.ID.String()),
		zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(user.Email)),
	)

	return nil
}

// CompleteProfile updates user profile with first/last name and marks profile as complete
func (s *authService) CompleteProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string) error {
	logger := appLogger.FromContext(ctx)

	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("User not found for profile completion",
				zap.String(appLogger.FieldModule, "auth"),
				zap.String(appLogger.FieldFunction, "CompleteProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			return domain.ErrUserNotFound
		}

		logger.Error("Failed to load user for profile completion",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "CompleteProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfileCompleted = true

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		logger.Error("Failed to update user profile",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "CompleteProfile"),
			zap.String(appLogger.FieldUserID, userID.String()),
			zap.Error(err),
		)
		return domain.ErrProfileUpdateFailed
	}

	logger.Info("User profile completed successfully",
		zap.String(appLogger.FieldModule, "auth"),
		zap.String(appLogger.FieldFunction, "CompleteProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)

	return nil
}
