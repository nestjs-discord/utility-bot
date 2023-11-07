package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/cache"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/archive"
	google_it "github.com/nestjs-discord/utility-bot/internal/discord/command/google-it"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/reference"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/solved"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/stats"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/interaction"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleInteractionApplicationCommand(s, i)
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
		handleInteractionApplicationCommandAutocomplete(s, i)
	}
}

func handleInteractionApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	userID := i.Member.User.ID

	log.Debug().Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", userID).
		Interface("options", i.ApplicationCommandData().Options).
		Msg("event: interaction application command")

	if checkRatelimit(userID) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: config.GetConfig().Ratelimit.Message,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
		return
	}

	if interaction.ContentHandler(s, i) {
		return
	}

	switch name {
	case archive.Archive:
		interaction.ArchiveHandler(s, i)
		return

	case reference.Name:
		interaction.ReferenceHandler(s, i)
		return

	case solved.Solved:
		interaction.SolvedHandler(s, i)
		return

	case stats.Stats:
		interaction.StatHandler(s, i)
		return

	case google_it.Name:
		google_it.Handler(s, i)
		return
	}

	interaction.DefaultHandler(s, i)
}

func handleInteractionApplicationCommandAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	log.Debug().Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Interface("options", i.ApplicationCommandData().Options).
		Msg("event: interaction application command autocomplete")

	switch name {
	case reference.Name:
		interaction.ReferenceAutocompleteHandler(s, i)
		return
	case google_it.Name:
		google_it.AutocompleteHandler(s, i)
		return
	}
}

func checkRatelimit(userID string) bool {
	maxUsage := config.GetConfig().Ratelimit.Usage
	if !util.IsUserModerator(userID) {
		cache.Ratelimit.IncrementUsage(userID)

		if cache.Ratelimit.GetUsageCount(userID) > maxUsage {
			return true
		}
	}

	return false
}
