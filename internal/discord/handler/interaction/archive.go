package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/pkg/errors"
)

func ArchiveHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	currentChannelInfo, isValid := validateInteractionForThreadPost(s, i)
	if !isValid {
		return
	}

	content := "This post has been marked as \"archived\".\n" +
		"Please use it as a reference, but do not re-open it. " +
		"If you have a similar issue and cannot resolve it after reading this thread, please open a new post."

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
	if err != nil {
		util.InteractionRespondError(errors.Wrap(err, "failed to respond to interaction"), s, i)
		return
	}

	// Default options
	archived := true
	locked := true

	_, _ = s.ChannelEdit(currentChannelInfo.ID, &discordgo.ChannelEdit{
		Archived: &archived,
		Locked:   &locked,
	})
}
