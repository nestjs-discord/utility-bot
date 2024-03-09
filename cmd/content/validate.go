package content

import (
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/cache"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Validate = &cobra.Command{
	Use:   "content:validate",
	Short: "Validates the Markdown content in the configuration to be the correct length",
	Run: func(cmd *cobra.Command, args []string) {
		err := cache.Content()
		if err != nil {
			log.Fatal().Err(err).Send()
			return
		}

		log.Info().
			Int("content-validated", len(config.GetYaml().Commands)).
			Msg("Good job! everything looks fine :)")
	},
}
