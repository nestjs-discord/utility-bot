package util

import "github.com/nestjs-discord/utility-bot/internal/config"

func IsUserModerator(userId string) bool {
	for _, id := range config.GetConfig().Moderators {
		if id == userId {
			return true
		}
	}

	return false
}
