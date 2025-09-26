package model

import (
	"time"

	"github.com/google/uuid"
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
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (u *PantryUser) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
