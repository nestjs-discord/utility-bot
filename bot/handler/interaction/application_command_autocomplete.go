package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/command/google_it"
	"github.com/nestjs-discord/utility-bot/bot/command/reference"
	"log/slog"
)

func (h *Handler) ApplicationCommandAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	h.logger.Debug("interaction application command autocomplete",
		slog.String("name", data.Name),
		slog.Any("options", data.Options),
	)

	switch data.Name {
	case reference.Name:
		reference.AutocompleteHandler(s, i)
		return
	case google_it.Name:
		google_it.AutocompleteHandler(s, i)
		return
	}
}
