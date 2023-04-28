package stats

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/config"
)

const Stats = "stats"

var permission = config.ProtectedContentPermission

var Command = &discordgo.ApplicationCommand{
	Name:                     Stats,
	Description:              "General statistics about the bot instance",
	DefaultMemberPermissions: &permission,
}
