package archive

import dgo "github.com/bwmarrin/discordgo"

func (a *Archive) validateChannelType(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	if channel.Type == dgo.ChannelTypeGuildPublicThread ||
		channel.Type == dgo.ChannelTypeGuildPrivateThread {
		return true
	}

	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "⚠️ You can only use this command in forum posts.",
			Flags:   dgo.MessageFlagsEphemeral,
		},
	})

	return false
}

func (a *Archive) validateThreadLock(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	if !channel.ThreadMetadata.Locked {
		return true
	}

	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "⚠️ Cannot perform this action on the locked forum posts",
			Flags:   dgo.MessageFlagsEphemeral,
		},
	})

	return false
}
