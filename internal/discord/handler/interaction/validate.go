package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
)

func validateInteractionForThreadPost(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.Channel, bool) {
	currentChannelInfo, err := s.Channel(i.ChannelID)
	if err != nil {
		util.InteractionRespondError(
			fmt.Errorf("failed to get current channel info: %s", err),
			s, i)

		return nil, false
	}

	if currentChannelInfo.Type != discordgo.ChannelTypeGuildPublicThread &&
		currentChannelInfo.Type != discordgo.ChannelTypeGuildPrivateThread {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":warning: You can only use this command in forum posts.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		return nil, false
	}

	if currentChannelInfo.ThreadMetadata.Locked {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":warning: Cannot perform this action on a locked thread post",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}

	return currentChannelInfo, true
}
