package discord

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/spf13/cobra"
)

var Clean = &cobra.Command{
	Use:   "discord:clean",
	Short: "Cleans the registered slash commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewHandler()

		bot, err := bot.NewBot(
			bot.WithToken(config.GetBotToken()),
			bot.WithAppId(config.GetAppID()),
			bot.WithGuildId(config.GetGuildID()),
			bot.WithHandler(h),
		)
		if err != nil {
			return err
		}

		return bot.CleanApplicationCommands()
	},
}
