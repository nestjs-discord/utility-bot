//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/antispam"
	"github.com/nestjs-discord/utility-bot/bot/commands"
	"github.com/nestjs-discord/utility-bot/bot/commands/archive"
	"github.com/nestjs-discord/utility-bot/bot/commands/dont_ping_mods"
	"github.com/nestjs-discord/utility-bot/bot/commands/google_it"
	"github.com/nestjs-discord/utility-bot/bot/commands/reference"
	"github.com/nestjs-discord/utility-bot/bot/commands/solved"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/bot/handler/interaction"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/bot/rate_limit"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
)

type App struct {
	bot *bot.Bot
}

func newApp(
	b *bot.Bot,
	h *handler.Handler,
	_ *commands.Commands,
) *App {
	b.ApplyHandler(h)

	return &App{
		bot: b,
	}
}

func initializeApp() (*App, error) {
	panic(wire.Build(
		// infra/config/env
		wire.NewSet(
			env.NewStageConfig,
			env.NewDiscordConfig,
		),

		// infra/config/yaml
		wire.NewSet(
			wire.NewSet(
				wire.Value(yaml.Path("config.yml")),
				yaml.NewConfig,
			),
			wire.NewSet(
				yaml.NewCommands,
				yaml.NewModerators,
				yaml.NewRateLimit,
				yaml.NewAntispam,
				yaml.NewForms,
			),
		),

		// bot
		wire.NewSet(
			bot.NewBot,
			bot.ProvideSession,
		),

		// bot features
		antispam.NewAntispam,
		markdown.NewMarkdown,
		moderators.NewModerators,
		rate_limit.NewRateLimit,
		forms.NewForms,

		// commands
		archive.New,
		reference.New,
		solved.New,
		dont_ping_mods.NewDontPingMods,
		google_it.NewGoogleIt,

		interaction.NewHandler,
		handler.NewHandler,

		commands.NewCommands,

		newApp,
	))
}
