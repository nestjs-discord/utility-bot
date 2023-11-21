package message

import (
	"fmt"
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
		log.Debug().
			Str("channel-id", i.ChannelID).
			Msg("auto mod: channel id is not trackable, skipping...")
		return
	}

	// Check if the author is a moderator; if true, skip further processing.
	if config.GetConfig().AutoMod.ModeratorsBypass && util.IsUserModerator(i.Author.ID) {
		return
	}

	userId := automod.UserId(i.Author.ID)

	if cache.AutoMod.IsUserInDeniedList(userId) {

		// Delete their message
		_ = s.ChannelMessageDelete(i.ChannelID, i.ID)

		// Try to ban them again
		_ = s.GuildBanCreateWithReason(i.GuildID, i.Author.ID, "spam", 7)

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

	// Delete their previous messages in another go routine
	go func() {
		userMessages := cache.AutoMod.GetUserMessages(userId)
		for chId, msgId := range userMessages {
			err = s.ChannelMessageDelete(chId, msgId)
			if err != nil {
				log.Err(err).
					Str("channel-id", chId).
					Str("message-id", msgId).
					Msg("auto mod: failed to delete the message")

				return
			}

			log.Debug().
				Str("channel-id", chId).
				Str("message-id", msgId).
				Msg("auto mod: message delete success")
		}
	}()

	// Add user to the denied list
	cache.AutoMod.AddUserToDeniedList(userId)

	logChannelId := config.GetConfig().AutoMod.LogChannelId
	_, err = s.ChannelMessageSendComplex(logChannelId, cache.AutoMod.GenerateAlertMessage(i))
	if err != nil {
		log.Err(err).Msg("auto mod: failed to notify log channel about the ongoing spam")
	}

	// Ban their account
	err = s.GuildBanCreateWithReason(i.GuildID, i.Author.ID, "spam", 7)
	if err != nil {
		log.Err(err).Str("user-id", i.Author.ID).Msg("auto mod: failed to ban the user")
		_, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf(":hammer: Failed to ban the spammer: `%s`", err.Error()))
		return
	}

	log.Info().Str("user-id", i.Author.ID).Msg("auto mod: banned user")

	_, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf(":hammer: Member banned: `%s`", i.Author.ID))

	// for debugging purposes only
	// jsonStr, _ := json.MarshalIndent(cache.AutoMod, "", "  ")
	// _, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf("```json\n%s\n```", string(jsonStr)))
}
