package npm

import "github.com/bwmarrin/discordgo"

const (
	SearchCommandName = "npm-search"
	SearchOptionName  = "text"
	SearchOptionSort  = "sort"
)

type SortOptionValue int64

const (
	SearchSortPopularity  SortOptionValue = 1
	SearchSortQuality     SortOptionValue = 2
	SearchSortMaintenance SortOptionValue = 3
)

var Search = &discordgo.ApplicationCommand{
	Name:        SearchCommandName,
	Description: "Search public packages on NPM registry",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        SearchOptionName,
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Package name (for example \"scope:nestjs\" to list all the official nestjs packages)",
			Required:    true,
		},
		{
			Name:        SearchOptionSort,
			Type:        discordgo.ApplicationCommandOptionInteger,
			Description: "Sort by (default optimal)",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "Popularity", Value: SearchSortPopularity},
				{Name: "Quality", Value: SearchSortQuality},
				{Name: "Maintenance", Value: SearchSortMaintenance},
			},
			Required: false,
		},
	},
}
