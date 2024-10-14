package bot

import "github.com/bwmarrin/discordgo"

const (
	permission = discordgo.PermissionViewChannel |
		discordgo.PermissionSendMessages |
		discordgo.PermissionSendMessagesInThreads |
		discordgo.PermissionAttachFiles |
		discordgo.PermissionEmbedLinks |
		discordgo.PermissionUseSlashCommands |
		discordgo.PermissionAddReactions |
		discordgo.PermissionUseExternalEmojis |
		discordgo.PermissionManageThreads |
		discordgo.PermissionManageMessages |
		discordgo.PermissionReadMessageHistory |
		discordgo.PermissionKickMembers |
		discordgo.PermissionBanMembers
	permissionDefault   = discordgo.PermissionUseSlashCommands
	permissionProtected = discordgo.PermissionManageMessages | discordgo.PermissionUseSlashCommands
	intents             = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
)
