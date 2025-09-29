package models

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeObject   Type = "object"
	TypeStandard Type = "standard"
)

type Configurations struct {
	ID        uuid.UUID `gorm:"primarykey"`
	ClientID  string
	Name      string `gorm:"size:100;uniqueIndex:idx_name_version"`
	Type      Type
	Schema    string `gorm:"type:TEXT;check:json_valid(schema)"`
	Input     string `gorm:"type:TEXT;check:json_valid(input)"`
	Version   int    `gorm:"uniqueIndex:idx_name_version"`
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	IsActive  int
}
