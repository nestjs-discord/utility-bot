package cmd

import (
	"github.com/nestjs-discord/utility-bot/cmd/discord"
	"github.com/nestjs-discord/utility-bot/core/config"
	util "github.com/nestjs-discord/utility-bot/core/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	// Used for flags.
	debug   bool
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "nestjs-utility",
		Short: "NestJS Utility",
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

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug log level")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")

	rootCmd.AddCommand(discord.Run)
	rootCmd.AddCommand(discord.Clean)
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	util.RegisterLogger(debug)

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
		log.Fatal().Err(err).Msg("config unmarshal failed")
	}

	// Cache local content on memory
	for _, c := range config.GetConfig().Commands {
		// Ignore non markdown files
		if !strings.HasSuffix(c.Content, ".md") {
			continue
		}

		p := c.Content
		data, err := os.ReadFile(p)
		if err != nil {
			log.Error().Err(err).Str("path", p).Msg("failed to read file content")
		}

		c.Content = string(data)

		// Slash commands can have a maximum of 4000 characters for combined name, description,
		// and value properties for each command, its options (including subcommands and groups), and choices.
		if len(c.Content) > 3500 {
			log.Fatal().Str("path", p).Msg("file content contains too many characters, please consider making it shorter.")
			return
		}

		log.Debug().Str("path", p).Msg("loaded file content")
	}
}
