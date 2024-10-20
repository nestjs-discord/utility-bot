package app

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/commands"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/bot/session"
	"github.com/nestjs-discord/utility-bot/modules/cron"
)

type Options struct {
	Bot      *bot.Bot
	Session  *session.Session
	Handler  *handler.Handler
	Cron     *cron.Cron
	Commands *commands.Commands
}

type App struct {
	Opts Options
}

func NewApp(opts Options) *App {
	return &App{
		Opts: opts,
	}
}
