package interaction

import (
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

	if h.rateLimit.CheckRateLimit(userID) {
		h.rateLimit.ForbidInteraction(s, i)
		return
	}

	staticHandlers := applicationCommandHandlersMap{
		solved.Name:  h.solved.Handler,
		archive.Name: h.archive.Handler,
		// ...
	}

	if handler, ok := staticHandlers[data.Name]; ok {
		err := handler(s, i)
		if err != nil {
			h.respondError(err, s, i)
		}
		return
	}

	switch data.Name {
	case reference.Name:
		reference.Handler(s, i)
		return
	case google_it.Name:
		google_it.Handler(s, i)
		return
	case dont_ping_mods.Name:
		h.dontPingMods.Handler(s, i)
		return
	}

	if h.markdown.ContentHandler(s, i) {
		return
	}

	h.applicationCommandUnknownHandler(s, i)
}

func (h *Handler) applicationCommandUnknownHandler(s *dgo.Session, i *dgo.InteractionCreate) {
	h.logger.Error("unknown application command",
		slog.Any("interaction", *i),
	)

	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "Unknown application command.",
			Flags:   dgo.MessageFlagsEphemeral,
		},
	})
}
