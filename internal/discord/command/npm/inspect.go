package npm

import "github.com/bwmarrin/discordgo"

const (
	InspectCommandName   = "inspect"
	InspectOptionPackage = "package"
	InspectOptionVersion = "version"
)

var Inspect = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        InspectCommandName,
	Description: "Inspect a specific public package on NPM registry",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        InspectOptionPackage,
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Package name (for example \"@nestjs/common\")",
			Required:    true,
		},
		{
			Name:        InspectOptionVersion,
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specific version of the package (default latest)",
			Required:    false,
		},
	},
}
