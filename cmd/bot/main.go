package main

import (
	"context"
	"moderari/internal/commands"
	"moderari/internal/config"
	"moderari/internal/embeds"
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
		PrefixCheck: gommand.MultiplePrefixCheckers(gommand.MentionPrefix, gommand.StaticPrefix(config.C.Prefix)),
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
	commands.Register(router)

	router.Hook(client)
	chk(client.StayConnectedUntilInterrupted(context.Background()))
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
