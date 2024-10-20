package moderators

import (
	"encoding/base64"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
)

type Options struct {
	EncodedUserIds yaml.Moderators
}

type Moderators struct {
	userIds []string
}

func NewModerators(opts Options) (*Moderators, error) {
	m := &Moderators{}
	for index := range opts.EncodedUserIds {
		decStr, err := m.DecodeUserId(opts.EncodedUserIds[index])
		if err != nil {
			return nil, err
		}

		m.userIds = append(m.userIds, decStr)
	}

	return m, nil
}

func (mod *Moderators) DecodeUserId(encStr string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encStr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (mod *Moderators) IsUserModerator(userId string) bool {
	for _, uid := range mod.userIds {
		if uid == userId {
			return true
		}
	}
	return false
}

func (mod *Moderators) UserIds() []string {
	return mod.userIds
}
