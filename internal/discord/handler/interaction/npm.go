package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/npm"
)

func NpmHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case npm.SearchCommandName:
			NpmSearchHandler(s, i)
			return
		case npm.InspectCommandName:
			NpmInspectHandler(s, i)
			return
		}
	}
}
