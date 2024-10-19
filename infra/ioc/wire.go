//go:build wireinject
// +build wireinject

package ioc

import (
	w "github.com/google/wire"
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

func InitializeApp() (*app.App, func(), error) {
	panic(w.Build(
		app.NewApp,

		// bot
		w.NewSet(
			bot.NewBot,
			bot.ProvideSession,
			session.ProvideSession,
			w.NewSet(
				w.NewSet(antispam.NewAntispam, w.Struct(new(antispam.Options), "*")),
				w.NewSet(auto_mod.NewAutoMod, w.Struct(new(auto_mod.Options), "*")),
				w.NewSet(
					commands.NewCommands, w.Struct(new(commands.Options), "*"),
					w.NewSet(
						archive.NewArchive,
						credits.NewCredits,
						dont_ping_mods.NewDontPingMods,
						google_it.NewGoogleIt,
						reference.NewReference,
						solved.NewSolved,
					),
				),
				w.NewSet(forms.NewForms, w.Struct(new(forms.Options), "*")),
				w.NewSet(
					w.NewSet(interaction.NewInteractionHandler, w.Struct(new(interaction.Options), "*")),
					w.NewSet(handler.NewHandler, w.Struct(new(handler.Options), "*")),
				),
				w.NewSet(markdown.NewMarkdown, w.Struct(new(markdown.Options), "*")),
				w.NewSet(moderators.NewModerators, w.Struct(new(moderators.Options), "*")),
				w.NewSet(rate_limit.NewRateLimit, w.Struct(new(rate_limit.Options), "*")),
				w.NewSet(status.NewStatus, w.Struct(new(status.Options), "*")),
			),
		),

		// infra
		w.NewSet(

			// infra/config/env
			w.NewSet(
				env.NewStageConfig,
				env.NewDiscordConfig,
				env.ProvideGuildId,
			),

			// infra/config/yaml
			w.NewSet(
				w.NewSet(
					w.Value(yaml.Path("config.yml")),
					yaml.NewConfig,
				),
				w.NewSet(
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
		),

		// modules
		w.NewSet(
			w.NewSet(cron.NewCron, w.Struct(new(cron.Option), "*")),
		),
	))
}
