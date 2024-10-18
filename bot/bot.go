package bot

import (
	"github.com/nestjs-discord/utility-bot/bot/session"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Bot struct {
	logger  *slog.Logger
	session *session.Session
}

func NewBot(discordCfg *env.DiscordConfig, _ *logger.Logger) (*Bot, error) {
	s, err := session.NewSession(discordCfg)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		logger:  logger.NewWithSubsystem("bot"),
		session: s,
	}

	return b, nil
}

func ProvideSession(b *Bot) *session.Session {
	return b.session
}

func (b *Bot) Open() error {
	return b.session.OpenWebsocketConnection()
}

func (b *Bot) Close() error {
	return b.session.Close()
}
