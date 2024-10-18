package credits

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
)

const Name = "credits"

type Credits struct{}

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
	respond.InteractionWithEphemeralMessage(s, i, "TBD") // TODO: logic
	return nil
}
