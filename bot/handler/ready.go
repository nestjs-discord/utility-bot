package handler

import (
	dgo "github.com/bwmarrin/discordgo"
	"log/slog"
)

func (h *Handler) Ready(_ *dgo.Session, r *dgo.Ready) {
	h.logger.Info("ready",
		slog.String("id", r.User.ID),
		slog.String("user", r.User.String()),
	)
}
