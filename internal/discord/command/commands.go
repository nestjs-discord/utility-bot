package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/npm"
)

var (
	StaticCommands = []*discordgo.ApplicationCommand{
		npm.Subcommand,
	}
	// RegisteredCommands stores both static and dynamic commands
	// that can be easily configured after the bot is launched
	RegisteredCommands []*discordgo.ApplicationCommand
)
