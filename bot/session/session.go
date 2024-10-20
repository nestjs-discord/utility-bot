package session

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/bot/permissions"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Session struct {
	discordCfg *env.DiscordConfig
	logger     *slog.Logger
	session    *dgo.Session
}

// NewSession creates a new Discord session with the provided token.
func NewSession(discordCfg *env.DiscordConfig) (*Session, error) {
	session, err := dgo.New("Bot " + discordCfg.Token)
	if err != nil {
		return nil, fmt.Errorf("unable to create the session: %v", err)
	}
	setIntents(session)
	s := &Session{
		discordCfg: discordCfg,
		logger:     logger.NewWithSubsystem("bot", "session"),
		session:    session,
	}
	return s, nil
}

func setIntents(session *dgo.Session) {
	session.Identify.Intents = permissions.BotIntents
}

func (s *Session) ApplyHandler(h *handler.Handler) {
	s.session.AddHandler(h.Ready)
	s.session.AddHandler(h.InteractionCreate)
	s.session.AddHandler(h.MessageCreate)
}

func ProvideSession(s *Session) *dgo.Session {
	return s.session
}

// OpenWebsocketConnection creates a websocket connection to Discord.
// See: https://discord.com/developers/docs/topics/gateway#connecting
func (s *Session) OpenWebsocketConnection() error {
	err := s.session.Open()
	if err != nil {
		return fmt.Errorf("unable to open the session: %v", err)
	}

	s.logger.Info("session opened")

	return nil
}

// Close closes the websocket and stops all listening/heartbeat goroutines.
func (s *Session) Close() error {
	err := s.session.Close()
	if err != nil {
		return fmt.Errorf("unable to close the session: %v", err)
	}

	s.logger.Info("session closed")

	return nil
}
