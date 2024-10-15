package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/stats"

	"log/slog"
)

func (h *Handler) MessageCreate(s *discordgo.Session, i *discordgo.MessageCreate) {
	if i.Message.Author.Bot {
		return
	}

	h.logger.Debug("message create",
		slog.String("id", i.Message.ID),
		slog.String("content", i.Message.Content),
	)

	if h.autoMod.Enabled() {
		h.autoMod.Handler(s, i)
	}

	if !h.isModerator(i.Author.ID) {
		return
	}

	switch i.Content {
	case "!automod":
		h.autoMod.TrackHandler(s, i)
		return
	case "!stats":
		stats.Handler(s, i)
		return
	}
}
