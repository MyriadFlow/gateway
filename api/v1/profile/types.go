package profile

type GetProfilePayload struct {
	Name              string `json:"name,omitempty"`
	WalletAddress     string `json:"walletAddress"`
	ProfilePictureUrl string `json:"profilePictureUrl,omitempty"`
	CoverPictureUrl   string `json:"coverPictureUrl,omitempty"`
	Location          string `json:"location,omitempty"`
	FacebookId        string `json:"facebook_id,omitempty"`
	InstagramId       string `json:"instagram_id,omitempty"`
	TwitterId         string `json:"twitter_id,omitempty"`
	DiscordId         string `json:"discord_id,omitempty"`
	TelegramId        string `json:"telegram_id,omitempty"`
	Email             string `json:"email,omitempty"`
	Bio               string `json:"bio,omitempty"`
	InstagramVerified bool   `json:"instagramVerified,omitempty"`
	FacebookVerified  bool   `json:"facebookVerified,omitempty"`
	TwitterVerified   bool   `json:"twitterVerified,omitempty"`
	DiscordVerified   bool   `json:"discordVerified,omitempty"`
	TelegramVerified  bool   `json:"telegramVerified,omitempty"`
	Plan              string `json:"plan,omitempty"`
}

type verifySocialPayload struct {
	SocialName string `json:"socialName"`
}

type createProfilePayload struct {
	Name              string `json:"name"`
	Bio               string `json:"bio"`
	Email             string `json:"email"`
	Location          string `json:"location"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	CoverPictureUrl   string `json:"coverPictureUrl"`
	InstagramId       string `json:"instagram_id,omitempty"`
	TwitterId         string `json:"twitter_id,omitempty"`
	DiscordId         string `json:"discord_id,omitempty"`
}

type updateProfilePayload struct {
	Name              string `json:"name"`
	Bio               string `json:"bio"`
	Email             string `json:"email"`
	Location          string `json:"location"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	CoverPictureUrl   string `json:"coverPictureUrl"`
	InstagramId       string `json:"instagram_id,omitempty"`
	TwitterId         string `json:"twitter_id,omitempty"`
	DiscordId         string `json:"discord_id,omitempty"`
}
