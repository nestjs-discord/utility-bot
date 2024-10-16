package markdown

import (
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Markdown struct {
	logger   *slog.Logger
	commands yaml.Commands
}

func NewMarkdown(commands yaml.Commands) *Markdown {
	return &Markdown{
		logger:   logger.NewWithSubsystem("bot", "markdown"),
		commands: commands,
	}
}
