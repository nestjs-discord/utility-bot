package interaction

import "github.com/bwmarrin/discordgo"

func DefaultHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Content not found.",
		},
	})
	if err != nil {
		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
	}

	// Delete the slash command when it doesn't have any registered handler
	// _ = s.ApplicationCommandDelete(s.State.User.ID, config.GetConfig().GuildID, i.ID)
}
