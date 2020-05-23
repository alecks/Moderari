package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/getsentry/sentry-go"
	"math/rand"
	"moderari/internal/config"
	"moderari/internal/db"
	"moderari/internal/embeds"
	"moderari/internal/http"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateIdentifier() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func init() {
	events[disgord.EvtGuildMemberAdd] = func(session disgord.Session, evt *disgord.GuildMemberAdd) {
		guildString, err := db.Client.Get("guild:" + evt.Member.GuildID.String()).Result()
		chk(err, session)

		guildDoc := db.GuildModel{}
		_ = json.Unmarshal([]byte(guildString), &guildDoc)
		if !guildDoc.Gotcha {
			return
		}
		guild, err := session.GetGuild(context.Background(), evt.Member.GuildID)
		chk(err, session)

		identifier := generateIdentifier()
		message := disgord.NewMessage()
		message.Embeds = append(message.Embeds, embeds.Info(
			"Hey! "+guild.Name+" is protected with Gotcha.",
			fmt.Sprintf(
				"To gain access, see %s/verify/%s.\n[Here](https://github.com/fjah/gotcha)'s the source.",
				config.C.URL,
				identifier,
			),
			"You have one hour to do this.",
		))
		if _, _, err := evt.Member.User.SendMsg(context.Background(), session, message); err != nil {
			// The user most likely has DMs off or blocked the bot. No need to log the error.
			return
		}

		go func() {
			message = disgord.NewMessage()
			message.Embeds = append(message.Embeds, embeds.Info("You're in!", "We hope you enjoy your time.", ""))
			stat := http.Gotcha.Await(identifier)
			switch stat {
			case 1:
				message.Embeds[0] = embeds.ErrorString(
					"Uhh...",
					"You didn't respond in one hour. Please contact server staff.",
				)
			case 2:
				message.Embeds[0] = embeds.ErrorString(
					"Uhh...",
					"You seem to be on our blocklist. Have you joined the same server already today?",
				)
			}

			if _, _, err := evt.Member.User.SendMsg(context.Background(), session, message); err != nil {
				return
			}
		}()
	}
}

func chk(err error, session disgord.Session) {
	if err != nil {
		if err == db.Nil {
			return
		}
		sentry.CaptureException(err)
		session.Logger().Error(err)
		return
	}
}
