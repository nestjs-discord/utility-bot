package moderator

import "encoding/base64"

type Moderator struct {
	userIds []string
}

func NewModerator(encUserIds []string) (*Moderator, error) {
	m := &Moderator{}
	for index := range encUserIds {
		decStr, err := m.decodeUserId(encUserIds[index])
		if err != nil {
			return nil, err
		}

		m.userIds = append(m.userIds, decStr)
	}

	return m, nil
}

func (mod *Moderator) decodeUserId(encStr string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encStr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (mod *Moderator) IsUserModerator(userId string) bool {
	for _, uid := range mod.userIds {
		if uid == userId {
			return true
		}
	}
	return false
}
