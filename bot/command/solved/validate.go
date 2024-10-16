package solved

import (
	dgo "github.com/bwmarrin/discordgo"
)

func (c *Solved) validateChannelType(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
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

func (c *Solved) validateThreadLock(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
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

func (c *Solved) validateChannelOwner(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	postOwnerId := channel.OwnerID
	executedById := i.Member.User.ID
	if postOwnerId == executedById ||
		c.moderators.IsUserModerator(executedById) {
		return true
	}

	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "⚠️ Only forum post owner and moderators can use this command.",
			Flags:   dgo.MessageFlagsEphemeral,
		},
	})

	return false
}
