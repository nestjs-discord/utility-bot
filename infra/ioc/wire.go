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
	"github.com/nestjs-discord/utility-bot/modules/cron"
)

var botSet = wire.NewSet(
	wire.NewSet(
		wire.Struct(new(antispam.Options), "*"),
		antispam.NewAntispam,
	),
	wire.NewSet(
		wire.Struct(new(auto_mod.Options), "*"),
		auto_mod.NewAutoMod,
	),
	wire.NewSet(
		botCommands,
		wire.Struct(new(commands.Options), "*"),
		commands.NewCommands,
	),
	wire.NewSet(
		wire.Struct(new(forms.Options), "*"),
		forms.NewForms,
	),
	wire.NewSet(
		wire.NewSet(
			wire.Struct(new(interaction.Options), "*"),
			interaction.NewInteractionHandler,
		),
		handler.NewHandler,
	),
	markdown.NewMarkdown,
	moderators.NewModerators,
	rate_limit.NewRateLimit,
	status.NewStatus,
)

var botCommands = wire.NewSet(
	archive.NewArchive,
	credits.NewCredits,
	dont_ping_mods.NewDontPingMods,
	google_it.NewGoogleIt,
	reference.NewReference,
	solved.NewSolved,
)

var infra = wire.NewSet(
	// infra/config/env
	wire.NewSet(
		env.NewStageConfig,
		env.NewDiscordConfig,
		env.ProvideGuildId,
	),
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
	// infra/logger
	logger.NewLogger,
)

var modules = wire.NewSet(
	wire.NewSet(
		wire.Struct(new(cron.Option), "*"),
		cron.NewCron,
	),
)

func InitializeApp() (*app.App, func(), error) {
	panic(wire.Build(
		app.NewApp,
		wire.NewSet(
			bot.NewBot,
			bot.ProvideSession,
			session.ProvideSession,
			botSet,
		),
		infra,
		modules,
	))
}
