package util

import "github.com/nestjs-discord/utility-bot/config"

func IsUserModerator(userId string) bool {
	for _, id := range config.Yaml().Moderators {
		if id == userId {
			return true
		}
	}

	return false
}
