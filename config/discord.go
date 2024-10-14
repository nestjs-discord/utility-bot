package config

import "github.com/bwmarrin/discordgo"

// TODO: remove this file after migrating to "bot/permission"

// BotProtectedContentPermission represents the permissions required for protected content.
const BotProtectedContentPermission int64 = discordgo.PermissionManageMessages |
	discordgo.PermissionUseSlashCommands

// BotDefaultContentPermission represents the default permissions required for dynamic slash-commands.
const BotDefaultContentPermission int64 = discordgo.PermissionUseSlashCommands
