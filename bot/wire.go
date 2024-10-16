package bot

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/google/wire"
)

func ProvideSession(b *Bot) *dgo.Session {
	return b.session
}

var Set = wire.NewSet(
	NewBot,
	ProvideSession,
)
