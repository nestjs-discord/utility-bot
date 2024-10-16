package automod

func (a *Antispam) AddUserToDeniedList(userId UserId) {
	a.sync.Lock()
	defer a.sync.Unlock()

	a.deniedList.SetWithTTL(string(userId), true, 1, a.denyTTL)
}

func (a *Antispam) IsUserInDeniedList(userId UserId) bool {
	a.sync.Lock()
	defer a.sync.Unlock()

	_, hit := a.deniedList.Get(string(userId))
	return hit
}
