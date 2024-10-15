package automod

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/config"
	"log/slog"
)

type Message struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"createdAt"` // Unix timestamp
	Content   string `json:"content"`
}

func NewMessage(ID string, content string) (Message, error) {
	msgTimestamp, err := discordgo.SnowflakeTimestamp(ID)
	if err != nil {
		return Message{}, fmt.Errorf("failed to get snowflake timestamp: %s", err)
	}

	return Message{
		ID:        ID,
		CreatedAt: msgTimestamp.UTC().Unix(),
		Content:   content,
	}, nil
}

func (a *AutoMod) StoreMessage(userId UserId, channelId string, message Message) {
	a.sync.Lock()
	defer a.sync.Unlock()

	if _, ok := a.userMap[userId]; !ok {
		a.userMap[userId] = make(map[string]Message)
	}

	a.userMap[userId][channelId] = message
}

func (a *AutoMod) GetUserUniqueMessages(userId UserId) map[string]string {
	a.sync.Lock()
	defer a.sync.Unlock()

	uniqueMap := map[string]string{}
	for _, message := range a.userMap[userId] {
		uniqueMap[message.Content] = message.ID
	}

	return uniqueMap
}

// GetUserMessages retrieves the messages associated with a user and organizes them in a map.
// The keys of the map represent the channel IDs, and the corresponding values are the message IDs.
//
// Note: This function is designed to be used with an AutoMod instance and requires a valid UserId parameter.
//
// Parameters:
//   - userId: The unique identifier of the user for whom messages are to be retrieved.
//
// Returns:
//   - map[string]string: A map where keys are channel IDs, and values are message IDs.
func (a *AutoMod) GetUserMessages(userId UserId) map[string]string {
	a.sync.Lock()
	defer a.sync.Unlock()

	res := map[string]string{}

	for channelId, message := range a.userMap[userId] {
		res[string(channelId)] = message.ID
	}

	return res
}

func (a *AutoMod) Handler(s *discordgo.Session, i *discordgo.MessageCreate) {
	channelId := i.ChannelID

	// Skip executing auto-mod logic if the provided channel ID is not in the list of channels being tracked.
	// This check ensures that auto-mod actions are only applied to channels marked for moderation.
	if !a.IsChannelIdTrackable(channelId) {
		a.logger.Debug("auto mod: channel id is not trackable, skipping...",
			slog.String("channelId", channelId),
		)
		return
	}

	// Check if the author is a moderator; if true, skip further processing.
	if a.cfg.ModeratorsBypass && a.moderator.IsUserModerator(i.Author.ID) {
		return
	}

	userId := UserId(i.Author.ID)

	if a.IsUserInDeniedList(userId) {

		// Delete their message
		_ = s.ChannelMessageDelete(i.ChannelID, i.ID)

		// Try to ban them again
		_ = s.GuildBanCreateWithReason(i.GuildID, i.Author.ID, "The antispam feature flagged this user!", 7)

		return
	}

	message, err := NewMessage(i.ID, i.Content)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to create message: %s", err))
		return
	}

	// Store the user message in the AutoMod cache.
	a.StoreMessage(userId, channelId, message)

	// Check if the user has sent messages to an excessive number of channels within the defined maximum channels limit.
	// If true, further processing is skipped.
	if a.IsUserWithinMaxChannelsLimit(userId) {
		return
	}

	// Delete their previous messages in another go routine
	go func() {
		userMessages := a.GetUserMessages(userId)
		for chId, msgId := range userMessages {
			err = s.ChannelMessageDelete(chId, msgId)
			if err != nil {
				a.logger.Error(fmt.Sprintf("failed to delete message: %s", err),
					slog.String("channelId", chId),
					slog.String("messageId", msgId),
				)
				return
			}
			a.logger.Debug(fmt.Sprintf("deleted message"),
				slog.String("channelId", chId),
				slog.String("messageId", msgId),
			)
		}
	}()

	// Add user to the denied list
	a.AddUserToDeniedList(userId)

	logChannelId := config.Yaml().AutoMod.LogChannelId
	_, err = s.ChannelMessageSendComplex(logChannelId, a.GenerateAlertMessage(i))
	if err != nil {
		a.logger.Error("failed to alert moderators about the ongoing spam",
			slog.Any("err", err),
		)
	}

	// Ban their account
	err = s.GuildBanCreateWithReason(i.GuildID, i.Author.ID, "spam", 7)
	if err != nil {
		a.logger.Error("failed to ban a spammer",
			slog.String("userId", i.Author.ID),
			slog.Any("err", err),
		)
		_, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf(":hammer: Failed to ban the spammer: `%s`", err.Error()))
		return
	}

	a.logger.Info("banned a user",
		slog.String("userId", i.Author.ID),
	)

	_, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf(":hammer: Member banned: `%s`", i.Author.ID))

	// for debugging purposes only
	// jsonStr, _ := json.MarshalIndent(cache.AutoMod, "", "  ")
	// _, _ = s.ChannelMessageSend(logChannelId, fmt.Sprintf("```json\n%s\n```", string(jsonStr)))
}

func (a *AutoMod) TrackHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
	content := "### Antispam feature is tracking the following channels: 👇\n"
	for _, channelId := range a.cfg.ChannelIds {
		content += fmt.Sprintf("- <#%s>\n", channelId)
	}

	_, _ = s.ChannelMessageSend(i.ChannelID, content)
}
