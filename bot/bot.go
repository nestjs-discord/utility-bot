package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Options struct {
	BotToken string
}

type Bot struct {
	opts    Options
	session *discordgo.Session
}

func NewBot(opts Options) (*Bot, error) {
	dg, err := discordgo.New("Bot " + opts.BotToken)
	if err != nil {
		return nil, fmt.Errorf("unable to create the discord session: %v", err)
	}

	return &Bot{
		opts:    opts,
		session: dg,
	}, nil
}
