package discord

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/spf13/cobra"
)

var Clean = &cobra.Command{
	Use:   "discord:clean",
	Short: "Cleans the registered slash commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := "Bot " + config.GetBotToken()
		dg, err := discordgo.New(token)
		if err != nil {
			return fmt.Errorf("failed to create the discord session: %v", err)
		}

		appId := config.GetAppID()
		guildId := config.GetGuildID()
		emptyCmd := make([]*discordgo.ApplicationCommand, 0)
		_, err = dg.ApplicationCommandBulkOverwrite(appId, guildId, emptyCmd)
		if err != nil {
			return fmt.Errorf("failed to bulk overwrite app commands: %v", err)
		}

		slog.Info("removed application commands")

		return nil
	},
}
