//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/antispam"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/bot/handler/interaction"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/bot/rate_limit"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
)

//type App struct {
//	bot *bot.Bot
//}
//
//func NewApp(b *bot.Bot) *App {
//	return &App{
//		bot: b,
//	}
//}

func InitializeDependencies() (*bot.Bot, error) {
	wire.Build(
		env.NewDiscordConfig,
		yaml.Set,
		markdown.NewMarkdown,
		moderators.NewModerators,
		rate_limit.NewRateLimit,
		antispam.NewAntispam,
		forms.NewForms,
		interaction.NewHandler,
		handler.NewHandler,
		bot.Set,
		//NewApp,
	)
	return &bot.Bot{}, nil
}
