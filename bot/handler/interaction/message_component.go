package interaction

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"log/slog"
)

type messageComponentHandlersMap map[string]func(*dgo.Session, *dgo.InteractionCreate, *components.CustomID) error

func (h *Handler) MessageComponent(s *dgo.Session, i *dgo.InteractionCreate) {
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

	messageComponentHandlers := messageComponentHandlersMap{
		forms.OpenModalButton:       h.forms.OpenModalButtonClicked,
		forms.ModeratorAcceptButton: h.forms.ModAcceptButtonClicked,
		forms.ModeratorRejectButton: h.forms.ModRejectButtonClicked,
		forms.ModeratorBanButton:    h.forms.ModBanButtonClicked,
	}

	handler, ok := messageComponentHandlers[customId.Action]
	if !ok {
		return
	}

	err = handler(s, i, customId)
	if err != nil {
		h.RespondError(err, s, i)
	}
}
