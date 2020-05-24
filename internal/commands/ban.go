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
		Name:                 "ban",
		Description:          "Bans a member.",
		Usage:                "<member> [days of messages to delete] [reason]",
		Category:             modCategory,
		PermissionValidators: []gommand.PermissionValidator{gommand.BAN_MEMBERS},
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: false,
				Function: gommand.MemberTransformer,
			},
			{
				Optional: true,
				Function: gommand.IntTransformer,
			},
			{
				Optional:  true,
				Remainder: true,
				Function:  gommand.StringTransformer,
			},
		},
		Function: ban,
	})
}

func ban(ctx *gommand.Context) error {
	member := ctx.Args[0].(*disgord.Member)
	deleteMessageDays := 0
	reason := "Banned by " + ctx.Message.Author.Username

	argDeleteMessageDays := ctx.Args[1]
	if ctx.Args[1] != nil {
		deleteMessageDays = argDeleteMessageDays.(int)
	}
	argReason := ctx.Args[2]
	if ctx.Args[2] != nil {
		reason = argReason.(string)
	}

	if err := ctx.Session.BanMember(
		context.Background(),
		ctx.Message.GuildID,
		member.User.ID,
		&disgord.BanMemberParams{
			Reason:            reason,
			DeleteMessageDays: deleteMessageDays,
		},
	); err != nil {
		return err
	}

	_, err := ctx.Reply(embeds.Info(
		"Banned",
		fmt.Sprintf("**%s** has been banned by **%s**.", member.User.Username, ctx.Message.Author.Username), ""),
	)
	return err
}
