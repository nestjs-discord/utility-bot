package reference

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/algolia"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/common"
)

const Name = "reference"
const QueryOption = "query"

var options = []*discordgo.ApplicationCommandOption{
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

var Subcommand = &discordgo.ApplicationCommand{
	Name:        Name,
	Description: "reference related sub-commands",
	Options:     []*discordgo.ApplicationCommandOption{},
}

func init() {
	for slug, app := range algolia.Apps {
		Subcommand.Options = append(Subcommand.Options, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        slug,
			Description: "Display docs for " + string(app),
			Options:     options,
		})
	}
}
