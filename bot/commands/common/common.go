package common

import "github.com/bwmarrin/discordgo"

const (
	OptionTarget = "mention"
	OptionHide   = "hide"
)

var TargetOption = &discordgo.ApplicationCommandOption{
	Name:        OptionTarget,
	Description: "User a mention",
	Type:        discordgo.ApplicationCommandOptionUser,
}

var HideOption = &discordgo.ApplicationCommandOption{
	Name:        OptionHide,
	Description: "Make the command output only visible to you!",
	Type:        discordgo.ApplicationCommandOptionBoolean,
}
