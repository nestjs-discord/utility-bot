package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/config"
)

func NewSession() (*discordgo.Session, error) {
	token := "Bot " + config.GetBotToken()

	dg, err := discordgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %v", err)
	}

	return dg, nil
}
