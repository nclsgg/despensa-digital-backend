package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditTransaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WalletID    uuid.UUID `gorm:"type:uuid;index;not null" json:"wallet_id"`
	UserID      uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Amount      int       `gorm:"not null" json:"amount"`
	Type        string    `gorm:"type:varchar(16);not null" json:"type"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (t *CreditTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
