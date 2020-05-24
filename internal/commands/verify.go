package commands

import (
	"context"
	"encoding/json"
	"moderari/internal/db"
	"moderari/internal/embeds"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "verify",
		Aliases:     []string{"trust", "allow"},
		Usage:       "<member>",
		Description: "Bypasses user verification by adding the member role.",
		Category:    modCategory,
		ArgTransformers: []gommand.ArgTransformer{
			{
				Function: gommand.MemberTransformer,
			},
		},
		PermissionValidators: []gommand.PermissionValidator{gommand.MANAGE_ROLES},
		Function:             verify,
	})
}

func verify(ctx *gommand.Context) error {
	guildString, err := db.Client.Get("guild:" + ctx.Message.GuildID.String()).Result()
	if err != nil {
		return err
	}
	guild := db.GuildModel{}
	_ = json.Unmarshal([]byte(guildString), &guild)

	member := ctx.Args[0].(*disgord.Member)
	if err := ctx.Session.AddGuildMemberRole(
		context.Background(),
		ctx.Message.GuildID,
		member.User.ID,
		disgord.ParseSnowflakeString(guild.MemberRole),
	); err != nil {
		return err
	}

	_, err = ctx.Reply(embeds.Info(
		"Verified",
		member.User.Username+" has been verified.",
		"",
	))
	return err
}
