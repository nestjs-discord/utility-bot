package handler

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/stats"

	"log/slog"
)

func (h *Handler) MessageCreate(s *dgo.Session, i *dgo.MessageCreate) {
	message := i.Message
	if message.Author.Bot {
		return
	}

	h.logger.Debug("message create",
		slog.String("id", message.ID),
		slog.String("content", message.Content),
	)

	if h.opts.Antispam.Enabled() {
		h.opts.Antispam.Handler(s, i)
	}

	if check := h.opts.AIForums.Handle(s, i); check {
		return
	}

	if !h.opts.Moderators.IsUserModerator(i.Author.ID) {
		return
	}

	switch i.Content {
	case "!antispam":
		h.opts.Antispam.TrackHandler(s, i)
		return
	case "!stats":
		stats.Handler(s, i)
		return
	}
}
