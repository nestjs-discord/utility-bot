package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/cache"
	"github.com/nestjs-discord/utility-bot/internal/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/automod"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
)

func AutoModHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
	channelId := automod.ChannelId(i.ChannelID)

	// Skip executing auto-mod logic if the provided channel ID is not in the list of channels being tracked.
	// This check ensures that auto-mod actions are only applied to channels marked for moderation.
	if !cache.AutoMod.IsChannelIdTrackable(channelId) {
		log.Debug().Interface("channel-id", channelId).Msg("auto mod: channel id is not trackable, skipping...")
		return
	}

	// Check if the author is a moderator; if true, skip further processing.
	if config.GetConfig().AutoMod.ModeratorsBypass && util.IsUserModerator(i.Author.ID) {
		return
	}

	userId := automod.UserId(i.Author.ID)

	if cache.AutoMod.IsUserInDeniedList(userId) {
		// _ = s.ChannelMessageDelete(i.ChannelID, i.ID)
		return
	}

	message, err := automod.NewMessage(i.ID, i.Content)
	if err != nil {
		log.Err(err).Msg("auto mod: failed to init new message")
		return
	}

	// Store the user message in the AutoMod cache.
	cache.AutoMod.StoreMessage(userId, channelId, message)

	// Check if the user has sent messages to an excessive number of channels within the defined maximum channels limit.
	// If true, further processing is skipped.
	if cache.AutoMod.IsUserWithinMaxChannelsLimit(userId) {
		return
	}

	// Add user to the denied list
	cache.AutoMod.AddUserToDeniedList(userId)

	logChannelId := config.GetConfig().AutoMod.LogChannelId
	_, err = s.ChannelMessageSendComplex(logChannelId, cache.AutoMod.GenerateAlertMessage(i))
	if err != nil {
		log.Err(err).Msg("auto mod: failed to notify log channel about the ongoing spam")
	}

	// for debugging purposes only
	// jsonStr, _ := json.MarshalIndent(cache.AutoMod, "", "  ")
	// _, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf("```json\n%s\n```", string(jsonStr)))
}
