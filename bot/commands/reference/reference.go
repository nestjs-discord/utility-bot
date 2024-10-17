package reference

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/common"
	"github.com/nestjs-discord/utility-bot/infra/services/algolia"
)

const Name = "reference"
const QueryOption = "query"

type Reference struct{}

func New() *Reference {
	return &Reference{}
}

func (r *Reference) Command() *discordgo.ApplicationCommand {
	options := []*discordgo.ApplicationCommandOption{
		{
			Name:         QueryOption,
			Type:         discordgo.ApplicationCommandOptionString,
			Description:  "The query to search for",
			Required:     true,
			Autocomplete: true,
		},
		common.TargetOption,
		common.HideOption,
	}

	subcommand := &discordgo.ApplicationCommand{
		Name:        Name,
		Description: "reference related sub-commands",
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	for slug, app := range algolia.Apps {
		subcommand.Options = append(subcommand.Options, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        slug,
			Description: "Display docs for " + string(app),
			Options:     options,
		})
	}

	return subcommand
}
