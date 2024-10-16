package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	"github.com/nestjs-discord/utility-bot/infra/logger"
)

var (
	stage          = flag.String("stage", "dev", "prod,dev")
	yamlConfigPath = flag.String("yaml-config-path", "./config.yml", "")
)

func init() {
	flag.Parse()

	err := logger.Initialize(*stage)
	if err != nil {
		log.Fatal(err)
	}
}

func initDependencies() *bot.Bot {
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

	// Initialize the Discord bot
	b, err := bot.NewBot(discordCfg)
	if err != nil {
		log.Fatal(err)
	}

	iMarkdown := markdown.NewMarkdown(yamlCfg.Commands)

	iModerators, err := moderators.NewModerators(yamlCfg.Moderators)
	if err != nil {
		log.Fatal(err)
	}

	iRateLimit := rate_limit.NewRateLimit(yamlCfg.RateLimit, iModerators)

	iAntispam, err := antispam.NewAntispam(yamlCfg.Antispam, iModerators)
	if err != nil {
		log.Fatal(err)
	}

	session := bot.ProvideSession(b)
	iForms, err := forms.NewForms(yamlCfg.Forms, session)
	if err != nil {
		log.Fatal(err)
	}

	interactionHandler := interaction.NewHandler(iForms, iModerators, iRateLimit, iMarkdown)

	b.ApplyHandler(
		handler.NewHandler(
			interactionHandler,
			iAntispam,
			iForms,
			iMarkdown,
			iModerators,
		),
	)

	err = b.RegisterApplicationCommands(yamlCfg.Commands)
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

	// Graceful shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM) // TODO: remove unnecessary ones
	signalReceived := <-sc

	slog.Info("signal received",
		slog.String("signal", signalReceived.String()),
	)

	err = b.Close()
	if err != nil {
		log.Fatalf("bot close failed: %v", err)
	}

	slog.Info("shutdown completed")
}
