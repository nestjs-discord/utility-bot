package commands

import "github.com/bwmarrin/discordgo"

const (
	NpmSearch           = "npm-search"
	NpmSearchNameOption = "name"
	NpmSearchSortOption = "sort"
)

type NpmSortOptionValue string

const (
	NpmSearchSortPopularity  NpmSortOptionValue = "popularity"
	NpmSearchSortQuality     NpmSortOptionValue = "quality"
	NpmSearchSortMaintenance NpmSortOptionValue = "maintenance"
)

var npmSearch = &discordgo.ApplicationCommand{
	Name:        NpmSearch,
	Description: "Search packages in NPM registry",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        NpmSearchNameOption,
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Package name (for example \"@nestjs/\" to list all the official nestjs packages)",
			Required:    true,
		},
		{
			Name:        NpmSearchSortOption,
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Sort by (default optimal)",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "Popularity", Value: NpmSearchSortPopularity},
				{Name: "Quality", Value: NpmSearchSortQuality},
				{Name: "Maintenance", Value: NpmSearchSortMaintenance},
			},
			Required: false,
		},
	},
}
