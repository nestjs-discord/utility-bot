package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/npm"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/interaction"
	"github.com/rs/zerolog/log"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	log.Debug().Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", i.Member.User.ID).
		Msg("event: interaction create")

	if interaction.ContentHandler(s, i) {
		return
	}

	switch name {
	// Npm subcommand
	case npm.Name:
		for _, option := range i.ApplicationCommandData().Options {

			switch option.Name {
			case npm.SearchCommandName:
				interaction.NpmSearchHandler(s, i)
				return
			case npm.InspectCommandName:
				interaction.NpmInspectHandler(s, i)
				return
			}
		}
	}

	interaction.DefaultHandler(s, i)
}
