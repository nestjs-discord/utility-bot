package common

import (
	dgo "github.com/bwmarrin/discordgo"
)

const (
	OptionTarget = "mention"
	OptionHide   = "hide"
)

var TargetOption = &dgo.ApplicationCommandOption{
	Name:        OptionTarget,
	Description: "User a mention",
	Type:        dgo.ApplicationCommandOptionUser,
}

var HideOption = &dgo.ApplicationCommandOption{
	Name:        OptionHide,
	Description: "Make the command output only visible to you!",
	Type:        dgo.ApplicationCommandOptionBoolean,
}
