package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type UserService interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type UserHandler interface {
	GetUser(c *gin.Context)
	GetCurrentUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
}
