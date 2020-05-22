package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"moderari/internal/db"
	"moderari/internal/embeds"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
	"github.com/google/uuid"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:                 "warn",
		Description:          "Warns a member.",
		Aliases:              []string{"strike", "warning"},
		Usage:                "<member> [severity] [reason]",
		Category:             modCategory,
		PermissionValidators: []gommand.PermissionValidator{gommand.MANAGE_MESSAGES},
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: false,
				Function: gommand.MemberTransformer,
			},
			{
				Optional: true,
				Function: gommand.UIntTransformer,
			},
			{
				Optional:  true,
				Remainder: true,
				Function:  gommand.StringTransformer,
			},
		},
		Function: warn,
	})
}

func warn(ctx *gommand.Context) error {
	member := ctx.Args[0].(*disgord.Member)

	// TODO: Improve this.
	var severity uint64 = 0
	argSeverity := ctx.Args[1]
	if argSeverity != nil {
		severity = argSeverity.(uint64)
	}
	reason := "Missing reason."
	argReason := ctx.Args[2]
	if argReason != nil {
		reason = argReason.(string)
	}

	id := uuid.New().String()
	go func() {
		key := "user:" + member.User.ID.String()
		oldString, err := db.Client.Get(key).Result()
		if err == db.Nil {
			oldString = "{}"
		}

		old := db.UserModel{}
		_ = json.Unmarshal([]byte(oldString), &old)

		warning := db.WarnModel{
			Reason:    reason,
			Severity:  severity,
			Moderator: ctx.Message.Author.ID.String(),
			Time:      time.Now().Unix(),
		}

		if old.Warns == nil {
			old.Warns = map[string]map[string]db.WarnModel{}
		}
		guildID := ctx.Message.GuildID.String()
		if old.Warns[guildID] == nil {
			old.Warns[guildID] = map[string]db.WarnModel{}
		}

		old.Warns[guildID][id] = warning

		go func() {
			guildString, err := db.Client.Get("guild:" + ctx.Message.GuildID.String()).Result()
			if err != nil {
				return
			}

			guild := db.GuildModel{}
			_ = json.Unmarshal([]byte(guildString), &guild)
			if guild.BanThreshold <= 0 {
				return
			}

			if len(old.Warns) >= guild.BanThreshold {
				_ = ctx.Session.BanMember(
					context.Background(),
					ctx.Message.GuildID,
					member.User.ID,
					&disgord.BanMemberParams{Reason: reason},
				)
			}
		}()

		bytes, _ := json.Marshal(old)
		db.Client.Set(key, bytes, 0)
	}()

	_, err := ctx.Reply(embeds.Info(
		"Warned",
		fmt.Sprintf("**%s** has been warned by **%s** for \"%s\".", member.User.Username, ctx.Message.Author.Username, reason),
		id,
	))
	return err
}
