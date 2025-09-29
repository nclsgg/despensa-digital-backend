package model

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ShoppingList struct {
	ID            uuid.UUID          `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID          `gorm:"type:uuid;not null" json:"user_id"`
	PantryID      *uuid.UUID         `gorm:"type:uuid;index" json:"pantry_id"`
	Name          string             `gorm:"not null" json:"name"`
	Status        string             `gorm:"default:'pending'" json:"status"` // pending, completed, cancelled
	TotalBudget   float64            `gorm:"type:numeric" json:"total_budget"`
	EstimatedCost float64            `gorm:"type:numeric" json:"estimated_cost"`
	ActualCost    float64            `gorm:"type:numeric" json:"actual_cost"`
	GeneratedBy   string             `gorm:"default:'manual'" json:"generated_by"` // manual, ai
	Items         []ShoppingListItem `gorm:"foreignKey:ShoppingListID" json:"items"`
	CreatedAt     time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt     `gorm:"index" json:"deleted_at"`
}

type ShoppingListItem struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ShoppingListID uuid.UUID      `gorm:"type:uuid;not null" json:"shopping_list_id"`
	Name           string         `gorm:"not null" json:"name"`
	Quantity       float64        `gorm:"not null" json:"quantity"`
	Unit           string         `gorm:"not null" json:"unit"`
	EstimatedPrice float64        `gorm:"type:numeric" json:"estimated_price"`
	ActualPrice    float64        `gorm:"type:numeric" json:"actual_price"`
	Category       string         `json:"category"`
	Brand          string         `json:"brand"`
	Priority       int            `gorm:"default:3" json:"priority"` // 1=high, 2=medium, 3=low
	Purchased      bool           `gorm:"default:false" json:"purchased"`
	Notes          string         `json:"notes"`
	Source         string         `json:"source"` // pantry_history, ai_suggestion, manual
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (s *ShoppingList) BeforeCreate(tx *gorm.DB) (err error) {
	__logParams := map[string]any{"s": s, "tx": tx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingList.BeforeCreate"), zap.Any("result", err), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingList.BeforeCreate"), zap.Any("params", __logParams))
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

func (s *ShoppingListItem) BeforeCreate(tx *gorm.DB) (err error) {
	__logParams := map[string]any{"s": s, "tx": tx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ShoppingListItem.BeforeCreate"), zap.Any("result", err), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ShoppingListItem.BeforeCreate"), zap.Any("params", __logParams))
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
