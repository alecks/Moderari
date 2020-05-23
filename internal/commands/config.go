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
		Function:    configCmd,
	})
}

func configCmd(ctx *gommand.Context) error {
	menu := gommand.NewEmbedMenu(embeds.Info("Configuration", "", ""), ctx)

	prefixMenu := menu.NewChildMenu(
		&gommand.ChildMenuOptions{
			Embed: embeds.Info(
				"Prefix",
				"Enter a new prefix.",
				"",
			),
			Button: &gommand.MenuButton{
				Emoji:       "🔠",
				Name:        "Prefix",
				Description: "The prefix changes how you execute commands. If it's `|`, you'd use `|ping`.",
			},
		},
	)
	{
		prefixMenu.AddBackButton()
		prefixMenu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"User Prefix",
					"Enter a new prefix.",
					"",
				),
				Button: &gommand.MenuButton{
					Emoji:       "😛",
					Name:        "User",
					Description: "This changes the prefix for *you*, not the entire server. It'll apply everywhere.",
				},
			},
		).AddBackButton()
		prefixMenu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Server Prefix",
					"Enter a new prefix.",
					"",
				),
				Button: &gommand.MenuButton{
					Emoji:       "👨‍👩‍👧‍👦",
					Name:        "Server",
					Description: "This changes the prefix for everyone in the server. It'll only apply here.",
				},
			},
		).AddBackButton()
	}

	menu.NewChildMenu(
		&gommand.ChildMenuOptions{
			Embed: embeds.Info(
				"Ban Threshold",
				"Enter a new ban threshold.",
				"",
			),
			Button: &gommand.MenuButton{
				Emoji:       "🔨",
				Name:        "Ban Threshold",
				Description: "The ban threshold is the amount of warnings at which a member is banned.",
			},
		},
	).AddBackButton()

	_ = ctx.DisplayEmbedMenu(menu)
	return nil
}
