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
		Name:                 "nick",
		Aliases:              []string{"nickname"},
		Description:          "Changes a member's nickname.",
		Usage:                "<member> <nickname>",
		Category:             modCategory,
		PermissionValidators: []gommand.PermissionValidator{gommand.MANAGE_NICKNAMES},
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: false,
				Function: gommand.MemberTransformer,
			},
			{
				Optional:  false,
				Remainder: true,
				Function:  gommand.StringTransformer,
			},
		},
		Function: nick,
	})
}

func nick(ctx *gommand.Context) error {
	member := ctx.Args[0].(*disgord.Member)
	nickname := ctx.Args[1].(string)

	oldNickname := member.Nick
	if err := member.UpdateNick(context.Background(), ctx.Session, nickname); err != nil {
		// TODO: Fix httd.ErrREST (Unknown Guild). Seems to be a Disgord issue.
		return err
	}

	_, err := ctx.Reply(embeds.Info(
		"Nickname Updated",
		fmt.Sprintf("%s ~> %s", oldNickname, member.Nick), ""),
	)
	return err
}
