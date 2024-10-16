package archive

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/config"
)

const Name = "archive"

var permission = config.BotProtectedContentPermission

var Command = &discordgo.ApplicationCommand{
	Name:                     Name,
	Description:              "Close and lock a forum post.",
	DefaultMemberPermissions: &permission,
}
