package reference

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/common"
	"github.com/nestjs-discord/utility-bot/infra/services/algolia"
)

const Name = "reference"
const QueryOption = "query"

type Reference struct{}

func NewReference() *Reference {
	return &Reference{}
}

func (r *Reference) Command() *dgo.ApplicationCommand {
	options := []*dgo.ApplicationCommandOption{
		{
			Name:         QueryOption,
			Type:         dgo.ApplicationCommandOptionString,
			Description:  "The query to search for",
			Required:     true,
			Autocomplete: true,
		},
		common.TargetOption,
		common.HideOption,
	}

	subcommand := &dgo.ApplicationCommand{
		Name:        Name,
		Description: "reference related sub-commands",
		Options:     []*dgo.ApplicationCommandOption{},
	}

	for slug, app := range algolia.Apps {
		subcommand.Options = append(subcommand.Options, &dgo.ApplicationCommandOption{
			Type:        dgo.ApplicationCommandOptionSubCommand,
			Name:        slug,
			Description: "Display docs for " + string(app),
			Options:     options,
		})
	}

	return subcommand
}
