package respond

import dgo "github.com/bwmarrin/discordgo"

func InteractionWithEphemeralMessage(s *dgo.Session, i *dgo.InteractionCreate, message string) {
	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: message,
			Flags:   dgo.MessageFlagsEphemeral,
		},
	})
}
