package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
}

type AuthService interface {
	Register(ctx context.Context, user *model.User) error
}

type AuthHandler interface {
	Register(c *gin.Context)
}
