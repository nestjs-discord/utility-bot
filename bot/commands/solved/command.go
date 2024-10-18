package solved

import dgo "github.com/bwmarrin/discordgo"

func (c *Solved) Command() *dgo.ApplicationCommand {
	return &dgo.ApplicationCommand{
		Name:        Name,
		Description: c.cfg.Description,
		Options: []*dgo.ApplicationCommandOption{
			{
				Required:    false,
				Name:        AutoClose,
				Description: c.cfg.Description,
				Type:        dgo.ApplicationCommandOptionInteger,
				Choices: []*dgo.ApplicationCommandOptionChoice{
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
}
