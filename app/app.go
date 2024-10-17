package app

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/commands"
	"github.com/nestjs-discord/utility-bot/bot/handler"
)

type App struct {
	Bot *bot.Bot
}

func NewApp(
	b *bot.Bot,
	h *handler.Handler,
	_ *commands.Commands,
) *App {
	b.ApplyHandler(h)

	return &App{
		Bot: b,
	}
}
