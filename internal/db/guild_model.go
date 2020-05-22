package db

// GuildModel is the database model for Discord guilds.
type GuildModel struct {
	BanThreshold int `json:"ban_threshold"`
}
