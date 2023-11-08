package discord

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/nestjs-discord/utility-bot/internal/config"
)

func GenerateInviteLink() string {
	qs := &struct {
		ClientID    string `url:"client_id"`
		Permissions int    `url:"permissions"`
		Scope       string `url:"scope"`
	}{
		ClientID:    config.GetAppID(),
		Permissions: config.BotPermissions,
		Scope:       "bot applications.commands",
	}

	v, err := query.Values(qs)
	if err != nil {
		panic(fmt.Sprintf("failed to generate querystring: %v", err))
	}

	return "https://discord.com/api/oauth2/authorize?" + v.Encode()
}
