package model

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Pantry struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type PantryUser struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PantryID  uuid.UUID      `gorm:"type:uuid;not null" json:"pantry_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Role      string         `gorm:"default:'member'" json:"role"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type PantryUserInfo struct {
	ID        uuid.UUID `json:"id"`
	PantryID  uuid.UUID `json:"pantry_id"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
}

type PantryWithItemCount struct {
	Pantry    *Pantry `json:"pantry"`
	ItemCount int     `json:"item_count"`
}

func (u *Pantry) BeforeCreate(tx *gorm.DB) (err error) {
	__logParams := map[string]any{"u": u, "tx": tx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*Pantry.BeforeCreate"), zap.Any("result", err), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*Pantry.BeforeCreate"), zap.Any("params", __logParams))
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (u *PantryUser) BeforeCreate(tx *gorm.DB) (err error) {
	__logParams := map[string]any{"u": u, "tx": tx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*PantryUser.BeforeCreate"), zap.Any("result", err), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*PantryUser.BeforeCreate"), zap.Any("params", __logParams))
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
