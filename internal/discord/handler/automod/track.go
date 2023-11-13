package automod

import "github.com/bwmarrin/discordgo"

func (a *AutoMod) SetChannels(channels []*discordgo.Channel) {
	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		a.trackedChannelsIds = append(a.trackedChannelsIds, ChannelId(channel.ID))
	}
}

func (a *AutoMod) IsChannelIdTrackable(channelId ChannelId) bool {
	for _, trackedID := range a.trackedChannelsIds {
		if trackedID == channelId {
			return true
		}
	}
	return false
}

func (a *AutoMod) GetTrackedChannelIds() []ChannelId {
	return a.trackedChannelsIds
}
