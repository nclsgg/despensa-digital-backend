package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// models/item.go
type Item struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PantryID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_item_pantry,priority:1" json:"pantry_id"`
	AddedBy       uuid.UUID      `gorm:"type:uuid;not null" json:"added_by"`
	CategoryID    *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	Name          string         `gorm:"not null;index:idx_item_name" json:"name"`
	Quantity      float64        `gorm:"not null" json:"quantity"`
	TotalPrice    float64        `gorm:"->;type:numeric" json:"total_price"`
	PricePerUnit  float64        `gorm:"type:numeric;not null" json:"price_per_unit"`
	PriceQuantity float64        `gorm:"type:numeric;default:1" json:"price_quantity"`
	Unit          string         `gorm:"not null" json:"unit"`
	ExpiresAt     *time.Time     `gorm:"type:timestamp;index" json:"expires_at"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index:idx_item_pantry,priority:2" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (i *Item) ApplyUpdate(input dto.UpdateItemDTO) {
	__logParams := map[string]any{"i": i, "input": input}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*Item.ApplyUpdate"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*Item.ApplyUpdate"), zap.Any("params", __logParams))
	if input.Name != nil {
		i.Name = *input.Name
	}
	if input.Quantity != nil {
		i.Quantity = *input.Quantity
	}
	if input.PriceQuantity != nil {
		i.PriceQuantity = *input.PriceQuantity
	}
	if input.Unit != nil {
		i.Unit = *input.Unit
	}
	if input.CategoryID != nil {
		parsedUUID := uuid.MustParse(*input.CategoryID)
		i.CategoryID = &parsedUUID
	}
	if input.ExpiresAt != "" {
		layout := "2006-01-02"
		parsedTime, err := time.Parse(layout, input.ExpiresAt)
		if err == nil {
			i.ExpiresAt = &parsedTime
		}
	}
}
