//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/nestjs-discord/utility-bot/app"
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/antispam"
	"github.com/nestjs-discord/utility-bot/bot/auto_mod"
	"github.com/nestjs-discord/utility-bot/bot/commands"
	"github.com/nestjs-discord/utility-bot/bot/commands/archive"
	"github.com/nestjs-discord/utility-bot/bot/commands/credits"
	"github.com/nestjs-discord/utility-bot/bot/commands/dont_ping_mods"
	"github.com/nestjs-discord/utility-bot/bot/commands/google_it"
	"github.com/nestjs-discord/utility-bot/bot/commands/reference"
	"github.com/nestjs-discord/utility-bot/bot/commands/solved"
	"github.com/nestjs-discord/utility-bot/bot/cron"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/bot/handler/interaction"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/bot/rate_limit"
	"github.com/nestjs-discord/utility-bot/bot/session"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
)

func InitializeApp() (*app.App, func(), error) {
	panic(wire.Build(
		// infra/config/env
		wire.NewSet(
			env.NewStageConfig,
			env.NewDiscordConfig,
			env.ProvideGuildId,
		),

		logger.Initialize,

		// infra/config/yaml
		wire.NewSet(
			wire.NewSet(
				wire.Value(yaml.Path("config.yml")),
				yaml.NewConfig,
			),
			wire.NewSet(
				yaml.NewModerators,
				yaml.NewRateLimit,
				yaml.NewAntispam,
				yaml.NewForms,
				yaml.NewArchiveCommand,
				yaml.NewSolvedCommand,
				yaml.NewCommands,
			),
		),

		// bot features
		antispam.NewAntispam,
		auto_mod.NewAutoMod,
		markdown.NewMarkdown,
		moderators.NewModerators,
		rate_limit.NewRateLimit,
		forms.NewForms,
		status.NewStatus,

		// commands
		archive.New,
		credits.NewCredits,
		reference.New,
		solved.New,
		dont_ping_mods.NewDontPingMods,
		google_it.NewGoogleIt,

		interaction.NewHandler,
		handler.NewHandler,

		wire.NewSet(
			bot.NewBot,
			bot.ProvideSession,
			session.ProvideSession,
		),

		commands.NewCommands,

		wire.NewSet(
			wire.Struct(new(cron.Option), "*"),
			cron.NewCron,
		),

		app.NewApp,
	))
}
