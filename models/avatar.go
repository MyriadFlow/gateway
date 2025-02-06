package models

import (
	"time"

	"github.com/google/uuid"
)

type Avatar struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AvatarID         string
	URL              string
	UserID           string
	BrandID          string
	PhygitalID       string
	AvatarVoice      string
	AvatarName       string
	AvatarDescription string
	ChaintypeID      uuid.UUID `gorm:"type:uuid"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}