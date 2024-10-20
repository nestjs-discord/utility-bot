package interaction

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"log/slog"
)

func (h *Handler) ModalSubmit(s *dgo.Session, i *dgo.InteractionCreate) {
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

	switch customId.Action {
	case forms.Modal:
		err = h.opts.Forms.ModalSubmitted(s, i, customId)
		if err != nil {
			h.respondError(err, s, i)
		}
	case forms.ModeratorBanConfirmModal:
		err = h.opts.Forms.ModBanModalSubmit(s, i, customId)
		if err != nil {
			h.respondError(err, s, i)
		}
		return
	}
}
