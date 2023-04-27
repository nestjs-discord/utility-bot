package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Register(debug bool) {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	// Default level for this example is info, unless debug flag is present
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug mode activated")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
