package commands

import (
	"moderari/internal/embeds"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "userinfo",
		Aliases:     []string{"ui", "profile", "user", "uinfo"},
		Description: "Shows a user's profile.",
		Category:    infoCategory,
		Function:    userInfo,
	})
}

func userInfo(ctx *gommand.Context) error {
	_, err := ctx.Reply(embeds.Info(
		"TODO",
		"",
		"",
	))
	return err
}
