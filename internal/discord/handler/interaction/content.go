package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/erosdesire/discord-nestjs-utility-bot/core/config"
	"github.com/rs/zerolog/log"
)

func ContentHandler(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	name := i.ApplicationCommandData().Name

	// Check cache if the tag exists
	cmd, cmdExist := config.GetConfig().Commands[name]
	if !cmdExist {
		return false
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cmd.Content,
		},
	})
	if err == nil {
		return true
	}

	log.Error().
		Err(err).
		Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", i.Member.User.ID).
		Msg("failed to respond to interaction")

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Something went wrong",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return false
}