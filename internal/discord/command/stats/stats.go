package stats

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/config"
)

const Name = "stats"

var permission = config.ProtectedContentPermission

var Command = &discordgo.ApplicationCommand{
	Name:                     Name,
	Description:              "General statistics about the bot instance",
	DefaultMemberPermissions: &permission,
}
