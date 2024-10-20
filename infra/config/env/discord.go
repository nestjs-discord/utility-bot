package env

import (
	"errors"
	"os"
)

type DiscordConfig struct {
	Token   string
	AppId   string
	GuildId GuildId
}

func NewDiscordConfig() (*DiscordConfig, error) {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		return nil, errors.New("DISCORD_BOT_TOKEN environment variable not set")
	}

	appId := os.Getenv("DISCORD_APP_ID")
	if appId == "" {
		return nil, errors.New("DISCORD_APP_ID environment variable not set")
	}

	guildId := os.Getenv("DISCORD_GUILD_ID")
	if guildId == "" {
		return nil, errors.New("DISCORD_GUILD_ID environment variable not set")
	}

	return &DiscordConfig{
		Token:   token,
		AppId:   appId,
		GuildId: GuildId(guildId),
	}, nil
}

type GuildId string

func (id GuildId) String() string {
	return string(id)
}

func ProvideGuildId(c *DiscordConfig) GuildId {
	return c.GuildId
}
