package handler

import (
	"github.com/nestjs-discord/utility-bot/bot/automod"
	"github.com/nestjs-discord/utility-bot/logger"
	"log/slog"
)

type Handler struct {
	logger  *slog.Logger
	autoMod *automod.AutoMod

	// TODO: dependencies...
}

func NewHandler(
	autoMod *automod.AutoMod,
) *Handler {
	return &Handler{
		logger:  logger.NewWithSubsystem("bot", "handler"),
		autoMod: autoMod,
	}
}
