package commands

import (
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

	key := "user:" + member.User.ID.String()
	oldString, err := db.Client.Get(key).Result()
	if err == db.Nil {
		oldString = "{}"
	} else if err != nil {
		return err
	}

	old := db.UserModel{}
	if err := json.Unmarshal([]byte(oldString), &old); err != nil {
		return err
	}

	id := uuid.New().String()
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

	new, err := json.Marshal(old)
	if err != nil {
		return err
	}
	go db.Client.Set(key, new, 0)

	_, err = ctx.Reply(embeds.Info(
		"Warned",
		fmt.Sprintf("**%s** has been warned by **%s** for \"%s\".", member.User.Username, ctx.Message.Author.Username, warning.Reason),
		id,
	))
	return err
}
