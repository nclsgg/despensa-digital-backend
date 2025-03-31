package model

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Email     string `gorm:"unique_index" json:"email"`
	Password  string `json:"password"`
	Name      string `gorm:"not null" json:"name"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt string `gorm:"index" json:"deleted_at"`
}
