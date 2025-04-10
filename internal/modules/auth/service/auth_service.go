package service

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/internal/utils"
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

func (s *authService) Register(ctx context.Context, user *model.User) (string, string, error) {
	if !utils.IsEmailValid(user.Email) {
		return "", "", errors.New("invalid email")
	}

	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return "", "", err
	}
	user.Password = hashedPassword

	if err := s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", "", errors.New("email already registered")
		}

		return "", "", errors.New("failed to create user")
	}

	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.GetUser(ctx, email)
	if err != nil {
		return "", "", err
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return "", "", err
	}

	if user.Password != hashedPassword {
		return "", "", errors.New("invalid password")
	}

	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return s.redis.Del(ctx, key).Err()
}

func (s *authService) HashPassword(password string) (string, error) {
	hashedPassword, _ := pbkdf2.Key(sha1.New, password, []byte("salt"), 4096, 32)
	encoded := base64.StdEncoding.EncodeToString(hashedPassword)
	return encoded, nil
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
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *authService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := uuid.NewString()
	key := fmt.Sprintf("refresh_token:%s", token)

	err := s.redis.Set(ctx, key, userID.String(), time.Hour*24*7).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)

	userIDStr, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return "", "", err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", errors.New("invalid user id in token")
	}

	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// Invalida o token antigo
	s.redis.Del(ctx, key)

	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
