package auto_mod

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Options struct {
	Session *dgo.Session
	GuildId env.GuildId
}

type AutoMod struct {
	opts   Options
	logger *slog.Logger

	keywordFilter map[string]string
	regexPatterns map[string]string
}

func NewAutoMod(opts Options) (*AutoMod, error) {
	a := &AutoMod{
		logger: logger.NewWithSubsystem("bot", "auto-mod"),
		opts:   opts,
	}

	a.emptyCachedRules()

	return a, nil
}

func (a *AutoMod) emptyCachedRules() {
	a.keywordFilter = make(map[string]string)
	a.regexPatterns = make(map[string]string)
}
