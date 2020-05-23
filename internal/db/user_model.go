package db

// UserModel is the database model for Discord users.
type UserModel struct {
	Warns  map[string]map[string]WarnModel `json:"warns"`
	Prefix string                          `json:"prefix"`
}

// WarnModel the database model for a warning/strike.
type WarnModel struct {
	Reason    string `json:"reason"`
	Moderator string `json:"moderator"`
	Time      int64  `json:"time"`
	Severity  uint64 `json:"severity"`
}
