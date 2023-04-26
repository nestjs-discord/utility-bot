package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	internalDiscord "github.com/nestjs-discord/utility-bot/internal/discord"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var Run = &cobra.Command{
	Use:   "discord:run",
	Short: "Runs the Discord bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		dg, err := internalDiscord.NewSession()
		if err != nil {
			return err
		}

		// Discord event handlers
		dg.AddHandler(handler.Ready)
		//dg.AddHandler(handlers.MessageCreate)
		dg.AddHandler(handler.InteractionCreate)

		// We only care about receiving message events
		dg.Identify.Intents = discordgo.IntentsGuildMessages

		// Open a websocket connection to Discord and begin listening
		err = dg.Open()
		if err != nil {
			return fmt.Errorf("failed to open Discord connection: %v", err)
		}

		// Graceful shutdown
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
		signal := <-sc

		log.Warn().Str("signal", signal.String()).Msg("shutting down")

		// Cleanly close down the Discord session
		return dg.Close()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		log.Warn().Msg("discord session closed")
	},
}
