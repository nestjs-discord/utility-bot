package moderator

type Moderator struct {
	userIds []string
}

func NewModerator(userIds []string) *Moderator {
	return &Moderator{
		userIds: userIds,
	}
}

func (mod *Moderator) IsUserModerator(userId string) bool {
	for _, uid := range mod.userIds {
		if uid == userId {
			return true
		}
	}
	return false
}
