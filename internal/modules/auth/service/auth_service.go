package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authService struct {
	repo  domain.AuthRepository
	cfg   *config.Config
	redis *redis.Client
}

func NewAuthService(repo domain.AuthRepository, cfg *config.Config, redis *redis.Client) (result0 domain.AuthService) {
	__logParams := map[string]any{"repo": repo, "cfg": cfg,

		// HashPassword creates a hash of the password (simplified for OAuth users)
		"redis": redis}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewAuthService"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewAuthService"), zap.Any("params", __logParams))
	result0 = &authService{repo, cfg, redis}
	return
}

func (s *authService) HashPassword(password string) (result0 string, result1 error) {
	__logParams := map[string]any{"s": s, "password": password}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authService.HashPassword"), zap.Any("result", map[string]any{

			// GenerateAccessToken generates JWT token for OAuth authenticated users
			"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authService.HashPassword"), zap.Any("params", __logParams))
	hasher := sha256.New()
	hasher.Write([]byte(password + "salt"))
	hashedPassword := hasher.Sum(nil)
	encoded := base64.StdEncoding.EncodeToString(hashedPassword)
	result0 = encoded
	result1 = nil
	return
}

func (s *authService) GenerateAccessToken(user *model.User) (result0 string, result1 error) {
	__logParams := map[string]any{"s": s, "user": user}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authService.GenerateAccessToken"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authService.GenerateAccessToken"), zap.Any("params", __logParams))
	jwtExpiration, err := time.ParseDuration(s.cfg.JWTExpiration)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*authService.GenerateAccessToken"), zap.Error(err), zap.Any("params", __logParams))
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
		zap.L().Error("function.error", zap.String("func", "*authService.GenerateAccessToken"), zap.Error(result1), zap.Any("params", __logParams))
		result0 = ""
		return
	}
	return
}

// OAuth specific methods
func (s *authService) GetUserByEmail(ctx context.Context, email string) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "email": email}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authService.GetUserByEmail"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}(

	// For OAuth users, we don't need password hashing since they authenticate via OAuth
	// But we'll set a placeholder password just in case
	)
	zap.L().Info("function.entry", zap.String("func", "*authService.GetUserByEmail"), zap.Any("params", __logParams))
	user, err := s.repo.GetUser(ctx, email)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*authService.GetUserByEmail"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = nil
			result1 = domain.ErrUserNotFound
			return
		}
		result0 = nil
		result1 = err
		return
	}
	result0 = user
	result1 = nil
	return
}

func (s *authService) CreateUserOAuth(ctx context.Context, user *model.User) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "user": user}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authService.CreateUserOAuth"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authService.CreateUserOAuth"), zap.Any("params", __logParams))

	if user.Password == "" {
		hashedPassword, err := s.HashPassword("oauth-no-password")
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "*authService.CreateUserOAuth"), zap.Error(err), zap.Any("params", __logParams))
			result0 = err
			return
		}
		user.Password = hashedPassword
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		zap.L().Error("function.error", zap.String("func", "*authService.CreateUserOAuth"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			result0 = domain.ErrEmailAlreadyRegistered
			return
		}
		result0 = err
		return
	}
	result0 = nil
	return
}

// CompleteProfile updates user profile with first/last name and marks profile as complete
func (s *authService) CompleteProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string) (result0 error) {
	__logParams := map[string]any{"s": s, "ctx": ctx, "userID": userID, "firstName": firstName, "lastName": lastName}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*authService.CompleteProfile"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*authService.CompleteProfile"), zap.Any("params", __logParams))
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*authService.CompleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result0 = domain.ErrUserNotFound
			return
		}
		result0 = err
		return
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfileCompleted = true

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		zap.L().Error("function.error", zap.String("func", "*authService.CompleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		result0 = domain.ErrProfileUpdateFailed
		return
	}
	result0 = nil
	return
}
