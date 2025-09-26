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
	"gorm.io/gorm"
)

type authService struct {
	repo  domain.AuthRepository
	cfg   *config.Config
	redis *redis.Client
}

func NewAuthService(repo domain.AuthRepository, cfg *config.Config, redis *redis.Client) domain.AuthService {
	return &authService{repo, cfg, redis}
}

// HashPassword creates a hash of the password (simplified for OAuth users)
func (s *authService) HashPassword(password string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(password + "salt"))
	hashedPassword := hasher.Sum(nil)
	encoded := base64.StdEncoding.EncodeToString(hashedPassword)
	return encoded, nil
}

// GenerateAccessToken generates JWT token for OAuth authenticated users
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
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// OAuth specific methods
func (s *authService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *authService) CreateUserOAuth(ctx context.Context, user *model.User) error {
	// For OAuth users, we don't need password hashing since they authenticate via OAuth
	// But we'll set a placeholder password just in case
	if user.Password == "" {
		hashedPassword, err := s.HashPassword("oauth-no-password")
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrEmailAlreadyRegistered
		}
		return err
	}

	return nil
}

// CompleteProfile updates user profile with first/last name and marks profile as complete
func (s *authService) CompleteProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string) error {
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfileCompleted = true

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return domain.ErrProfileUpdateFailed
	}

	return nil
}
