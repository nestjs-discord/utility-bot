package main

import (
	"flag"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/handler/interaction"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/moderator"
	"github.com/nestjs-discord/utility-bot/bot/rate_limit"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/automod"
	"github.com/nestjs-discord/utility-bot/bot/handler"
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

func main() {
	// Environment variables
	discordCfg, err := env.NewDiscordConfig()
	if err != nil {
		log.Fatal(err)
	}

	// YAML configuration
	yamlCfg, err := yaml.NewConfig(*yamlConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the Discord bot
	b, err := bot.NewBot(discordCfg)
	if err != nil {
		log.Fatal(err)
	}

	iMarkdown := markdown.NewMarkdown()
	err = iMarkdown.CacheCommands(yamlCfg.Commands)
	if err != nil {
		log.Fatal(err)
	}

	iModerator := moderator.NewModerator(yamlCfg.Moderators)
	iRateLimit := rate_limit.NewRateLimit(yamlCfg.RateLimit, iModerator)

	iAutoMod, err := automod.NewAutoMod(yamlCfg.AutoMod, iModerator)
	if err != nil {
		log.Fatal(err)
	}

	iForms, err := forms.NewForms(yamlCfg.Forms, b.Session()) // TODO: find a way not to pass the session
	if err != nil {
		log.Fatal(err)
	}

	interactionHandler := interaction.NewHandler(iModerator, iRateLimit, iMarkdown)

	b.ApplyHandler(
		handler.NewHandler(
			interactionHandler,
			iAutoMod,
			iForms,
			iMarkdown,
			iModerator,
		),
	)

	err = b.RegisterApplicationCommands(yamlCfg.Commands)
	if err != nil {
		log.Fatal(err)
	}

	// Open a websocket connection to Discord and begin listening
	err = b.Open()
	if err != nil {
		log.Fatalf("failed to open Discord connection: %v", err)
	}

	// Graceful shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	signalReceived := <-sc

	slog.Info("signal received",
		slog.String("signal", signalReceived.String()),
	)

	// Cleanly close down the Discord session
	err = b.Close()
	if err != nil {
		log.Fatalf("unable to close the session: %v", err)
	}

	slog.Info("shutdown completed")
}
