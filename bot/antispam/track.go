package antispam

func (a *Antispam) IsChannelIdTrackable(channelId string) bool {
	for _, cid := range a.opts.Cfg.TrackedChannelIds {
		if cid == channelId {
			return true
		}
	}
	return false
}

func (a *Antispam) GetTrackedChannelIds() []string {
	return a.opts.Cfg.TrackedChannelIds
}
