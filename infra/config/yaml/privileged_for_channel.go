package yaml

import "errors"

type PrivilegedForChannel map[string][]string

func (p PrivilegedForChannel) IsUserPrivilegedInChannel(userID, channelID string) bool {
	if users, ok := p[channelID]; ok {
		for _, user := range users {
			if user == userID {
				return true
			}
		}
	}
	return false
}

func (p PrivilegedForChannel) validate() error {
	if len(p) == 0 {
		return errors.New("at least 1 privileged user is required")
	}

	return nil
}

func NewPrivilegedForChannel(config *Config) PrivilegedForChannel {
	return config.PrivilegedForChannel
}
