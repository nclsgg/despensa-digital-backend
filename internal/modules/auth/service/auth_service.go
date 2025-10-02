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
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authService struct {
	repo domain.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo domain.AuthRepository, cfg *config.Config) (result0 domain.AuthService) {
	result0 = &authService{repo, cfg}
	return
}

func (s *authService) GenerateAccessToken(user *model.User) (result0 string, result1 error) {
	jwtExpiration, err := time.ParseDuration(s.cfg.JWTExpiration)
	if err != nil {
		zap.L().Error("auth_service.generate_access_token.parse_duration_failed",
			zap.String("jwt_expiration", s.cfg.JWTExpiration),
			zap.Error(err),
		)
		result0 = ""
		result1 = err
		return
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
	result0, result1 = token.SignedString([]byte(s.cfg.JWTSecret))
	if result1 != nil {
		zap.L().Error("auth_service.generate_access_token.sign_failed",
			zap.String("user_id", user.ID.String()),
			zap.Error(result1),
		)
		result0 = ""
		return
	}
	return
}

// OAuth specific methods
func (s *authService) GetUserByEmail(ctx context.Context, email string) (result0 *model.User, result1 error) {
	user, err := s.repo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Info("auth_service.get_user_by_email.not_found",
				zap.String("email", email),
			)
			result0 = nil
			result1 = domain.ErrUserNotFound
			return
		}
		zap.L().Error("auth_service.get_user_by_email.failed",
			zap.String("email", email),
			zap.Error(err),
		)
		result0 = nil
		result1 = err
		return
	}
	result0 = user
	result1 = nil
	return
}

func (s *authService) CreateUserOAuth(ctx context.Context, user *model.User) (result0 error) {
	if err := s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			zap.L().Warn("auth_service.create_oauth_user.duplicate_email",
				zap.String("email", user.Email),
			)
			result0 = domain.ErrEmailAlreadyRegistered
			return
		}
		zap.L().Error("auth_service.create_oauth_user.persist_failed",
			zap.String("email", user.Email),
			zap.Error(err),
		)
		result0 = err
		return
	}
	result0 = nil
	return
}

// CompleteProfile updates user profile with first/last name and marks profile as complete
func (s *authService) CompleteProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string) (result0 error) {
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Info("auth_service.complete_profile.user_not_found",
				zap.String("user_id", userID.String()),
			)
			result0 = domain.ErrUserNotFound
			return
		}
		zap.L().Error("auth_service.complete_profile.load_failed",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		result0 = err
		return
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfileCompleted = true

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		zap.L().Error("auth_service.complete_profile.persist_failed",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		result0 = domain.ErrProfileUpdateFailed
		return
	}
	result0 = nil
	return
}
