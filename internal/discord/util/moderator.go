package util

import "github.com/erosdesire/discord-nestjs-utility-bot/core/config"

func IsUserModerator(userId string) bool {
	for _, id := range config.GetConfig().Moderators {
		if id == userId {
			return true
		}
	}

	return false
}
