package cmd

import (
	"github.com/nestjs-discord/utility-bot/cmd/content"
	"github.com/nestjs-discord/utility-bot/cmd/discord"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/nestjs-discord/utility-bot/core/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "sets log level to debug")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")

	rootCmd.AddCommand(content.Validate)

	rootCmd.AddCommand(discord.Clean)
	rootCmd.AddCommand(discord.Run)
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	logger.Register(debug)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("viper: config read in failed")
	}

	log.Debug().Str("path", viper.ConfigFileUsed()).Msg("config: read success")

	if err := config.Unmarshal(); err != nil {
		log.Fatal().Err(err).Msg("config: unmarshal failed")
	}

	if err := config.ValidateConfig(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
