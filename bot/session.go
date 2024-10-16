package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// newSession creates a new Discord session with the provided token.
func (b *Bot) newSession() error {
	session, err := discordgo.New("Bot " + b.discordCfg.Token)
	if err != nil {
		return fmt.Errorf("unable to create the session: %v", err)
	}

	// set intents
	session.Identify.Intents = intents

	b.session = session

	return nil
}

// OpenWebsocketConnection creates a websocket connection to Discord.
// See: https://discord.com/developers/docs/topics/gateway#connecting
func (b *Bot) OpenWebsocketConnection() error {
	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("unable to open the session: %v", err)
	}

	b.logger.Info("session opened")

	return nil
}

// Close closes the websocket and stops all listening/heartbeat goroutines.
func (b *Bot) Close() error {
	err := b.session.Close()
	if err != nil {
		return fmt.Errorf("unable to close the session: %v", err)
	}

	b.logger.Info("session closed")

	return nil
}
