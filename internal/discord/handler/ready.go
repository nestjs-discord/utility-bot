package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Info().
		Str("id", s.State.User.ID).
		Str("username", fmt.Sprintf("%s#%s", m.User.Username, m.User.Discriminator)).
		Msg("logged in as")

	if err := updateStatus(s); err != nil {
		log.Panic().Err(err).Msg("failed to update status")
		return
	}
	log.Debug().Msg("status updated")

	log.Info().Msg("ready")
}

func updateStatus(s *discordgo.Session) error {
	return s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "slash commands",
				Type: discordgo.ActivityTypeListening,
			},
		},
	})
}
