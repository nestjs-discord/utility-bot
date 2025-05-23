package solved

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
)

func (c *Solved) validateChannelType(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	if channel.Type == dgo.ChannelTypeGuildPublicThread ||
		channel.Type == dgo.ChannelTypeGuildPrivateThread {
		return true
	}

	respond.InteractionWithEphemeralMessage(s, i, "⚠️ You can only use this command in forum posts")

	return false
}

func (c *Solved) validateThreadLock(s *dgo.Session, i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	if !channel.ThreadMetadata.Locked {
		return true
	}

	respond.InteractionWithEphemeralMessage(s, i,
		"⚠️ Cannot perform this action on the locked forum posts",
	)

	return false
}

func (c *Solved) validateChannelOwner(i *dgo.InteractionCreate, channel *dgo.Channel) bool {
	postOwnerId := channel.OwnerID
	executedById := i.Member.User.ID
	if postOwnerId == executedById ||
		c.opts.Moderators.IsUserModerator(executedById) {
		return true
	}

	return c.opts.CfgPrivileged.IsUserRolesPrivilegedInChannel(i.Member.Roles, channel.ParentID)
}
