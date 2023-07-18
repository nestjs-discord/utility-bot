package archive

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/config"
)

const Archive = "archive"

var permission = config.ProtectedContentPermission

var Command = &discordgo.ApplicationCommand{
	Name:                     Archive,
	Description:              "Close and lock a forum post.",
	DefaultMemberPermissions: &permission,
}
