package solved

import (
	"github.com/bwmarrin/discordgo"
)

const Solved = "solved"

var Command = &discordgo.ApplicationCommand{
	Name:        Solved,
	Description: "Close and mark a forum post as solved.",
}
