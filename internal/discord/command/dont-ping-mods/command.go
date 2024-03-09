package dont_ping_mods

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/common"
)

const Name = "dont-ping-mods"

var permission = config.BotProtectedContentPermission

var Command = &discordgo.ApplicationCommand{
	Name:                     Name,
	Description:              "Tell someone to stop pinging mods for help",
	DefaultMemberPermissions: &permission,
	Options: []*discordgo.ApplicationCommandOption{
		common.TargetOption,
	},
}
