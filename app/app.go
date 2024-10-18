package app

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/commands"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/bot/session"
)

type App struct {
	Bot *bot.Bot
}

func NewApp(
	b *bot.Bot,
	s *session.Session,
	h *handler.Handler,
	_ *commands.Commands,
) *App {
	s.ApplyHandler(h)
	return &App{
		Bot: b,
	}
}
