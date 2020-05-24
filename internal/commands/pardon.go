package commands

import (
	"encoding/json"
	"moderari/internal/db"
	"moderari/internal/embeds"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:                 "pardon",
		Description:          "Removes warns. If warn IDs aren't provided, all of the member's warns are cleared.",
		Aliases:              []string{"rmwarn", "rmwarns", "delwarn", "delwarns", "rmstrike", "rmstrikes", "delstrike", "delstrikes"},
		Usage:                "<member> [warn IDs...]",
		Category:             modCategory,
		PermissionValidators: []gommand.PermissionValidator{gommand.MANAGE_MESSAGES},
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
		Function: pardon,
	})
}

func pardon(ctx *gommand.Context) error {
	key := "user:" + ctx.Args[0].(*disgord.Member).User.ID.String()
	oldString, err := db.Client.Get(key).Result()
	if err != nil {
		return err
	}

	old := db.UserModel{}
	if err := json.Unmarshal([]byte(oldString), &old); err != nil {
		return err
	}

	guildID := ctx.Message.GuildID.String()
	removedWarns := 0
	if ctx.Args[1] != nil {
		for _, v := range strings.Split(ctx.Args[1].(string), " ") {
			delete(old.Warns[guildID], v)
			removedWarns++
		}
	} else {
		_, err = ctx.Reply(embeds.Info("Are you sure?", "This will clear all warnings.", ""))
		if err != nil {
			return err
		}

		res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
			return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
		})
		if res.Content[0] != 'y' {
			ctx.Reply("Cancelled.")
			return nil
		}

		removedWarns = len(old.Warns[guildID])
		delete(old.Warns, guildID)
	}

	go func() {
		bytes, _ := json.Marshal(old)
		db.Client.Set(key, bytes, 0)
	}()

	_, err = ctx.Reply(embeds.Info("Pardoned", strconv.Itoa(removedWarns)+" warns were removed.", ""))
	return err
}
