package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique_index" json:"email"`
	Password  string         `json:"password"`
	Name      string         `gorm:"not null" json:"name"`
	Role      string         `gorm:"default:'free'" json:"role"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
