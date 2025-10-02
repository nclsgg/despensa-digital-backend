package model

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email            string         `gorm:"unique;index" json:"email"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Role             string         `gorm:"default:'free'" json:"role"`
	ProfileCompleted bool           `gorm:"default:false" json:"profile_completed"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	__logParams := map[string]any{"u": u, "tx": tx}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*User.BeforeCreate"), zap.Any("result", err), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*User.BeforeCreate"), zap.Any("params", __logParams))
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
