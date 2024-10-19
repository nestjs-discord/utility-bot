package interaction

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/google_it"
	"github.com/nestjs-discord/utility-bot/bot/commands/reference"
	"log/slog"
)

func (h *Handler) ApplicationCommandAutocomplete(s *dgo.Session, i *dgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	h.logger.Debug("interaction application command autocomplete",
		slog.String("name", data.Name),
	)

	switch data.Name {
	case reference.Name:
		h.opts.Reference.AutocompleteHandler(s, i)
		return
	case google_it.Name:
		h.opts.GoogleIt.AutocompleteHandler(s, i)
		return
	}
}
