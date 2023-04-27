package content

import (
	"github.com/nestjs-discord/utility-bot/core/cache"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Validate = &cobra.Command{
	Use:   "content:validate",
	Short: "Validates the Markdown content in the configuration to be the correct length",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cache.Content()
		if err != nil {
			return err
		}

		log.Info().
			Int("content-validated", len(config.GetConfig().Commands)).
			Msg("Good job! everything looks fine :)")

		return nil
	},
}
