package yaml

import "errors"

type PrivilegedForChannel map[string][]string

func (p PrivilegedForChannel) IsUserRolesPrivilegedInChannel(userRoleIDs []string, channelID string) bool {
	roles, ok := p[channelID]
	if !ok {
		return false
	}

	for _, role := range roles {
		for _, roleIDs := range userRoleIDs {
			if role == roleIDs {
				return true
			}
		}
	}
	return false
}

func (p PrivilegedForChannel) validate() error {
	if len(p) == 0 {
		return errors.New("at least 1 privileged role is required")
	}

	return nil
}

func NewPrivilegedForChannel(config *Config) PrivilegedForChannel {
	return config.PrivilegedForChannel
}
