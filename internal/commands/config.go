package commands

import (
	"context"
	"encoding/json"
	"moderari/internal/config"
	"moderari/internal/db"
	"moderari/internal/embeds"
	"strconv"

	"github.com/andersfylling/disgord"

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

// TODO: Message waiters don't stop when you go back to the parent page.
func configCmd(ctx *gommand.Context) error {
	guildKey := "guild:" + ctx.Message.GuildID.String()
	userKey := "user:" + ctx.Message.Author.ID.String()

	// TODO: This needs improved, uh, a lot.
	guildString, err := db.Client.Get(guildKey).Result()
	guild := db.GuildModel{
		BanThreshold: 0,
		Prefix:       config.C.Prefix,
	}
	if err != db.Nil {
		if err != nil {
			return err
		}
		_ = json.Unmarshal([]byte(guildString), &guild)
	}
	userString, err := db.Client.Get(userKey).Result()
	user := db.UserModel{Warns: map[string]map[string]db.WarnModel{}}
	if err != db.Nil {
		if err != nil {
			return err
		}
		_ = json.Unmarshal([]byte(userString), &user)
	}

	menu := gommand.NewEmbedMenu(embeds.Info("Configuration", "", ""), ctx)

	prefixMenu := menu.NewChildMenu(
		&gommand.ChildMenuOptions{
			Embed: embeds.Info(
				"Prefix",
				"There are two types of prefixes. Please select one.",
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
					"Current: "+user.Prefix,
				),
				Button: &gommand.MenuButton{
					Emoji:       "😛",
					Name:        "User",
					Description: "This changes the prefix for *you*, not the entire server. It'll apply everywhere.",
				},
				AfterAction: func() {
					// TODO: Possibly implement a get/set interface?
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)

					user.Prefix = res.Content
					go func() {
						marshalled, _ := json.Marshal(user)
						db.Client.Set(userKey, marshalled, 0)
					}()

					newPrefixMessage(res.Content, ctx)
				},
			},
		).AddBackButton()
		prefixMenu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Server Prefix",
					"Enter a new prefix.",
					"Current: "+guild.Prefix,
				),
				Button: &gommand.MenuButton{
					Emoji:       "👨‍👩‍👧‍👦",
					Name:        "Server",
					Description: "This changes the prefix for everyone in the server. It'll only apply here.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)

					guild.Prefix = res.Content
					go func() {
						marshalled, _ := json.Marshal(guild)
						db.Client.Set(guildKey, marshalled, 0)
					}()

					newPrefixMessage(res.Content, ctx)
				},
			},
		).AddBackButton()
	}

	menu.NewChildMenu(
		&gommand.ChildMenuOptions{
			Embed: embeds.Info(
				"Ban Threshold",
				"Enter a new ban threshold.",
				"Current: "+strconv.Itoa(guild.BanThreshold),
			),
			Button: &gommand.MenuButton{
				Emoji:       "🔨",
				Name:        "Ban Threshold",
				Description: "The ban threshold is the amount of warnings at which a member is banned.",
			},
			AfterAction: func() {
				res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
					return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
				})
				go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)

				newThreshold, err := strconv.Atoi(res.Content)
				if err != nil {
					_, _ = ctx.Reply("Invalid number.")
					return
				}
				guild.BanThreshold = newThreshold
				go func() {
					marshalled, _ := json.Marshal(guild)
					db.Client.Set(guildKey, marshalled, 0)
				}()

				_, _ = ctx.Reply("Ban threshold updated. Members will be banned once they reach", res.Content, "warnings.")
			},
		},
	).AddBackButton()

	menu.NewChildMenu(
		&gommand.ChildMenuOptions{
			Embed: embeds.Info(
				"Gotcha",
				"Gotcha has been toggled.",
				"",
			),
			Button: &gommand.MenuButton{
				Emoji:       "👀",
				Name:        "Gotcha",
				Description: "[Gotcha](https://github.com/fjah/gotcha) stops users on the same network from raiding your server.",
			},
			AfterAction: func() {
				guild.Gotcha = !guild.Gotcha
				message := "disabled"
				if guild.Gotcha {
					message = "enabled"
				}
				go func() {
					marshalled, _ := json.Marshal(guild)
					db.Client.Set(guildKey, marshalled, 0)
				}()

				_, _ = ctx.Reply("Gotcha is now", message+".", "Members will receive a verification link when they join.")
			},
		},
	).AddBackButton()

	_ = ctx.DisplayEmbedMenu(menu)
	return nil
}

func newPrefixMessage(prefix string, ctx *gommand.Context) {
	_, _ = ctx.Reply("Prefix updated. You can now use commands like so: `" + prefix + "config`.")
}
