package handler

import (
	"github.com/nestjs-discord/utility-bot/bot/automod"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/logger"
	"github.com/nestjs-discord/utility-bot/pkg/rate_limit"
	"log/slog"
)

type Handler struct {
	logger     *slog.Logger
	moderators []string
	rateLimit  *rate_limit.TTLMap // for the application commands

	autoMod  *automod.AutoMod
	forms    *forms.Forms
	markdown *markdown.Markdown
	// TODO: dependencies...
}

func NewHandler(
	moderators []string,
	rateLimit *rate_limit.TTLMap,
	autoMod *automod.AutoMod,
	forms *forms.Forms,
	markdown *markdown.Markdown,
) *Handler {
	return &Handler{
		logger:     logger.NewWithSubsystem("bot", "handler"),
		moderators: moderators,
		rateLimit:  rateLimit,
		autoMod:    autoMod,
		forms:      forms,
		markdown:   markdown,
	}
}

func (h *Handler) isModerator(userId string) bool {
	for _, id := range h.moderators {
		if id == userId {
			return true
		}
	}

	return false
}
