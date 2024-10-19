package handler

import (
	"github.com/nestjs-discord/utility-bot/bot/antispam"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/handler/interaction"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Options struct {
	InteractionHandler *interaction.Handler
	Antispam           *antispam.Antispam
	Forms              *forms.Forms
	Markdown           *markdown.Markdown
	Moderators         *moderators.Moderators
}

type Handler struct {
	opts   Options
	logger *slog.Logger
}

func NewHandler(opts Options) *Handler {
	h := &Handler{
		opts:   opts,
		logger: logger.NewWithSubsystem("bot", "handler"),
	}

	return h
}
