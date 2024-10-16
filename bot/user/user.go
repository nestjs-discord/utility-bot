package user

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
)

func FormatUsername(user *dgo.User) string {
	if user.Discriminator == "0" {
		return user.Username
	}

	return user.Username + "#" + user.Discriminator
}

func Mention(user *dgo.User) string {
	return fmt.Sprintf("<@%s>", user.ID)
	// return fmt.Sprintf("<@!%s>", user.ID) // TODO: is this silent mention?
}
