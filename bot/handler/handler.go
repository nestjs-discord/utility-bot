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

type Handler struct {
	logger             *slog.Logger
	interactionHandler *interaction.Handler
	antispam           *antispam.Antispam
	forms              *forms.Forms
	markdown           *markdown.Markdown
	moderators         *moderators.Moderators
}

func NewHandler(
	interactionHandler *interaction.Handler,
	antispam *antispam.Antispam,
	forms *forms.Forms,
	markdown *markdown.Markdown,
	moderators *moderators.Moderators,
) *Handler {
	h := &Handler{
		logger:             logger.NewWithSubsystem("bot", "handler"),
		interactionHandler: interactionHandler,
		antispam:           antispam,
		forms:              forms,
		markdown:           markdown,
		moderators:         moderators,
	}

	return h
}
