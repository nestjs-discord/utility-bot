package antispam

func (a *Antispam) IsChannelIdTrackable(channelId string) bool {
	for _, cid := range a.cfg.TrackedChannelIds {
		if cid == channelId {
			return true
		}
	}
	return false
}

func (a *Antispam) GetTrackedChannelIds() []string {
	return a.cfg.TrackedChannelIds
}
