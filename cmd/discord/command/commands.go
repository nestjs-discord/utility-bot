package commands

import "github.com/bwmarrin/discordgo"

const (
	NpmInspect           = "npm-inspect"
	NpmInspectNameOption = "name"
)

var (
	StaticCommands = []*discordgo.ApplicationCommand{
		npmSearch,
	}
	// RegisteredCommands stores both static and dynamic commands
	// that can be easily configured after the bot is launched
	RegisteredCommands []*discordgo.ApplicationCommand
)
