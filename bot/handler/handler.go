package handler

import (
	"github.com/nestjs-discord/utility-bot/bot/automod"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/handler/interaction"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/moderator"
	"github.com/nestjs-discord/utility-bot/logger"
	"github.com/nestjs-discord/utility-bot/pkg/rate_limit"
	"log/slog"
)

type Handler struct {
	logger             *slog.Logger
	interactionHandler *interaction.Handler
	rateLimit          *rate_limit.TTLMap // for the application commands
	autoMod            *automod.AutoMod
	forms              *forms.Forms
	markdown           *markdown.Markdown
	moderator          *moderator.Moderator
	// TODO: dependencies...
}

func NewHandler(
	interactionHandler *interaction.Handler,
	rateLimit *rate_limit.TTLMap,
	autoMod *automod.AutoMod,
	forms *forms.Forms,
	markdown *markdown.Markdown,
	moderator *moderator.Moderator,
) *Handler {
	h := &Handler{
		logger:             logger.NewWithSubsystem("bot", "handler"),
		interactionHandler: interactionHandler,
		rateLimit:          rateLimit,
		autoMod:            autoMod,
		forms:              forms,
		markdown:           markdown,
		moderator:          moderator,
	}

	return h
}
