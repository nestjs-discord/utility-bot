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
		dgo.PermissionUseSlashCommands |
		dgo.PermissionAddReactions |
		dgo.PermissionUseExternalEmojis |
		dgo.PermissionManageThreads |
		dgo.PermissionManageMessages |
		dgo.PermissionReadMessageHistory |
		dgo.PermissionKickMembers |
		dgo.PermissionBanMembers
	DefaultCommands   = dgo.PermissionUseSlashCommands
	ProtectedCommands = dgo.PermissionManageMessages |
		dgo.PermissionUseSlashCommands
	BotIntents = dgo.IntentsGuildMessages | dgo.IntentsMessageContent | dgo.IntentAutoModerationExecution
)
