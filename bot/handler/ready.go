package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/bot/user"
	"log/slog"
)

func (h *Handler) Ready(s *discordgo.Session, r *discordgo.Ready) {
	h.logger.Info("ready",
		slog.String("id", r.User.ID),
		slog.String("username", user.FormatUsername(r.User)),
	)

	err := status.Update(s)
	if err != nil {
		h.logger.Error("update status failed",
			slog.Any("err", err),
		)
		return
	}

	h.logger.Info("status updated")
}
