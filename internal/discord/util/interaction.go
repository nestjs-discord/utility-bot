package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func InteractionRespondError(err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := "Something went wrong."
	if IsUserModerator(i.User.ID) {
		content += "\nBecause you're marked as a moderator, here's the internal error message:\n\n" + err.Error()
	}

	e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if e != nil {
		log.Error().Err(e).Msg("failed to respond interaction error")
	}
}
