package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"log/slog"
)

func (h *Handler) MessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: refactor this, it should only be called for the moderator actions, not the public interactive button!
	if h.forms.RaceConditionCheck(i.Message.ID) {
		h.forms.RaceConditionRespond(s, i)
		return
	}

	data := i.MessageComponentData()

	if data.CustomID == "" {
		h.logger.Error("received a message component with empty custom id",
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
	case forms.ModeratorAcceptButton:
		err = h.forms.ModAcceptButtonClicked(s, i, customId)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return
	case forms.ModeratorRejectButton:
		err = h.forms.ModRejectButtonClicked(s, i, customId)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return
	case forms.ModeratorBanButton:
		err = h.forms.ModBanButtonClicked(s, i, customId)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return
	case forms.OpenModalButton:
		err = h.forms.OpenModalButtonClicked(s, i, customId)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return
	}
}
