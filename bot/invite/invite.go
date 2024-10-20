package invite

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/bot/permissions"
)

func Link(appId string) string {
	scope := "bot+applications.commands"
	return fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&permissions=%d&scope=%s",
		appId,
		permissions.InviteLink,
		scope,
	)
}
