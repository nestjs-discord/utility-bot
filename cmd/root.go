package cmd

import (
	"github.com/erosdesire/discord-nestjs-utility-bot/cmd/discord"
	"github.com/erosdesire/discord-nestjs-utility-bot/config"
	util "github.com/erosdesire/discord-nestjs-utility-bot/logger"
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

	rootCmd.AddCommand(discord.RunCmd)
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
		log.Debug().Str("path", p).Msg("loaded file content")
	}
}
