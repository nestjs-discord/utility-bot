package credits

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
)

const Name = "credits"

type Options struct {
	Moderators *moderators.Moderators
}

type Credits struct {
	opts Options
}

func NewCredits() *Credits {
	return &Credits{}
}

func (c *Credits) Command() *dgo.ApplicationCommand {
	return &dgo.ApplicationCommand{
		Name:        Name,
		Description: "Bot developer credits",
	}
}

func (c *Credits) Handler(s *dgo.Session, i *dgo.InteractionCreate) error {
	text := `**Credits** 🎉

This bot has been carefully crafted and continuously maintained by <@%s>, focusing on serving the NestJS community.

💻 **GitHub:** [github.com/nestjs-discord/utility-bot](https://github.com/nestjs-discord/utility-bot) ⭐️

A big thank you to everyone on the NestJS team and the Discord server moderators for their valuable feedback, innovative ideas, and unwavering support! ❤️

If you're looking for a custom bot like this or have any suggestions for improvements, don't hesitate to reach out to the developer!`
	id, err := c.opts.Moderators.DecodeUserId("MzcwOTAzNzAwOTc5MzE4Nzg0")
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: fmt.Sprintf(text, id),
		},
	})
}
