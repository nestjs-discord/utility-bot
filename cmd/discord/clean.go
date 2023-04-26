package discord

import (
	"github.com/nestjs-discord/utility-bot/core/config"
	internalDiscord "github.com/nestjs-discord/utility-bot/internal/discord"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Clean = &cobra.Command{
	Use:   "discord:clean",
	Short: "Clean registered slash commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		dg, err := internalDiscord.NewSession()
		if err != nil {
			return err
		}

		appId := config.GetAppID()
		guildId := config.GetGuildID()

		// Fetch registered slash commands
		registeredCommands, err := dg.ApplicationCommands(appId, guildId)
		if err != nil {
			return err
		}

		log.Warn().Int("len", len(registeredCommands)).Msg("fetched registered slash commands")

		// Remove registered slash command
		for _, c := range registeredCommands {
			if err := dg.ApplicationCommandDelete(appId, guildId, c.ID); err != nil {
				log.Error().Err(err).Str("name", c.Name).Msg("failed to remove slash command")
				continue
			}

			log.Info().Str("name", c.Name).Msg("removed slash command")
		}

		log.Info().Msg("done")

		return nil
	},
}
