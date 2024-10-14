package main

import (
	"github.com/joho/godotenv"
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/cache"
	"github.com/nestjs-discord/utility-bot/internal/discord/command"
	"github.com/nestjs-discord/utility-bot/internal/discord/forms"
	"github.com/nestjs-discord/utility-bot/internal/logger"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = godotenv.Load()
	logger.Register()

	botCfg, err := config.NewBotConfig()
	if err != nil {
		log.Fatal(err)
	}

	yamlCfg, err := config.NewYamlConfig(config.YamlFile)
	if err != nil {
		log.Fatal(err)
	}

	err = cache.Content(yamlCfg.Commands) // Cache Markdown content // TODO: avoid global instance
	if err != nil {
		log.Fatal(err)
	}

	cache.InitRateLimit(config.Yaml().RateLimit.TTL)

	cache.InitAutoMod()

	b, err := bot.NewBot(botCfg, nil) // TODO: init handler
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
	channels, err := session.GuildChannels(botCfg.GuildId)
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

	slog.Info("shutting down",
		slog.String("signal", signalReceived.String()),
	)

	// Cleanly close down the Discord session
	err = session.Close()
	if err != nil {
		log.Fatalf("failed to close Discord connection: %v", err)
	}
}
