package main

import (
	"flag"
	"github.com/nestjs-discord/utility-bot/config/env"
	"github.com/nestjs-discord/utility-bot/config/yaml"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/cache"
	"github.com/nestjs-discord/utility-bot/internal/discord/command"
	"github.com/nestjs-discord/utility-bot/internal/discord/forms"
	"github.com/nestjs-discord/utility-bot/internal/logger"
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

	// Yaml configuration
	yamlCfg, err := yaml.NewConfig(*yamlConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	err = cache.MarkdownContent(yamlCfg.Commands) // TODO: avoid global instance
	if err != nil {
		log.Fatal(err)
	}

	cache.InitRateLimit(config.Yaml().RateLimit.TTL)

	cache.InitAutoMod()

	b, err := bot.NewBot(discordCfg, nil) // TODO: init handler
	if err != nil {
		log.Fatal(err)
	}

	session := b.Session() // TODO: remove?

	err = forms.Init(yamlCfg.Forms, session)
	if err != nil {
		log.Fatalf("failed to init forms: %s", err)
	}

	command.RegisterApplicationCommands(session)

	// Fetch all the channels
	channels, err := session.GuildChannels(discordCfg.GuildId)
	if err != nil {
		log.Fatalf("failed to fetch guild channels: %s", err)
	}
	cache.AutoMod.SetChannels(channels)

	// Open a websocket connection to Discord and begin listening
	err = session.Open()
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
	err = session.Close()
	if err != nil {
		log.Fatalf("unable to close the session: %v", err)
	}

	slog.Info("shutdown completed")
}
