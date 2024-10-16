package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
)

func ArchiveHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, isValid := validateInteractionForThreadPost(s, i)
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
		util.InteractionRespondError(fmt.Errorf("failed to respond to interaction: %s", err), s, i)
		return
	}

	archived := true
	locked := true

	_, _ = s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
		Archived: &archived,
		Locked:   &locked,
	})
}
