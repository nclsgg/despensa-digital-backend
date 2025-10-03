package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StringArray []string

func (sa *StringArray) Scan(value interface{}) (result0 error) {
	__logParams := map[string]any{"sa": sa, "value": value}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*StringArray.Scan"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*StringArray.Scan"), zap.Any("params", __logParams))
	if value == nil {
		*sa = nil
		result0 = nil
		return
	}

	switch v := value.(type) {
	case string:
		result0 = json.Unmarshal([]byte(v), sa)
		return
	case []byte:
		result0 = json.Unmarshal(v, sa)
		return
	}
	result0 = nil
	return
}

func (sa StringArray) Value() (result0 driver.Value, result1 error) {
	__logParams := map[string]any{"sa": sa}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "StringArray.Value"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "StringArray.Value"), zap.Any("params", __logParams))
	if sa == nil {
		result0 = nil
		result1 = nil
		return
	}
	result0, result1 = json.Marshal(sa)
	return
}

type ShoppingList struct {
	ID                  uuid.UUID          `gorm:"type:uuid;primary_key" json:"id"`
	UserID              uuid.UUID          `gorm:"type:uuid;not null;index:idx_shopping_list_user,priority:1" json:"user_id"`
	PantryID            *uuid.UUID         `gorm:"type:uuid;index" json:"pantry_id"`
	Name                string             `gorm:"not null" json:"name"`
	Status              string             `gorm:"default:'pending';index" json:"status"` // pending, completed, cancelled
	TotalBudget         float64            `gorm:"type:numeric" json:"total_budget"`
	EstimatedCost       float64            `gorm:"type:numeric" json:"estimated_cost"`
	ActualCost          float64            `gorm:"type:numeric" json:"actual_cost"`
	GeneratedBy         string             `gorm:"default:'manual'" json:"generated_by"` // manual, ai
	HouseholdSize       int                `gorm:"default:1" json:"household_size"`
	MonthlyIncome       float64            `gorm:"type:numeric" json:"monthly_income"`
	DietaryRestrictions StringArray        `gorm:"type:text" json:"dietary_restrictions"`
	Items               []ShoppingListItem `gorm:"foreignKey:ShoppingListID" json:"items"`
	CreatedAt           time.Time          `gorm:"autoCreateTime;index:idx_shopping_list_user,priority:2" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt     `gorm:"index" json:"deleted_at"`
}

type ShoppingListItem struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ShoppingListID uuid.UUID      `gorm:"type:uuid;not null;index:idx_shopping_item_list,priority:1" json:"shopping_list_id"`
	Name           string         `gorm:"not null" json:"name"`
	Quantity       float64        `gorm:"not null" json:"quantity"`
	Unit           string         `gorm:"not null" json:"unit"`
	EstimatedPrice float64        `gorm:"type:numeric" json:"estimated_price"`
	ActualPrice    float64        `gorm:"type:numeric" json:"actual_price"`
	Category       string         `json:"category"`
	Priority       int            `gorm:"default:3" json:"priority"` // 1=high, 2=medium, 3=low
	Purchased      bool           `gorm:"default:false;index:idx_shopping_item_list,priority:2" json:"purchased"`
	Source         string         `json:"source"` // pantry_history, ai_suggestion, manual
	PantryItemID   *uuid.UUID     `gorm:"type:uuid;index" json:"pantry_item_id"`
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
