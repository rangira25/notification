package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FullName     string         `gorm:"type:varchar(120);not null" json:"full_name"`
	Email        string         `gorm:"type:citext;uniqueIndex;not null" json:"email"`
	Phone        *string        `gorm:"type:varchar(32);index" json:"phone,omitempty"`
	PasswordHash *string        `gorm:"type:text" json:"-"`
	Role         string         `gorm:"type:varchar(32);default:user;index" json:"role"`
	Status       string         `gorm:"type:varchar(16);default:active;index" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
