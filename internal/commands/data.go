package commands

import (
	"context"
	"moderari/internal/db"
	"moderari/internal/embeds"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "data",
		Description: "Direct messages you with everything in our database linked to your user.",
		Category:    utilCategory,
		Function:    data,
	})
}

func data(ctx *gommand.Context) error {
	data, err := db.Client.Get("user:" + ctx.Message.Author.ID.String()).Result()
	if err != nil {
		return err
	}

	_, _, err = ctx.Message.Author.SendMsgString(context.Background(), ctx.Session, "```json\n"+data+"\n```")
	if err != nil {
		return err
	}
	_, err = ctx.Reply(embeds.Info("Check DMs", "I've sent you the data we store about you.", ""))
	return err
}
