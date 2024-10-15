package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func UnknownHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Unknown slash command.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	// Delete the slash command when it doesn't have any registered handler
	log.Error().
		Str("app-id", i.AppID).
		Str("guild-id", i.GuildID).
		Str("id", i.ID).
		Str("interaction-id", i.Interaction.ID).
		Interface("interaction-data", i.Data).
		Msg("received unknown slash command, consider running the discord:clean command.")
}
