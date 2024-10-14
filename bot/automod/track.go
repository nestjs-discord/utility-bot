package automod

func (a *AutoMod) setChannels(channelIds []string) {
	for _, id := range channelIds {
		a.trackedChannelsIds = append(a.trackedChannelsIds, ChannelId(id))
	}
}

func (a *AutoMod) IsChannelIdTrackable(channelId ChannelId) bool {
	for _, cid := range a.trackedChannelsIds {
		if cid == channelId {
			return true
		}
	}
	return false
}

func (a *AutoMod) GetTrackedChannelIds() []ChannelId {
	return a.trackedChannelsIds
}
