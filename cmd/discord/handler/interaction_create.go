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
	if true {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: name,
			},
		})
	} else {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Content not found.",
			},
		})
		_ = s.ApplicationCommandDelete(s.State.User.ID, config.GetConfig().GuildID, i.ID)
	}
}
