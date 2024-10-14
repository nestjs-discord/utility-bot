package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

const subsystemKey = "subsystem"

var Level = slog.LevelDebug // TODO: set log level from config

func NewWithSubsystem(s ...string) *slog.Logger {
	return slog.With(
		slog.String(subsystemKey, strings.Join(s, "/")),
	)
}

func Initialize(stage string) error {
	switch stage {
	case "prod":
		slog.SetDefault(newProductionLogger())
		return nil
	case "dev":
		slog.SetDefault(newDevelopmentLogger())
		return nil
	}

	return fmt.Errorf("invalid stage: '%s'", stage)
}

func newProductionLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: Level,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

func newDevelopmentLogger() *slog.Logger {
	opts := &tint.Options{
		Level:      Level,
		TimeFormat: time.Kitchen,
	}
	handler := tint.NewHandler(os.Stdout, opts)
	return slog.New(handler)
}
