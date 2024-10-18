package dont_ping_mods

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/common"
	"github.com/nestjs-discord/utility-bot/bot/permissions"
)

const Name = "dont-ping-mods"

func (d *DontPingMods) Command() *discordgo.ApplicationCommand {
	perm := int64(permissions.ProtectedCommands)

	return &discordgo.ApplicationCommand{
		Name:                     Name,
		Description:              "Tell someone to stop pinging mods for help",
		DefaultMemberPermissions: &perm,
		Options: []*discordgo.ApplicationCommandOption{
			common.TargetOption,
		},
	}
}
