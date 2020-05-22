package commands

import (
	"moderari/internal/embeds"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "config",
		Description: "Configures the bot.",
		Category:    utilCategory,
		Function:    config,
	})
}

func config(ctx *gommand.Context) error {
	menu := gommand.NewEmbedMenu(embeds.Info("Configuration", "", ""), ctx)

	menu.NewChildMenu(
		embeds.Info(
			"Prefix",
			"Enter a new prefix.",
			"",
		),
		gommand.MenuButton{
			Emoji:       "🔠",
			Name:        "Prefix",
			Description: "The prefix changes how you execute commands. If it's `|`, you'd use `|ping`.",
		},
	).AddBackButton()
	menu.NewChildMenu(
		embeds.Info(
			"Warn Threshold",
			"Enter a new warn threshold.",
			"",
		),
		gommand.MenuButton{
			Emoji:       "⚠️",
			Name:        "Warn Threshold",
			Description: "The warn threshold is the amount of warnings at which a member is banned.",
		},
	).AddBackButton()

	_ = ctx.DisplayEmbedMenu(menu)
	return nil
}
