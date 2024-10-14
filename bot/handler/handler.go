package handler

import (
	"github.com/nestjs-discord/utility-bot/internal/logger"
	"log/slog"
)

type Handler struct {
	logger *slog.Logger
	// TODO: dependencies...
}

func NewHandler() *Handler {
	return &Handler{
		logger: logger.NewWithSubsystem("bot", "handler"),
	}
}
