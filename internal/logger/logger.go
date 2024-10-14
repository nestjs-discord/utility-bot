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

func Register() {
	w := os.Stdout
	level := slog.LevelDebug // TODO: set log level from config
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// TODO: remove zero log dependency
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stdout,
	})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
