package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/cache"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/npm"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/stats"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/interaction"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	userID := i.Member.User.ID

	log.Debug().Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", userID).
		Interface("options", i.ApplicationCommandData().Options).
		Msg("event: interaction create")

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
	// Npm subcommand
	case npm.Name:
		for _, option := range i.ApplicationCommandData().Options {

			switch option.Name {
			case npm.SearchCommandName:
				interaction.NpmSearchHandler(s, i)
				return
			case npm.InspectCommandName:
				interaction.NpmInspectHandler(s, i)
				return
			}
		}
	case stats.Stats:
		interaction.StatHandler(s, i)
		return
	}

	interaction.DefaultHandler(s, i)
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
