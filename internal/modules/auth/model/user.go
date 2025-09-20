package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email            string         `gorm:"unique;index" json:"email"`
	Password         string         `json:"password"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Role             string         `gorm:"default:'free'" json:"role"`
	ProfileCompleted bool           `gorm:"default:false" json:"profile_completed"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
