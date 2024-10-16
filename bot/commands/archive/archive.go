package archive

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/permissions"
)

const Name = "archive"

var perm = int64(permissions.ProtectedCommands)

var Command = &discordgo.ApplicationCommand{
	Name:                     Name,
	Description:              "Close and lock a forum post.",
	DefaultMemberPermissions: &perm,
}
