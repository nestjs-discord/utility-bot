package interaction

import (
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"

	"github.com/nestjs-discord/utility-bot/bot/moderator"
)

type Handler struct {
	logger    *slog.Logger
	moderator *moderator.Moderator
}

func NewHandler(
	moderator *moderator.Moderator,
) *Handler {
	return &Handler{
		logger:    logger.NewWithSubsystem("bot", "handler", "interaction"),
		moderator: moderator,
	}
}
