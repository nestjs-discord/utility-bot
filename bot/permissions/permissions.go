package permissions

import (
	dgo "github.com/bwmarrin/discordgo"
)

const (
	InviteLink = dgo.PermissionViewChannel |
		dgo.PermissionSendMessages |
		dgo.PermissionSendMessagesInThreads |
		dgo.PermissionAttachFiles |
		dgo.PermissionEmbedLinks |
		dgo.PermissionUseApplicationCommands |
		dgo.PermissionAddReactions |
		dgo.PermissionUseExternalEmojis |
		dgo.PermissionManageThreads |
		dgo.PermissionManageMessages |
		dgo.PermissionReadMessageHistory |
		dgo.PermissionKickMembers |
		dgo.PermissionBanMembers
	DefaultCommands   = dgo.PermissionUseApplicationCommands
	ProtectedCommands = dgo.PermissionManageMessages | dgo.PermissionUseApplicationCommands
	BotIntents = dgo.IntentsGuildMessages | dgo.IntentsMessageContent | dgo.IntentAutoModerationExecution
)
