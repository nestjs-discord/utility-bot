package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/config/env"
	"github.com/nestjs-discord/utility-bot/logger"
	"log/slog"
)

type Bot struct {
	discordCfg *env.DiscordConfig
	session    *discordgo.Session
	logger     *slog.Logger
}

func NewBot(discordCfg *env.DiscordConfig, handler *handler.Handler) (*Bot, error) {
	bot := &Bot{
		discordCfg: discordCfg,
		logger:     logger.NewWithSubsystem("bot"),
	}

	err := bot.newSession()
	if err != nil {
		return nil, err
	}

	bot.applyHandler(handler)

	err = bot.logServerInviteLink()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (b *Bot) applyHandler(h *handler.Handler) {
	b.session.AddHandler(h.Ready)
	b.session.AddHandler(h.InteractionCreate)
	b.session.AddHandler(h.MessageCreate)
}
