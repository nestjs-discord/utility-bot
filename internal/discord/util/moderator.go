package util

import (
	"github.com/nestjs-discord/utility-bot/infra/config"
)

func IsUserModerator(userId string) bool { // TODO: remove this
	for _, id := range config.Yaml().Moderators {
		if id == userId {
			return true
		}
	}

	return false
}
