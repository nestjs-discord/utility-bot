package google_it

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/services/google"
)

const (
	Name         = "google-it"
	OptionSuffix = "-suggestion"
)

var googleClient = google.NewGoogle() // TODO: don't use global instance

func init() {
	elements := []string{"first", "second", "third", "fourth"}
	minLength := 3

	for i, opt := range elements {
		Command.Options = append(Command.Options, &discordgo.ApplicationCommandOption{
			Name:         opt + OptionSuffix,
			Description:  opt + " suggestion",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     i == 0, // only the first item is required
			Autocomplete: true,
			MinLength:    &minLength,
		})
	}
}

var Command = &discordgo.ApplicationCommand{
	Name:        Name,
	Description: "Tell someone to search on Google or StackOverflow!",
	Type:        discordgo.ChatApplicationCommand,
}
