package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/message"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
)

func MessageCreate(s *discordgo.Session, i *discordgo.MessageCreate) {
	// Check if the message is sent by a bot; if true, skip further processing.
	if i.Message.Author.Bot {
		return
	}

	log.Debug().Str("id", i.ID).Str("content", i.Content).Msg("event: message create")

	// Check if auto-mod is enabled in the configuration.
	if config.GetConfig().AutoMod.Enabled {
		message.AutoModHandler(s, i)
	}

	if !util.IsUserModerator(i.Author.ID) {
		return
	}

	switch i.Content {
	case "!stats":
		message.StatsHandler(s, i)
		return
	case "!automod":
		message.AutoModTrackHandler(s, i)
		return
	}
}
