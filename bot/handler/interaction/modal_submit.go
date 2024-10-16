package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"log/slog"
)

func (h *Handler) ModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	if data.CustomID == "" {
		h.logger.Error("received a modal submit with empty custom id",
			slog.Any("interaction", i),
		)
		return
	}

	customId, err := components.DecodeCustomId(data.CustomID)
	if err != nil {
		h.logger.Error("failed to decode custom id",
			slog.String("customId", data.CustomID),
			slog.Any("interaction", i),
		)
		return
	}

	// this may be replaced this with a switch statement in the future..!
	if customId.Action != forms.Modal {
		return
	}

	err = h.forms.ModalSubmitted(s, i, customId)
	if err != nil {
		h.respondError(err, s, i)
	}
}
