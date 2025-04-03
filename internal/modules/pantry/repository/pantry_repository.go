package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/pantry/model"
	"gorm.io/gorm"
)

type pantryRepository struct {
	db *gorm.DB
}

func NewPantryRepository(db *gorm.DB) domain.PantryRepository {
	return &pantryRepository{db}
}

func (r *pantryRepository) Create(ctx context.Context, pantry *model.Pantry) (*model.Pantry, error) {
	err := r.db.WithContext(ctx).Create(pantry).Error
	return pantry, err
}

func (r *pantryRepository) Delete(ctx context.Context, pantryID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Pantry{}, "id = ?", pantryID).Error
}

func (r *pantryRepository) Update(ctx context.Context, pantry *model.Pantry) error {
	return r.db.WithContext(ctx).Save(pantry).Error
}

func (r *pantryRepository) GetByID(ctx context.Context, pantryID uuid.UUID) (*model.Pantry, error) {
	var pantry model.Pantry
	err := r.db.WithContext(ctx).First(&pantry, "id = ?", pantryID).Error
	if err != nil {
		return nil, err
	}
	return &pantry, nil
}

func (r *pantryRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*model.Pantry, error) {
	var pantries []*model.Pantry
	err := r.db.WithContext(ctx).
		Joins("JOIN pantry_users ON pantries.id = pantry_users.pantry_id").
		Where("pantry_users.user_id = ? AND pantry_users.deleted_at IS NULL", userID).
		Find(&pantries).Error
	return pantries, err
}

func (r *pantryRepository) IsUserInPantry(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PantryUser{}).
		Where("pantry_id = ? AND user_id = ? AND deleted_at IS NULL", pantryID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *pantryRepository) IsUserOwner(ctx context.Context, pantryID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PantryUser{}).
		Where("pantry_id = ? AND user_id = ? AND role = ? AND deleted_at IS NULL", pantryID, userID, "owner").
		Count(&count).Error
	return count > 0, err
}

func (r *pantryRepository) AddUserToPantry(ctx context.Context, pantryUser *model.PantryUser) error {
	return r.db.WithContext(ctx).Create(pantryUser).Error
}

func (r *pantryRepository) RemoveUserFromPantry(ctx context.Context, pantryID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("pantry_id = ? AND user_id = ?", pantryID, userID).
		Delete(&model.PantryUser{}).Error
}

func (r *pantryRepository) ListUsersInPantry(ctx context.Context, pantryID uuid.UUID) ([]*model.PantryUserInfo, error) {
	var users []*model.PantryUserInfo
	err := r.db.WithContext(ctx).
		Table("pantry_users").
		Select("users.id as user_id, users.email, pantry_users.role").
		Joins("JOIN users ON users.id = pantry_users.user_id").
		Where("pantry_users.pantry_id = ?", pantryID).
		Scan(&users).Error

	return users, err
}
