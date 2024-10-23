package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string    `gorm:"uniqueIndex" json:"username"`
	Password     string    `json:"-"`
	Roles        string    `gorm:"type:text" json:"-"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RefreshToken string    `gorm:"type:text" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
