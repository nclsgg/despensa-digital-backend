package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/item/dto"
	"gorm.io/gorm"
)

// models/item.go
type Item struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PantryID     uuid.UUID      `gorm:"type:uuid;not null" json:"pantry_id"`
	AddedBy      uuid.UUID      `gorm:"type:uuid;not null" json:"added_by"`
	CategoryID   *uuid.UUID     `gorm:"type:uuid" json:"category_id"`
	Name         string         `gorm:"not null" json:"name"`
	Quantity     float64        `gorm:"not null" json:"quantity"`
	TotalPrice   float64        `gorm:"->;type:numeric" json:"total_price"`
	PricePerUnit float64        `gorm:"type:numeric;not null" json:"price_per_unit"`
	Unit         string         `gorm:"not null" json:"unit"`
	ExpiresAt    *time.Time     `gorm:"type:timestamp" json:"expires_at"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (i *Item) ApplyUpdate(input dto.UpdateItemDTO) {
	if input.Name != nil {
		i.Name = *input.Name
	}
	if input.Quantity != nil {
		i.Quantity = *input.Quantity
	}
	if input.PricePerUnit != nil {
		i.PricePerUnit = *input.PricePerUnit
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
