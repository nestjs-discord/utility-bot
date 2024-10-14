package handler

import (
	"github.com/nestjs-discord/utility-bot/bot/automod"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/logger"
	"github.com/nestjs-discord/utility-bot/pkg/rate_limit"
	"log/slog"
)

type Handler struct {
	logger    *slog.Logger
	rateLimit *rate_limit.TTLMap // for the application commands

	autoMod *automod.AutoMod
	forms   *forms.Forms
	// TODO: dependencies...
}

func NewHandler(
	rateLimit *rate_limit.TTLMap,
	autoMod *automod.AutoMod,
	forms *forms.Forms,
) *Handler {
	return &Handler{
		logger:    logger.NewWithSubsystem("bot", "handler"),
		rateLimit: rateLimit,
		autoMod:   autoMod,
		forms:     forms,
	}
}
