package cmd

import (
	"github.com/nestjs-discord/utility-bot/cmd/content"
	"github.com/nestjs-discord/utility-bot/cmd/discord"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	debug   bool
	cfgFile string

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
	cobra.OnInitialize(func() {
		logger.Register(debug)
		config.Bootstrap(cfgFile)
	})

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "sets log level to debug")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")

	rootCmd.AddCommand(content.Validate)

	rootCmd.AddCommand(discord.Clean)
	rootCmd.AddCommand(discord.Run)
}

func Execute() error {
	return rootCmd.Execute()
}
