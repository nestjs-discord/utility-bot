package google_it

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/pkg/search"
)

const (
	Name         = "google-it"
	OptionSuffix = "-suggestion"
)

var searchClient = search.NewSearch()

func init() {
	elements := []string{"first", "second", "third", "fourth"}
	minLength := 3

	for i, opt := range elements {
		Command.Options = append(Command.Options, &discordgo.ApplicationCommandOption{
			Name:         opt + OptionSuffix,
			Description:  opt + " suggestion",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     i == 0, // only the first element is required
			Autocomplete: true,
			MinLength:    &minLength,
		})
	}
}

var Command = &discordgo.ApplicationCommand{
	Name:        Name,
	Description: "Telling people to search Google/StackOverflow before asking on our server",
	Type:        discordgo.ChatApplicationCommand,
}
