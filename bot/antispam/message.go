package antispam

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"log/slog"
)

type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channelId"`
	CreatedAt int64  `json:"createdAt"` // Unix timestamp
	Content   string `json:"content"`
}

func NewMessage(messageCreateEvent *dgo.MessageCreate) (Message, error) {
	msgTimestamp, err := dgo.SnowflakeTimestamp(messageCreateEvent.ID)
	if err != nil {
		return Message{}, fmt.Errorf("failed to get snowflake timestamp: %s", err)
	}

	return Message{
		ID:        messageCreateEvent.ID,
		ChannelID: messageCreateEvent.ChannelID,
		CreatedAt: msgTimestamp.UTC().Unix(),
		Content:   messageCreateEvent.Content,
	}, nil
}

func (a *Antispam) StoreMessage(userId userIdType, message Message) {
	a.sync.Lock()
	defer a.sync.Unlock()

	a.userToMessagesMap[userId] = append(a.userToMessagesMap[userId], message)
}

func (a *Antispam) GetUserUniqueMessages(userId userIdType) map[string]string {
	a.sync.Lock()
	defer a.sync.Unlock()

	uniqueMap := map[string]string{}
	for _, message := range a.userToMessagesMap[userId] {
		uniqueMap[message.Content] = message.ID
	}

	return uniqueMap
}

// GetUserMessages retrieves the messages associated with a user and organizes them in a map.
// The keys of the map represent the channel IDs, and the corresponding values are the message IDs.
//
// Note: This function is designed to be used with an Antispam instance and requires a valid userIdType parameter.
//
// Parameters:
//   - userId: The unique identifier of the user for whom messages are to be retrieved.
//
// Returns:
//   - map[string]string: A map where keys are channel IDs, and values are message IDs.
func (a *Antispam) GetUserMessages(userId userIdType) []Message {
	a.sync.Lock()
	defer a.sync.Unlock()

	return a.userToMessagesMap[userId]
}

func (a *Antispam) Handler(s *dgo.Session, i *dgo.MessageCreate) {
	channelId := i.ChannelID

	// Skip executing auto-mod logic if the provided channel ID is not in the list of channels being tracked.
	// This check ensures that auto-mod actions are only applied to channels marked for moderation.
	if !a.IsChannelIdTrackable(channelId) {
		a.logger.Debug("channel id is not trackable",
			slog.String("channelId", channelId),
		)
		return
	}

	// Check if the author is a moderator; if true, skip further processing.
	if a.opts.Cfg.ModeratorsBypass && a.opts.Moderators.IsUserModerator(i.Author.ID) {
		return
	}

	userId := userIdType(i.Author.ID)

	if a.IsUserInDeniedList(userId) {

		// Delete their message
		_ = s.ChannelMessageDelete(i.ChannelID, i.ID)

		// Try to ban them again
		// _ = s.GuildBanCreateWithReason(i.GuildID, i.Author.ID, "The antispam feature flagged this user!", 7)

		return
	}

	message, err := NewMessage(i)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to create message: %s", err))
		return
	}

	// Cache the message.
	a.StoreMessage(userId, message)

	// Check if the user has sent messages to an excessive number of channels within the defined maximum channels limit.
	// If true, further processing is skipped.
	if !a.IsUserWithinMaxChannelsLimit(userId) {
		// Delete their previous messages in another go routine
		go func() {
			userMessages := a.GetUserMessages(userId)
			for _, usrMsg := range userMessages {

				err = s.ChannelMessageDelete(usrMsg.ChannelID, usrMsg.ID)
				if err != nil {
					a.logger.Error(fmt.Sprintf("failed to delete message: %s", err),
						slog.String("channelId", usrMsg.ChannelID),
						slog.String("messageId", usrMsg.ID),
					)
					return
				}
				a.logger.Debug("deleted message",
					slog.String("channelId", usrMsg.ChannelID),
					slog.String("messageId", usrMsg.ID),
				)
			}
		}()

		// Add user to the denied list
		a.AddUserToDeniedList(userId)

		logChannelId := a.opts.Cfg.LogChannelId
		_, err = s.ChannelMessageSendComplex(logChannelId, a.GenerateAlertMessage(i))
		if err != nil {
			a.logger.Error("failed to alert moderators about the ongoing spam",
				slog.Any("err", err),
			)
		}

		// Ban their account
		err = s.GuildBanCreateWithReason(i.GuildID, i.Author.ID, "spam", 7)
		if err != nil {
			a.logger.Error("failed to ban the spammer",
				slog.String("userId", i.Author.ID),
				slog.Any("err", err),
			)
			_, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf("🔨 Failed to ban the spammer: `%s`", err.Error()))
			return
		}

		a.logger.Info("banned a user",
			slog.String("userId", i.Author.ID),
		)

		_, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf("🔨 Member banned: `%s`", i.Author.ID))

		// for debugging purposes only
		// jsonStr, _ := json.MarshalIndent(cache.Antispam, "", "  ")
		// _, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf("```json\n%s\n```", string(jsonStr)))
	}

	if repeatedMessages := a.GetUserRepeatedMessages(userId); len(repeatedMessages) > 0 {
		// Add user to the denied list
		a.AddUserToDeniedList(userId)

		a.logger.Info("repeated messages", slog.Any("messages", repeatedMessages))

		logChannelId := a.opts.Cfg.LogChannelId
		alertMsg := a.GenerateRepeatedMessagesFoundAlertMessage(i, repeatedMessages)
		_, _ = s.ChannelMessageSendComplex(logChannelId, alertMsg)
	}
}

func (a *Antispam) TrackHandler(s *dgo.Session, i *dgo.MessageCreate) {
	content := "### Antispam feature is tracking the following channels: 👇\n"
	for _, channelId := range a.opts.Cfg.TrackedChannelIds {
		content += fmt.Sprintf("- <#%s>\n", channelId)
	}

	_, _ = s.ChannelMessageSend(i.ChannelID, content)
}
