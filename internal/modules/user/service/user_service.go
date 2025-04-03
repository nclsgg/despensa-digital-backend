package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo}
}

func (s *userService) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetAllUsers(ctx)
}
