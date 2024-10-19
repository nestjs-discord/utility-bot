package app

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/commands"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/bot/session"
	"github.com/nestjs-discord/utility-bot/modules/cron"
)

type App struct {
	Bot *bot.Bot
}

func NewApp(
	b *bot.Bot,
	s *session.Session,
	h *handler.Handler,
	c *cron.Cron,
	_ *commands.Commands,
) *App {
	s.ApplyHandler(h)
	c.Start()
	return &App{
		Bot: b,
	}
}
