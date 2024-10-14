package discord

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/cache"
	"github.com/nestjs-discord/utility-bot/internal/discord/command"
	"github.com/nestjs-discord/utility-bot/internal/discord/forms"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var Run = &cobra.Command{
	Use:   "discord:run",
	Short: "Starts the Discord bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		botCfg, err := config.NewBotConfig()
		if err != nil {
			return err
		}

		yamlCfg, err := config.NewYamlConfig(config.YamlFile)
		if err != nil {
			return err
		}

		err = cache.Content(yamlCfg.Commands) // Cache Markdown content // TODO: avoid global instance
		if err != nil {
			return err
		}

		cache.InitRateLimit(config.Yaml().RateLimit.TTL)

		cache.InitAutoMod()

		b, err := bot.NewBot(botCfg, nil) // TODO: init handler
		if err != nil {
			return err
		}

		session := b.Session() // TODO: remove?

		err = forms.Init(yamlCfg.Forms, session)
		if err != nil {
			return fmt.Errorf("failed to init forms: %s", err)
		}

		command.RegisterApplicationCommands(session)

		// Fetch all the channels
		channels, err := session.GuildChannels(botCfg.GuildId)
		if err != nil {
			return fmt.Errorf("failed to fetch guild channels: %s", err)
		}
		cache.AutoMod.SetChannels(channels)

		// Open a websocket connection to Discord and begin listening
		err = session.Open()
		if err != nil {
			return fmt.Errorf("failed to open Discord connection: %v", err)
		}

		// Graceful shutdown
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
		signal := <-sc

		log.Warn().Str("signal", signal.String()).Msg("shutting down")

		// Cleanly close down the Discord session
		return session.Close()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		log.Warn().Msg("discord session closed")
	},
}
