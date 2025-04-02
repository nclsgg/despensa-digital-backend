package service

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nclsgg/dispensa-digital/backend/config"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/dispensa-digital/backend/internal/utils"
	"github.com/redis/go-redis/v9"
)

type authService struct {
	repo  domain.AuthRepository
	cfg   *config.Config
	redis *redis.Client
}

func NewAuthService(repo domain.AuthRepository, cfg *config.Config, redis *redis.Client) domain.AuthService {
	return &authService{repo, cfg, redis}
}

func (s *authService) Register(ctx context.Context, user *model.User) error {
	isEmailValid := utils.IsEmailValid(user.Email)
	if !isEmailValid {
		return errors.New("invalid email")
	}

	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.repo.CreateUser(ctx, user)
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

func (s *authService) HashPassword(password string) (string, error) {
	hashedPassword, err := pbkdf2.Key(sha1.New, password, []byte("salt"), 4096, 32)
	if err != nil {
		return "", err
	}
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

func (s *authService) GenerateRefreshToken(ctx context.Context, userID uint64) (string, error) {
	token := uuid.NewString()
	key := fmt.Sprintf("refresh_token:%s", token)

	err := s.redis.Set(ctx, key, userID, time.Hour*24*7).Err()
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

	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return "", "", err
	}

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
