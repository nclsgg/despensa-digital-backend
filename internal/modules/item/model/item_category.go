package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"gorm.io/gorm"
)

type ItemCategory struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PantryID  uuid.UUID      `gorm:"type:uuid" json:"pantry_id"`
	AddedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"added_by"`
	Name      string         `gorm:"not null" json:"name"`
	Color     string         `gorm:"not null" json:"color"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (i *ItemCategory) ApplyUpdate(input dto.UpdateItemCategoryDTO) {
	if input.Name != nil {
		i.Name = *input.Name
	}
	if input.Color != nil {
		i.Color = *input.Color
	}
}
