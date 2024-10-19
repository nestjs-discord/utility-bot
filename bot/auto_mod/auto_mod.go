package auto_mod

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type AutoMod struct {
	logger  *slog.Logger
	session *dgo.Session
	guildId env.GuildId

	keywordFilter map[string]string
	regexPatterns map[string]string
}

func NewAutoMod(session *dgo.Session, guildId env.GuildId) (*AutoMod, error) {
	a := &AutoMod{
		logger:  logger.NewWithSubsystem("bot", "auto-mod"),
		session: session,
		guildId: guildId,
	}

	a.emptyCachedRules()

	return a, nil
}

func (a *AutoMod) emptyCachedRules() {
	a.keywordFilter = make(map[string]string)
	a.regexPatterns = make(map[string]string)
}
