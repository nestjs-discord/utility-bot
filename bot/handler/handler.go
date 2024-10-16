package handler

import (
	"github.com/nestjs-discord/utility-bot/bot/automod"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/handler/interaction"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/moderator"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Handler struct {
	logger             *slog.Logger
	interactionHandler *interaction.Handler
	antispam           *automod.Antispam
	forms              *forms.Forms
	markdown           *markdown.Markdown
	moderator          *moderator.Moderator
	// TODO: dependencies...
}

func NewHandler(
	interactionHandler *interaction.Handler,
	antispam *automod.Antispam,
	forms *forms.Forms,
	markdown *markdown.Markdown,
	moderator *moderator.Moderator,
) *Handler {
	h := &Handler{
		logger:             logger.NewWithSubsystem("bot", "handler"),
		interactionHandler: interactionHandler,
		antispam:           antispam,
		forms:              forms,
		markdown:           markdown,
		moderator:          moderator,
	}

	return h
}
