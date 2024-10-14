package config

import (
	"os"
)

// TODO: remove this file

func GetAppID() string {
	return os.Getenv("DISCORD_APP_ID")
}

func GetBotToken() string {
	return os.Getenv("DISCORD_BOT_TOKEN")
}

func GetGuildID() string {
	return os.Getenv("DISCORD_GUILD_ID")
}
