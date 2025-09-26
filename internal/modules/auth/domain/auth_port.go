package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
}

type AuthService interface {
	// Core methods (kept for OAuth)
	HashPassword(password string) (string, error)
	GenerateAccessToken(user *model.User) (string, error)

	// OAuth methods
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUserOAuth(ctx context.Context, user *model.User) error
	CompleteProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string) error
}

// AuthHandler interface removed - using OAuth only
