package interaction

import (
	"github.com/bwmarrin/discordgo"
)

func NpmInspectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: logic
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Something went wrong",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
