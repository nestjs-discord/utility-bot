package common

import "github.com/bwmarrin/discordgo"

const (
	OptionTarget = "target"
	OptionHide   = "hide"
)

var TargetOption = &discordgo.ApplicationCommandOption{
	Name:        OptionTarget,
	Description: "User to mention",
	Type:        discordgo.ApplicationCommandOptionUser,
}

var HideOption = &discordgo.ApplicationCommandOption{
	Name:        OptionHide,
	Description: "Hide commands output",
	Type:        discordgo.ApplicationCommandOptionBoolean,
}
