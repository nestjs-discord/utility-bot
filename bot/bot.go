package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/logger"
	"log/slog"
)

type Bot struct {
	cfg     *config.BotConfig
	session *discordgo.Session
	logger  *slog.Logger
}

func NewBot(cfg *config.BotConfig, handler *handler.Handler) (*Bot, error) {
	bot := &Bot{
		cfg:    cfg,
		logger: logger.NewWithSubsystem("bot"),
	}

	err := bot.newSession()
	if err != nil {
		return nil, err
	}

	bot.applyHandler(handler)

	return bot, nil
}

func (b *Bot) applyHandler(h *handler.Handler) {
	b.session.AddHandler(h.Ready)
	b.session.AddHandler(h.InteractionCreate)
	b.session.AddHandler(h.MessageCreate)
}
