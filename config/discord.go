package config

import "github.com/bwmarrin/discordgo"

// BotPermissions represents the permissions required by the bot.
const BotPermissions int = discordgo.PermissionViewChannel |
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

// BotIntents represents the intents the bot is interested in.
const BotIntents = discordgo.IntentsGuildMessages |
	discordgo.IntentsMessageContent

// BotProtectedContentPermission represents the permissions required for protected content.
const BotProtectedContentPermission int64 = discordgo.PermissionManageMessages |
	discordgo.PermissionUseSlashCommands

// BotDefaultContentPermission represents the default permissions required for dynamic slash-commands.
const BotDefaultContentPermission int64 = discordgo.PermissionUseSlashCommands
