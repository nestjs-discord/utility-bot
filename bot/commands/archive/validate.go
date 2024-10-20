package archive

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
)

func (a *Archive) validateChannelType(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	if channel.Type == dgo.ChannelTypeGuildPublicThread ||
		channel.Type == dgo.ChannelTypeGuildPrivateThread {
		return true
	}

	respond.InteractionWithEphemeralMessage(s, i, "⚠️ You can only use this command in forum posts.")

	return false
}

func (a *Archive) validateThreadLock(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	if !channel.ThreadMetadata.Locked {
		return true
	}

	respond.InteractionWithEphemeralMessage(s, i, "⚠️ Cannot perform this action on the locked forum posts.")

	return false
}
