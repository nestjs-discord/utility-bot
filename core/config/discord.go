package config

import "github.com/bwmarrin/discordgo"

const BotPermissions int = discordgo.PermissionViewChannel |
	discordgo.PermissionSendMessages |
	discordgo.PermissionSendMessagesInThreads |
	discordgo.PermissionAttachFiles |
	discordgo.PermissionEmbedLinks |
	discordgo.PermissionUseSlashCommands |
	discordgo.PermissionAddReactions |
	discordgo.PermissionUseExternalEmojis |
	discordgo.PermissionManageThreads

const BotIntents = discordgo.IntentsGuildMessages

const ProtectedContentPermission int64 = discordgo.PermissionManageMessages | discordgo.PermissionUseSlashCommands
const DefaultContentPermission int64 = discordgo.PermissionUseSlashCommands
