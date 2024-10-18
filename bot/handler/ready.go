package handler

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func (h *Handler) Ready(_ *discordgo.Session, r *discordgo.Ready) {
	h.logger.Info("ready",
		slog.String("id", r.User.ID),
		slog.String("user", r.User.String()),
	)

	h.status.StartUpdatingInBackground()
}
