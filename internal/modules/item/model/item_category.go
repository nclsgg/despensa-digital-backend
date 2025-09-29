package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"go.uber.org/zap"
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
	__logParams := map[string]any{"i": i, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ItemCategory.ApplyUpdate"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ItemCategory.ApplyUpdate"), zap.Any("params", __logParams))
	if input.Name != nil {
		i.Name = *input.Name
	}
	if input.Color != nil {
		i.Color = *input.Color
	}
}
