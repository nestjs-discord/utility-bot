package config

import (
	"github.com/rs/zerolog/log"
	"os"
)

var c config

func GetConfig() *config {
	return &c
}

func GetAppID() string {
	k := os.Getenv("DISCORD_APP_ID")
	if k == "" {
		log.Fatal().Msg("config: DISCORD_APP_ID environment variable is required")
	}
	return k
}

func GetBotToken() string {
	k := os.Getenv("DISCORD_BOT_TOKEN")
	if k == "" {
		log.Fatal().Msg("config: DISCORD_BOT_TOKEN environment variable is required")
	}
	return k
}

func GetGuildID() string {
	k := os.Getenv("DISCORD_GUILD_ID")
	if k == "" {
		log.Fatal().Msg("config: DISCORD_GUILD_ID environment variable is required")
	}
	return k
}

type config struct {
	Moderators []string            `mapstructure:"moderators" validate:"required,min=1,dive,min=1"`
	Ratelimit  Ratelimit           `mapstructure:"ratelimit" validate:"required,dive"`
	Commands   map[string]*Command `mapstructure:"commands" validate:"required,min=1,max=85,dive"`
}

type Ratelimit struct {
	TTL     int    `mapstructure:"ttl" validate:"required,min=1"`
	Usage   int    `mapstructure:"usage" validate:"required,min=2"`
	Message string `mapstructure:"message" validate:"required,min=3"`
}

type Command struct {
	Description string `mapstructure:"description" validate:"required,min=1,max=100"`
	Content     string `mapstructure:"content" validate:"required,min=1"`
	Protected   bool   `mapstructure:"protected" validate:"boolean"`
}
