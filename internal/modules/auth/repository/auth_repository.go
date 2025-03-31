package repository

import (
	"context"

	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/model"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
