package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditWallet struct {
	ID           uuid.UUID           `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID           `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Balance      int                 `gorm:"not null" json:"balance"`
	Transactions []CreditTransaction `gorm:"foreignKey:WalletID" json:"transactions"`
	CreatedAt    time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
}

func (w *CreditWallet) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	if w.Balance == 0 {
		w.Balance = 10
	}
	return nil
}
