package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"log"
	"log/slog"
)

func (h *Handler) Ready(s *discordgo.Session, m *discordgo.Ready) {
	h.logger.Info("ready",
		slog.String("id", m.User.ID),
		slog.String("username", util.FormatUsername(m.User)),
		// TODO: put a breakpoint and see what else is available to print
	)

	err := status.Update(s)
	if err != nil {
		log.Fatalf("error updating status: %v", err)
	}

	h.logger.Info("status updated")
}
