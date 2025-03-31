package service

import (
	"context"
	"log"

	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/model"
)

type authService struct {
	repo domain.AuthRepository
}

func NewAuthService(repo domain.AuthRepository) domain.AuthService {
	return &authService{repo}
}

func (s *authService) Register(ctx context.Context, user *model.User) error {
	//TODO - validate user and hash password
	log.Println("user: ", user)
	return s.repo.CreateUser(ctx, user)
}
