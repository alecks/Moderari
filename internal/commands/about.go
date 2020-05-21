package commands

import (
	"moderari/internal/embeds"
	"strings"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "about",
		Description: "Retrieves bot information.",
		Category:    infoCategory,
		Function:    about,
	})
}

func about(ctx *gommand.Context) error {
	latency, _ := ctx.Session.AvgHeartbeatLatency()

	_, err := ctx.Reply(embeds.Info(
		"About "+strings.Title(ctx.BotUser.Username),
		"This is an instance of **Moderari**, an open-source project.",
		latency.String(),
		embeds.Field("built with", "[andersfylling/disgord](https://github.com/andersfylling/disgord)", true),
		embeds.Field("using", "[auttaja/gommand](https://github.com/auttaja/gommand)", true),
	))
	return err
}
