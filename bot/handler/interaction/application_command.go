package interaction

import (
	"github.com/nestjs-discord/utility-bot/bot/commands/credits"
	"log/slog"

	dgo "github.com/bwmarrin/discordgo"

	"github.com/nestjs-discord/utility-bot/bot/commands/archive"
	"github.com/nestjs-discord/utility-bot/bot/commands/dont_ping_mods"
	"github.com/nestjs-discord/utility-bot/bot/commands/google_it"
	"github.com/nestjs-discord/utility-bot/bot/commands/reference"
	"github.com/nestjs-discord/utility-bot/bot/commands/solved"
)

type applicationCommandHandlersMap map[string]func(s *dgo.Session, i *dgo.InteractionCreate) error

func (h *Handler) ApplicationCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	userID := i.Member.User.ID

	h.logger.Debug("application command",
		slog.String("userId", userID),
		slog.String("channelId", i.ChannelID),
		slog.String("name", data.Name),
	)

	if h.opts.RateLimit.CheckRateLimit(userID) {
		h.opts.RateLimit.ForbidInteraction(s, i)
		return
	}

	staticHandlers := applicationCommandHandlersMap{
		archive.Name:        h.opts.Archive.Handler,
		credits.Name:        h.opts.Credits.Handler,
		dont_ping_mods.Name: h.opts.DontPingMods.Handler,
		google_it.Name:      h.opts.GoogleIt.Handler,
		reference.Name:      h.opts.Reference.Handler,
		solved.Name:         h.opts.Solved.Handler,
		// ...
	}

	if handler, ok := staticHandlers[data.Name]; ok {
		err := handler(s, i)
		if err != nil {
			h.respondError(err, s, i)
		}
		return
	}

	// content handling should happen always at the end
	// because it automatically responds to the interaction
	err := h.opts.Markdown.ContentHandler(s, i)
	if err != nil {
		h.respondError(err, s, i)
	}
}
