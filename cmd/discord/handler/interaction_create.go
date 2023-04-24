package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/erosdesire/discord-nestjs-utility-bot/config"
	"github.com/rs/zerolog/log"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	log.Debug().
		Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", i.Member.User.ID).
		Msg("event: interaction create")

	// check cache if the tag exists
	cmd, cmdExist := config.GetConfig().Commands[name]
	if cmdExist {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: cmd.Content,
			},
		})
		if err != nil {
			log.Error().
				Err(err).
				Str("name", name).
				Str("guild-id", i.GuildID).
				Str("channel-id", i.ChannelID).
				Str("user-id", i.Member.User.ID).
				Msg("failed to respond to interaction")
		}
		return
	}

	// Delete the slash if it doesn't have any registered handler

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Content not found.",
		},
	})
	_ = s.ApplicationCommandDelete(s.State.User.ID, config.GetConfig().GuildID, i.ID)
}
