package interaction

import (
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/rate_limit"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"

	"github.com/nestjs-discord/utility-bot/bot/moderators"
)

type Handler struct {
	logger     *slog.Logger
	forms      *forms.Forms
	moderators *moderators.Moderators
	rateLimit  *rate_limit.RateLimit
	markdown   *markdown.Markdown
}

func NewHandler(
	forms *forms.Forms,
	moderators *moderators.Moderators,
	rateLimit *rate_limit.RateLimit,
	markdown *markdown.Markdown,
) *Handler {
	return &Handler{
		logger:     logger.NewWithSubsystem("bot", "handler", "interaction"),
		forms:      forms,
		moderators: moderators,
		rateLimit:  rateLimit,
		markdown:   markdown,
	}
}
