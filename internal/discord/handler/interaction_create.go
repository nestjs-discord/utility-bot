package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/erosdesire/discord-nestjs-utility-bot/internal/discord/command/npm"
	interaction2 "github.com/erosdesire/discord-nestjs-utility-bot/internal/discord/handler/interaction"
	"github.com/rs/zerolog/log"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	log.Debug().Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", i.Member.User.ID).
		Msg("event: interaction create")

	if interaction2.ContentHandler(s, i) {
		return
	}

	switch name {
	case npm.SearchCommandName:
		interaction2.NpmSearchHandler(s, i)
		return
	}

	interaction2.DefaultHandler(s, i)
}
