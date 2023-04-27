package discord

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type inviteQuerystring struct {
	ClientID    string `url:"client_id"`
	Permissions int    `url:"permissions"`
	Scope       string `url:"scope"`
}

var Invite = &cobra.Command{
	Use:   "discord:invite",
	Short: "Generates an invite link to add the bot to servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		qs := &inviteQuerystring{
			ClientID:    config.GetAppID(),
			Permissions: config.BotPermissions,
			Scope:       "bot applications.commands",
		}
		v, err := query.Values(qs)
		if err != nil {
			return errors.Wrap(err, "failed to generate querystring")
		}

		fmt.Print("https://discord.com/api/oauth2/authorize?" + v.Encode())

		return nil
	},
}
