package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/model"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *model.Profile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error)
	Update(ctx context.Context, profile *model.Profile) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProfileService interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, input dto.CreateProfileDTO) (*dto.ProfileResponseDTO, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponseDTO, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error)
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}
