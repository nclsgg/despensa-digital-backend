package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
}

type AuthService interface {
	Register(ctx context.Context, user *model.User) (accessToken string, refreshToken string, err error)
	HashPassword(password string) (string, error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
	GenerateAccessToken(user *model.User) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (accessToken string, refreshTokenOut string, err error)
}

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
}
