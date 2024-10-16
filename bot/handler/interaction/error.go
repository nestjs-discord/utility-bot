package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func (h *Handler) RespondError(err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.logger.Error("interaction respond error",
		slog.String("reason", err.Error()),
		slog.Any("interaction", i),
	)

	content := "Something went wrong."
	if h.moderators.IsUserModerator(i.Member.User.ID) {
		content += fmt.Sprintf("\nHere's the internal error message: 🪲\n```\n%s\n```", err.Error())
	}

	e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if e == nil {
		return
	}

	h.logger.Error("interaction respond error failed",
		slog.String("reason", e.Error()),
	)
}
