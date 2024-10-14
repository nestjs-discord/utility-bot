package content

import (
	"github.com/nestjs-discord/utility-bot/internal/cache"
	"log/slog"

	"github.com/spf13/cobra"

	"log"
)

var Validate = &cobra.Command{
	Use:   "content:validate",
	Short: "Validates the Markdown content in the configuration to be the correct length",
	Run: func(cmd *cobra.Command, args []string) {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		// TODO: make sure slog is using text handler and the output of this commands looks readable

		err := cache.Content()
		if err != nil {
			log.Fatal(err)
		}

		slog.Info("Good job! everything looks fine :)")
	},
}
