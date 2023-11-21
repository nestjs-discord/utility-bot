package automod

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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

func (a *AutoMod) StoreMessage(userId UserId, channelId ChannelId, message Message) {
	a.sync.Lock()
	defer a.sync.Unlock()

	if _, ok := a.userMap[userId]; !ok {
		a.userMap[userId] = make(map[ChannelId]Message, 0)
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
