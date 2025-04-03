package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, id uint64) (*model.User, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
}

type AuthService interface {
	Register(ctx context.Context, user *model.User) error
	HashPassword(password string) (string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	GenerateAccessToken(user *model.User) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uint64) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}
