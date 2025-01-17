package models

import "github.com/google/uuid"

type User struct {
	UserID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email             string    `json:"email"`
	Name              string    `json:"name"`
	WalletAddress     string    `json:"wallet_address"`
	ProfilePicture    string    `json:"profilePictureUrl"`
	CoverPicture      string    `json:"coverPictureUrl"`
	Location          string    `json:"location"`
	FacebookId        string    `json:"facebook_id"`
	InstagramId       string    `json:"instagram_id"`
	TwitterId         string    `json:"twitter_id"`
	DiscordId         string    `json:"discord_id"`
	TelegramId        string    `json:"telegram_id"`
	UserType          string    `json:"user_type"`
	Bio               string    `json:"bio"`
	InstagramVerified bool      `json:"instagramVerified"`
	FacebookVerified  bool      `json:"facebookVerified"`
	TwitterVerified   bool      `json:"twitterVerified"`
	DiscordVerified   bool      `json:"discordVerified"`
	TelegramVerified  bool      `json:"telegramVerified"`
	Plan              string    `json:"plan"`
	ProfileCreated    bool      `json:"profileCreated"`
}
