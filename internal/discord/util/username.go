package util

import "github.com/bwmarrin/discordgo"

func FormatUsername(user *discordgo.User) string {
	if user.Discriminator == "0" {
		return user.Username
	}

	return user.Username + "#" + user.Discriminator
}
