package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray is a custom type to handle string arrays in different databases
type StringArray []string

func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), sa)
	case []byte:
		return json.Unmarshal(v, sa)
	}

	return nil
}

func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	return json.Marshal(sa)
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
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
