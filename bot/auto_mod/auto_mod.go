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

	rules []*dgo.AutoModerationRule
}

func NewAutoMod(session *dgo.Session, guildId env.GuildId) (*AutoMod, error) {
	a := &AutoMod{
		logger:  logger.NewWithSubsystem("bot", "auto-mod"),
		session: session,
		guildId: guildId,
	}

	return a, nil
}

func (a *AutoMod) ExecuteBackgroundJob() error {
	a.logger.Debug("executing background job")
	return a.syncRules()
}

func (a *AutoMod) syncRules() error {
	rules, err := a.session.AutoModerationRules(a.guildId.String())
	if err != nil {
		return err
	}

	a.rules = rules

	return nil
}
