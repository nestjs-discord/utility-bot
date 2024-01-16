package config

import (
	"github.com/rs/zerolog/log"
	"os"
)

var c Config

func GetConfig() *Config {
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

type Config struct {
	Moderators []string            `mapstructure:"moderators" validate:"required,min=1,dive,min=1"`
	Ratelimit  Ratelimit           `mapstructure:"ratelimit" validate:"required"`
	AutoMod    AutoMod             `mapstructure:"autoMod" validate:"required"`
	Commands   map[string]*Command `mapstructure:"commands" validate:"required,max-one-space-allowed,min=1,max=85,dive"`
}

type Ratelimit struct {
	TTL     int    `mapstructure:"ttl" validate:"required,min=1"`
	Usage   int    `mapstructure:"usage" validate:"required,min=2"`
	Message string `mapstructure:"message" validate:"required,min=3"`
}

type AutoMod struct {
	Enabled                 bool   `mapstructure:"enabled" validate:"boolean"`
	ModeratorsBypass        bool   `mapstructure:"moderatorsBypass" validate:"boolean"`
	LogChannelId            string `mapstructure:"logChannelId" validate:"required,min=1"`
	LogMentionRoleId        string `mapstructure:"logMentionRoleId"`
	MessageTTL              int    `mapstructure:"messageTTL" validate:"required,min=1"`
	MaxChannelsLimitPerUser int    `mapstructure:"maxChannelsLimitPerUser" validate:"required,min=1"`
	DenyTTL                 int    `mapstructure:"denyTTL" validate:"required,min=1"`
}

type Command struct {
	Description string      `mapstructure:"description" validate:"required,min=1,max=100"`
	Content     string      `mapstructure:"content" validate:"required,min=1"`
	Protected   bool        `mapstructure:"protected" validate:"boolean"`
	Buttons     [][]*Button `mapstructure:"buttons" validate:"min=0,max=8,dive,min=1,max=4,dive"`
}

type Button struct {
	Label string `mapstructure:"label" validate:"required,min=3,max=40"`
	URL   string `mapstructure:"url" validate:"required,url,min=3"`
	Emoji string `mapstructure:"emoji" validate:"regexp=^[\p{Emoji}]$"`
}
