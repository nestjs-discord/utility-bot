package solved

import (
	"github.com/bwmarrin/discordgo"
)

const (
	Solved    = "solved"
	AutoClose = "auto-close"
)

var Command = &discordgo.ApplicationCommand{
	Name:        Solved,
	Description: "Close and mark a forum post as solved.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Required:    false,
			Name:        AutoClose,
			Description: "If the 'auto-close' option isn't specified, the post will remain open after using the command.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Close right after",
					Value: 1,
				},
				{
					Name:  "In 1 hour",
					Value: 60,
				},
				{
					Name:  "In 24 hours",
					Value: 60 * 24,
				},
				{
					Name:  "In 3 days",
					Value: 60 * 24 * 3,
				},
				{
					Name:  "In a week",
					Value: 60 * 24 * 7,
				},
			},
		},
	},
}
