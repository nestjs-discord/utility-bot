package discord

import (
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/discord"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Clean = &cobra.Command{
	Use:   "discord:clean",
	Short: "Cleans the registered slash commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		dg, err := discord.NewSession()
		if err != nil {
			return err
		}

		appId := config.GetAppID()
		guildId := config.GetGuildID()

		_, err = dg.ApplicationCommandBulkOverwrite(appId, guildId, []*discordgo.ApplicationCommand{})
		if err != nil {
			return err
		}

		log.Info().Msg("removed application commands")

		return nil
	},
}
