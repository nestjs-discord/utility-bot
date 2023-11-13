package automod

func (a *AutoMod) AddUserToDeniedList(userId UserId) {
	a.sync.Lock()
	defer a.sync.Unlock()

	a.deniedList.SetWithTTL(string(userId), true, 1, a.denyTTL)
}

func (a *AutoMod) IsUserInDeniedList(userId UserId) bool {
	a.sync.Lock()
	defer a.sync.Unlock()

	_, hit := a.deniedList.Get(string(userId))
	return hit
}
