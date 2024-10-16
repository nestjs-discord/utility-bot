package interaction

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/nestjs-discord/utility-bot/bot/command/archive"
	"github.com/nestjs-discord/utility-bot/bot/command/dont_ping_mods"
	"github.com/nestjs-discord/utility-bot/bot/command/google_it"
	"github.com/nestjs-discord/utility-bot/bot/command/reference"
	"github.com/nestjs-discord/utility-bot/bot/command/solved"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/interaction"
)

func (h *Handler) ApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	userID := i.Member.User.ID

	h.logger.Debug("application command",
		slog.String("userId", userID),
		slog.String("channelId", i.ChannelID),
		slog.String("name", data.Name),
		slog.Any("options", i.ApplicationCommandData().Options),
	)

	if h.rateLimit.CheckRateLimit(userID) {
		h.rateLimit.ForbidInteraction(s, i)
		return
	}

	switch data.Name {
	case solved.Name:
		interaction.SolvedHandler(s, i) // TODO: replace this with the instantiated one!
		return
	case archive.Name:
		interaction.ArchiveHandler(s, i)
		return
	case reference.Name:
		reference.Handler(s, i)
		return
	case google_it.Name:
		google_it.Handler(s, i)
		return
	case dont_ping_mods.Name:
		dont_ping_mods.Handler(s, i)
		return
	}

	if h.markdown.ContentHandler(s, i) {
		return
	}

	h.applicationCommandUnknownHandler(s, i)
}

func (h *Handler) applicationCommandUnknownHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.logger.Error("unknown application command",
		slog.Any("interaction", i),
	)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Unknown application command.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
