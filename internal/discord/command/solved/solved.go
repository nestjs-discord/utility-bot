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
			Description: "If no option is specified, the post will be closed after using the command.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Do not close",
					Value: 0,
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
