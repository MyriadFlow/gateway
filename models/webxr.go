package models

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type WebXR struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ChainID     uuid.UUID `gorm:"type:uuid"`
	Image360    string
	Video360    string
	PhygitalID  string
	BrandID     string
	AvatarID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}