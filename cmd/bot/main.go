package main

import (
	"context"
	"encoding/json"
	"moderari/internal/commands"
	"moderari/internal/config"
	"moderari/internal/db"
	"moderari/internal/embeds"
	"moderari/internal/events"
	"moderari/internal/http"
	"reflect"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func main() {
	chk(sentry.Init(sentry.ClientOptions{
		Dsn: config.C.SentryDSN,
	}))
	defer sentry.Flush(2 * time.Second)

	router := gommand.NewRouter(&gommand.RouterConfig{
		PrefixCheck: gommand.MultiplePrefixCheckers(
			gommand.MentionPrefix,
			func(ctx *gommand.Context, r *gommand.StringIterator) bool {
				guildString, err := db.Client.Get("guild:" + ctx.Message.GuildID.String()).Result()
				guild := db.GuildModel{Prefix: config.C.Prefix}
				if err == nil {
					_ = json.Unmarshal([]byte(guildString), &guild)
				}

				bytes := []byte(guild.Prefix)
				l := len(bytes)

				if res := prefixIterator(l, bytes, r); !res {
					return false
				}
				ctx.Prefix = guild.Prefix
				return true
			},
			func(ctx *gommand.Context, r *gommand.StringIterator) bool {
				userString, err := db.Client.Get("user:" + ctx.Message.Author.ID.String()).Result()
				user := db.UserModel{Prefix: config.C.Prefix}
				if err == nil {
					_ = json.Unmarshal([]byte(userString), &user)
				}

				bytes := []byte(user.Prefix)
				l := len(bytes)

				if res := prefixIterator(l, bytes, r); !res {
					return false
				}
				ctx.Prefix = user.Prefix
				return true
			},
		),
	})
	client := disgord.New(disgord.Config{
		BotToken: config.C.Token,
		// We can actually just use this logger in events; it's a member of Client.
		Logger: logrus.New(),
	})

	router.AddErrorHandler(func(ctx *gommand.Context, err error) bool {
		switch err.(type) {
		case *gommand.CommandNotFound, *gommand.CommandBlank:
			return true
		case *gommand.InvalidTransformation:
			_, _ = ctx.Reply(embeds.Error("Invalid Type", err, false))
			return true
		case *gommand.IncorrectPermissions:
			_, _ = ctx.Reply(embeds.Error("Missing Permissions", err, false))
			return true
		case *gommand.InvalidArgCount:
			_, _ = ctx.Reply(embeds.Error("Missing Arguments", err, false))
			return true
		case *gommand.PanicError:
			eventID := sentry.CaptureException(err)
			_, _ = ctx.Reply(embeds.Error("Panic ("+string(*eventID)+")", err, true))
			// We also want to log to the console.
			return false
		default:
			eventID := sentry.CaptureException(err)
			_, _ = ctx.Reply(embeds.Error("Handled Error: "+reflect.TypeOf(err).String()+" ("+string(*eventID)+")", err, true))
			return false
		}
	})
	events.Register(client)
	commands.Register(router)

	router.Hook(client)
	go func() {
		chk(http.Serve())
	}()
	chk(client.StayConnectedUntilInterrupted(context.Background()))
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func prefixIterator(l int, bytes []byte, r *gommand.StringIterator) bool {
	i := 0
	for i != l {
		b, err := r.GetChar()
		if err != nil {
			return false
		}
		if b != bytes[i] {
			return false
		}
		i++
	}
	return true
}