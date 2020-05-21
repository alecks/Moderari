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
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "warns",
		Aliases:     []string{"strikes", "warnings"},
		Description: "Retrieves a member's warnings.",
		Usage:       "[member] [minimum severity]",
		Category:    modCategory,
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: true,
				Function: gommand.MemberTransformer,
			},
			{
				Optional: true,
				Function: gommand.UIntTransformer,
			},
		},
		Function: warns,
	})
}

func warns(ctx *gommand.Context) error {
	// TODO: Improve this.
	userID := ctx.Message.Author.ID.String()
	if ctx.Args[0] != nil {
		userID = ctx.Args[0].(*disgord.Member).User.ID.String()
	}
	var minSeverity uint64 = 0
	if ctx.Args[1] != nil {
		minSeverity = ctx.Args[1].(uint64)
	}

	userString, err := db.Client.Get("user:" + userID).Result()
	if err != nil {
		return err
	}

	user := db.UserModel{}
	if err := json.Unmarshal([]byte(userString), &user); err != nil {
		return err
	}

	paginatorEmbeds := []*disgord.Embed{}
	for id, v := range user.Warns[ctx.Message.GuildID.String()] {
		if v.Severity < minSeverity {
			continue
		}

		fields := []*disgord.EmbedField{
			embeds.Field("Severity", fmt.Sprintf("%d", v.Severity), false),
			embeds.Field("Time", time.Unix(v.Time, 0).String(), false),
		}

		moderator, err := ctx.Session.GetUser(context.Background(), disgord.ParseSnowflakeString(v.Moderator))
		if err == nil {
			fields = append(fields, embeds.Field("Moderator", moderator.Username, false))
		}

		paginatorEmbeds = append(paginatorEmbeds, embeds.Info(
			v.Reason,
			id,
			"",
			fields...,
		))
	}

	if len(paginatorEmbeds) <= 0 {
		_, err := ctx.Reply(embeds.Info("Clean Slate", "They don't have any punishments.", ""))
		return err
	}

	return gommand.EmbedsPaginator(ctx, paginatorEmbeds)
}
