package logger

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

const subsystemKey = "subsystem"

func NewWithSubsystem(s ...string) *slog.Logger {
	return slog.With(
		slog.String(subsystemKey, strings.Join(s, "/")),
	)
}

func Initialize(stageCfg env.Stage) error {
	switch stageCfg {
	case env.StageProd:
		slog.SetDefault(newProductionLogger())
		return nil
	case env.StageDev:
		slog.SetDefault(newDevelopmentLogger())
		return nil
	}

	return fmt.Errorf("invalid stage: '%s'", stageCfg)
}

func newProductionLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

func newDevelopmentLogger() *slog.Logger {
	opts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	}
	handler := tint.NewHandler(os.Stdout, opts)
	return slog.New(handler)
}
