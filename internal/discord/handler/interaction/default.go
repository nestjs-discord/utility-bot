package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
)

func DefaultHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Content not found.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		util.InteractionRespondError(err, s, i)
	}

	// Delete the slash command when it doesn't have any registered handler
	log.Warn().
		Str("app-id", i.AppID).
		Str("guild-id", i.GuildID).
		Str("id", i.ID).
		Str("interaction-id", i.Interaction.ID).
		Interface("interaction-data", i.Data).
		Msg("received unknown slash command, consider running the discord:clean command.")

	//err = s.ApplicationCommandDelete(i.AppID, i.GuildID, i.ID)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to delete default command")
	//}
}
