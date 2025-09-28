package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           uuid.UUID `gorm:"primarykey"`
	Username     string    `gorm:"uniqueIndex;size:100"`
	PasswordHash string
	Role         Role `gorm:"size:20;default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
