package automod

func (a *AutoMod) IsChannelIdTrackable(channelId string) bool {
	for _, cid := range a.cfg.ChannelIds {
		if cid == channelId {
			return true
		}
	}
	return false
}

func (a *AutoMod) GetTrackedChannelIds() []string {
	return a.cfg.ChannelIds
}
