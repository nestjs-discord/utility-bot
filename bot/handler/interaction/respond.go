package interaction

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"log/slog"
)

//func (h *Handler) RespondWithMessage(s *dgo.Session, i *dgo.InteractionCreate, message string) {
//	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
//		Type: dgo.InteractionResponseChannelMessageWithSource,
//		Data: &dgo.InteractionResponseData{
//			Content: message,
//			Flags:   dgo.MessageFlagsEphemeral,
//		},
//	})
//}

func (h *Handler) respondError(err error, s *dgo.Session, i *dgo.InteractionCreate) {
	h.logger.Error("interaction respond error",
		slog.String("reason", err.Error()),
		slog.Any("interaction", i),
	)

	content := "Something went wrong."
	if h.moderators.IsUserModerator(i.Member.User.ID) {
		content += fmt.Sprintf("\nHere's the internal error message: 🪲\n```\n%s\n```", err.Error())
	}

	e := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
			Flags:   dgo.MessageFlagsEphemeral,
		},
	})
	if e == nil {
		return
	}

	h.logger.Error("interaction respond error failed",
		slog.String("reason", e.Error()),
	)
}
