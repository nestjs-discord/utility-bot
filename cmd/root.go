package cmd

import (
	"github.com/nestjs-discord/utility-bot/cmd/content"
	"github.com/nestjs-discord/utility-bot/cmd/discord"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "utility-bot",
		Short: "NestJS Discord Utility Bot",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			HiddenDefaultCmd:    true,
			DisableDescriptions: true,
			DisableNoDescFlag:   true,
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	rootCmd.AddCommand(content.Validate)

	rootCmd.AddCommand(discord.Clean)
	rootCmd.AddCommand(discord.Run)
}

func Execute() error {
	return rootCmd.Execute()
}
