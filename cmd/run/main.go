package main

import (
	"flag"
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
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

var (
	yamlConfigPath = flag.String("yaml-config-path", "./config.yml", "")
)

func init() {
	flag.Parse()

}

func initDependencies() *bot.Bot {
	stageCfg, err := env.NewStageConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(stageCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Environment variables
	discordCfg, err := env.NewDiscordConfig()
	if err != nil {
		log.Fatal(err)
	}

	// YAML configuration
	yamlCfg, err := yaml.NewConfig(yaml.Path(*yamlConfigPath))
	if err != nil {
		log.Fatal(err)
	}
	yamlCommands := yaml.NewCommands(yamlCfg)
	yamlModerators := yaml.NewModerators(yamlCfg)
	yamlRateLimit := yaml.NewRateLimit(yamlCfg)
	yamlForms := yaml.NewForms(yamlCfg)
	yamlAntiSpam := yaml.NewAntispam(yamlCfg)

	// Initialize the Discord bot
	b, err := bot.NewBot(discordCfg)
	if err != nil {
		log.Fatal(err)
	}

	iMarkdown := markdown.NewMarkdown(yamlCommands)

	iModerators, err := moderators.NewModerators(yamlModerators)
	if err != nil {
		log.Fatal(err)
	}

	iRateLimit := rate_limit.NewRateLimit(yamlRateLimit, iModerators)

	iAntispam, err := antispam.NewAntispam(yamlAntiSpam, iModerators)
	if err != nil {
		log.Fatal(err)
	}

	session := bot.ProvideSession(b)
	iForms, err := forms.NewForms(yamlForms, session)
	if err != nil {
		log.Fatal(err)
	}

	iArchive := archive.New(iModerators)
	iReference := reference.New()
	iSolved := solved.New(iModerators)
	iDontPingMods := dont_ping_mods.NewDontPingMods(iModerators)
	iGoogleIt := google_it.NewGoogleIt()

	interactionHandler := interaction.NewHandler(
		iForms,
		iModerators,
		iRateLimit,
		iMarkdown,
		iArchive,
		iGoogleIt,
		iReference,
		iSolved,
		iDontPingMods,
	)

	b.ApplyHandler(
		handler.NewHandler(
			interactionHandler,
			iAntispam,
			iForms,
			iMarkdown,
			iModerators,
		),
	)

	_, err = commands.NewCommands(
		session,
		discordCfg,
		yamlCommands,
		iArchive,
		iGoogleIt,
		iReference,
		iSolved,
	)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func main() {
	b := initDependencies()
	err := b.OpenWebsocketConnection()
	if err != nil {
		log.Fatalf("bot open failed: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	slog.Info("shutting down")
	err = b.Close()
	if err != nil {
		log.Fatalf("bot close failed: %v", err)
	}
	slog.Info("shutdown done")
}
