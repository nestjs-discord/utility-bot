package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	commands "github.com/erosdesire/discord-nestjs-utility-bot/cmd/discord/command"
	"github.com/erosdesire/discord-nestjs-utility-bot/cmd/discord/handler"
	"github.com/erosdesire/discord-nestjs-utility-bot/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var removeSlashCommands bool

func init() {
	RunCmd.PersistentFlags().BoolVarP(&removeSlashCommands, "remove-slash-commands", "", true, "remove all slash commands before shutting down")
}

var RunCmd = &cobra.Command{
	Use:   "discord:run",
	Short: "Runs the Discord bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := "Bot " + os.Getenv("DISCORD_BOT_TOKEN")
		guildId := config.GetConfig().GuildID

		dg, err := discordgo.New(token)
		if err != nil {
			return fmt.Errorf("failed to create Discord session: %v", err)
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

		if removeSlashCommands {
			// Remove registered slash command before exiting
			for _, c := range commands.RegisteredCommands {
				err := dg.ApplicationCommandDelete(dg.State.User.ID, guildId, c.ID)
				if err != nil {
					log.Error().Err(err).Str("name", c.Name).Msg("failed to remove slash command")
					continue
				}
				log.Debug().Str("name", c.Name).Msg("removed slash command")
			}
		}

		// Cleanly close down the Discord session
		return dg.Close()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		log.Warn().Msg("discord session closed")
	},
}
