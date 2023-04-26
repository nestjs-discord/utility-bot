package npm

import (
	"github.com/bwmarrin/discordgo"
)

const Name = "npm"

var Subcommand = &discordgo.ApplicationCommand{
	// When a command has subcommands/subcommand groups
	// It must not have top-level options, they aren't accessible in the UI
	// in this case (at least not yet), so if a command has
	// subcommands/subcommand any groups registering top-level options
	// will cause the registration of the command to fail
	Name:        Name,
	Description: "NPM registry related commands",
	Options: []*discordgo.ApplicationCommandOption{
		Search,
		Inspect,
	},
}
