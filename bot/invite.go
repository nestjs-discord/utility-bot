package bot

import (
	"fmt"
	"github.com/google/go-querystring/query"
)

type inviteLinkQuery struct {
	ClientID    string `url:"client_id"`
	Permissions int    `url:"permissions"`
	Scope       string `url:"scope"`
}

func (b *Bot) InviteLink() (string, error) {
	qs := inviteLinkQuery{
		ClientID:    b.cfg.AppId,
		Permissions: permission,
		Scope:       "bot applications.commands",
	}

	v, err := query.Values(qs)
	if err != nil {
		return "", fmt.Errorf("failed to generate querystring: %v", err)
	}

	link := "https://discord.com/api/oauth2/authorize?" + v.Encode()
	return link, nil
}
