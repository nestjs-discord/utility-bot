package interaction

import (
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/rate_limit"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"

	"github.com/nestjs-discord/utility-bot/bot/moderator"
)

type Handler struct {
	logger    *slog.Logger
	moderator *moderator.Moderator
	rateLimit *rate_limit.RateLimit
	markdown  *markdown.Markdown
}

func NewHandler(
	moderator *moderator.Moderator,
	rateLimit *rate_limit.RateLimit,
	markdown *markdown.Markdown,
) *Handler {
	return &Handler{
		logger:    logger.NewWithSubsystem("bot", "handler", "interaction"),
		moderator: moderator,
		rateLimit: rateLimit,
		markdown:  markdown,
	}
}
