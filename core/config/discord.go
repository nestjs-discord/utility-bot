package config

import "github.com/bwmarrin/discordgo"

const BotPermissions int = discordgo.PermissionViewChannel |
	discordgo.PermissionSendMessages |
	discordgo.PermissionSendMessagesInThreads |
	discordgo.PermissionAttachFiles |
	discordgo.PermissionEmbedLinks |
	discordgo.PermissionUseSlashCommands |
	discordgo.PermissionAddReactions

const BotIntents = discordgo.IntentsGuildMessages
