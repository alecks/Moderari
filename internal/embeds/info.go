package embeds

import (
	"moderari/internal/config"

	"github.com/andersfylling/disgord"
)

// Info instantiates an informational embed.
func Info(title, description, footer string, fields ...*disgord.EmbedField) *disgord.Embed {
	return &disgord.Embed{
		Title:       title,
		Description: description,
		Footer:      &disgord.EmbedFooter{Text: footer},
		Color:       config.C.Colors.Info,
		Fields:      fields,
	}
}
