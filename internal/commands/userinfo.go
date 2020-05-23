package commands

import (
	"encoding/json"
	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
	"moderari/internal/config"
	"moderari/internal/db"
	"moderari/internal/embeds"
	"strconv"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "userinfo",
		Usage: "[member]",
		Aliases:     []string{"ui", "profile", "user", "uinfo"},
		Description: "Shows a user's profile.",
		Category:    infoCategory,
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: true,
				Function: gommand.MemberTransformer,
			},
		},
		Function: userInfo,
	})
}

func userInfo(ctx *gommand.Context) error {
	member := ctx.Message.Member
	if ctx.Args[0] != nil {
		member = ctx.Args[0].(*disgord.Member)
	}

	// TODO: This is really messy at the moment.
	userID := member.User.ID.String()
	userDocString, err := db.Client.Get("user:" + userID).Result()
	if err != nil {
		return err
	}
	userDoc := db.UserModel{}
	_ = json.Unmarshal([]byte(userDocString), &userDoc)

	avatarURL, _ := member.User.AvatarURL(2048, true)
	nickname := member.Nick
	if nickname == "" {
		nickname = member.User.Username
	}

	_, err = ctx.Reply(&disgord.Embed{
		Title: member.User.Tag(),
		Color: config.C.Colors.Info,
		Thumbnail: &disgord.EmbedThumbnail{
			URL: avatarURL,
		},
		Fields: []*disgord.EmbedField{
			embeds.Field("Mention", member.User.Mention(), true),
			embeds.Field("ID", userID, true),
			embeds.Field("Nickname", nickname, true),
			embeds.Field("Warnings", strconv.Itoa(len(userDoc.Warns[ctx.Message.GuildID.String()])), true),
			embeds.Field("Joined", member.JoinedAt.Format("15:04:05 on January 2, 2006"), true),
		},
	})
	return err
}
