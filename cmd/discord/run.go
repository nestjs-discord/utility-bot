package discord

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/cache"
	internalDiscord "github.com/nestjs-discord/utility-bot/internal/discord"
	"github.com/nestjs-discord/utility-bot/internal/discord/command"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var Run = &cobra.Command{
	Use:   "discord:run",
	Short: "Starts the Discord bot",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := cache.Content() // Cache Markdown content
		if err != nil {
			return err
		}

		cache.InitRatelimit(config.GetYaml().Ratelimit.TTL)

		cache.InitAutoMod()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		session, err := internalDiscord.NewSession()
		if err != nil {
			return fmt.Errorf("failed to create new Discord session: %s", err)
		}

		log.Info().Str("link", internalDiscord.GenerateInviteLink()).Msg("server invite")

		command.RegisterApplicationCommands(session)

		// Discord event handlers
		session.AddHandler(handler.Ready)

		session.AddHandler(handler.MessageCreate)
		session.AddHandler(handler.InteractionCreate)

		// We only care about receiving message events
		session.Identify.Intents = config.BotIntents

		// Fetch all the channels
		channels, err := session.GuildChannels(config.GetGuildID())
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
