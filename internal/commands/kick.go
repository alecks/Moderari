package commands

import (
	"context"
	"fmt"
	"moderari/internal/embeds"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:                 "kick",
		Description:          "Kicks a member.",
		Usage:                "<member> [reason]",
		Category:             modCategory,
		PermissionValidators: []gommand.PermissionValidator{gommand.KICK_MEMBERS},
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: false,
				Function: gommand.MemberTransformer,
			},
			{
				Optional:  true,
				Remainder: true,
				Function:  gommand.StringTransformer,
			},
		},
		Function: kick,
	})
}

func kick(ctx *gommand.Context) error {
	member := ctx.Args[0].(*disgord.Member)
	reason := "Kicked by " + ctx.Message.Author.Username
	if ctx.Args[1] != nil {
		reason = ctx.Args[1].(string)
	}

	if err := ctx.Session.KickMember(
		context.Background(),
		ctx.Message.GuildID,
		member.User.ID,
		reason,
	); err != nil {
		return err
	}

	_, err := ctx.Reply(embeds.Info(
		"Kicked",
		fmt.Sprintf("%s has been kicked by %s.", member.User.Username, ctx.Message.Author.Username), ""),
	)
	return err
}
