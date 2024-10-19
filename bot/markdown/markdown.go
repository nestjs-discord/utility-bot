package markdown

import (
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Options struct {
	Commands yaml.Commands
}

type Markdown struct {
	opts   Options
	logger *slog.Logger
}

func NewMarkdown(opts Options) *Markdown {
	return &Markdown{
		opts:   opts,
		logger: logger.NewWithSubsystem("bot", "markdown"),
	}
}
