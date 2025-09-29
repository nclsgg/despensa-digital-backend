package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StringArray is a custom type to handle string arrays in different databases
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

type Profile struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID              uuid.UUID      `gorm:"type:uuid;not null;unique" json:"user_id"`
	MonthlyIncome       float64        `gorm:"type:numeric" json:"monthly_income"`
	PreferredBudget     float64        `gorm:"type:numeric" json:"preferred_budget"`
	HouseholdSize       int            `gorm:"default:1" json:"household_size"`
	DietaryRestrictions StringArray    `gorm:"type:text" json:"dietary_restrictions"`
	PreferredBrands     StringArray    `gorm:"type:text" json:"preferred_brands"`
	ShoppingFrequency   string         `gorm:"default:'weekly'" json:"shopping_frequency"` // weekly, biweekly, monthly
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	__logParams := map[string]any{"p": p, "tx": tx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*Profile.BeforeCreate"), zap.Any("result", err), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*Profile.BeforeCreate"), zap.Any("params", __logParams))
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
