package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewWithSubsystem(s ...string) *slog.Logger {
	return slog.With(slog.String("subsystem", strings.Join(s, "/")))
}

func Register(debug bool) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // TODO: set log level from config
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	// The Default level for this example is info
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug mode activated")
		return
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
