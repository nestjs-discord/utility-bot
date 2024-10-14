package discord

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/spf13/cobra"
)

var Clean = &cobra.Command{
	Use:   "discord:clean",
	Short: "Cleans the registered slash commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		botCfg, err := config.NewBotConfig()
		if err != nil {
			return err
		}

		// handler := handler.NewHandler()
		bot, err := bot.NewBot(botCfg, nil)
		if err != nil {
			return err
		}

		return bot.CleanApplicationCommands()
	},
}
