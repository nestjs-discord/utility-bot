package config

import (
	"os"
)

var (
	discordAppId    = "DISCORD_APP_ID"
	discordBotToken = "DISCORD_BOT_TOKEN"
	discordGuildId  = "DISCORD_GUILD_ID"
	// List of environment variables to perform validation using the validateEnvVars function.
	toValidate = []string{
		discordAppId,
		discordBotToken,
		discordGuildId,
	}
)

func GetAppID() string {
	return os.Getenv(discordAppId)
}

func GetBotToken() string {
	return os.Getenv(discordBotToken)
}

func GetGuildID() string {
	return os.Getenv(discordGuildId)
}
