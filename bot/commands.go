package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/command/archive"
	"github.com/nestjs-discord/utility-bot/bot/command/common"
	"github.com/nestjs-discord/utility-bot/bot/command/dont_ping_mods"
	"github.com/nestjs-discord/utility-bot/bot/command/google_it"
	"github.com/nestjs-discord/utility-bot/bot/command/reference"
	"github.com/nestjs-discord/utility-bot/bot/command/solved"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/config/yaml"
)

var (
	commands = []*discordgo.ApplicationCommand{
		archive.Command,
		reference.Subcommand,
		solved.Command,
		google_it.Command,
		dont_ping_mods.Command,
	}
	defaultOptions = []*discordgo.ApplicationCommandOption{
		common.TargetOption,
		common.HideOption,
	}
)

// TODO: refactor this file

type subCommandsType = map[string]yaml.Commands

func (b *Bot) RegisterApplicationCommands() error {
	normalCmd, subCmd := generateCommandsToRegister()

	commands = append(commands, generateDynamicCommands(normalCmd)...)
	commands = append(commands, generateDynamicSubcommands(subCmd)...)

	_, err := b.session.ApplicationCommandBulkOverwrite(
		b.discordCfg.AppId,
		b.discordCfg.GuildId, // TODO: can we globally register the commands instead? (on production only)
		commands,
	)
	if err != nil {
		return fmt.Errorf("failed to bulk overwrite application commands: %s", err)
	}

	b.logger.Info("registered application commands",
		slog.Int("len", len(commands)),
	)
	return nil
}

func generateDynamicCommands(normalCommands yaml.Commands) (commands []*discordgo.ApplicationCommand) {
	for k, v := range normalCommands {
		perm := calculateCommandPermission(v)

		cmd := &discordgo.ApplicationCommand{
			Name:                     k,
			Description:              v.Description,
			DefaultMemberPermissions: &perm,
			Options:                  defaultOptions,
		}

		commands = append(commands, cmd)
	}

	return
}

func generateDynamicSubcommands(subCommands subCommandsType) (commands []*discordgo.ApplicationCommand) {
	for k, v := range subCommands {
		var perm int64 = 0
		var options []*discordgo.ApplicationCommandOption

		for s, sd := range v {
			if perm == 0 {
				perm = calculateCommandPermission(sd)
			}

			options = append(options, &discordgo.ApplicationCommandOption{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        s,
				Description: sd.Description,
				Options:     defaultOptions,
			})
		}

		cmd := &discordgo.ApplicationCommand{
			Name:                     k,
			Description:              "sub-commands related to " + k,
			Options:                  options,
			DefaultMemberPermissions: &perm,
		}

		commands = append(commands, cmd)
	}

	return
}

func generateCommandsToRegister() (yaml.Commands, subCommandsType) {
	subCommands := subCommandsType{}
	normalCommands := yaml.Commands{}

	for cmdName, cmdData := range config.Yaml().Commands {
		if !strings.Contains(cmdName, " ") {
			normalCommands[cmdName] = cmdData
			continue
		}

		parts := strings.Split(cmdName, " ")
		root := parts[0]
		subCmd := parts[1]

		if subCommands[root] == nil {
			subCommands[root] = make(yaml.Commands)
		}

		subCommands[root][subCmd] = cmdData
	}

	return normalCommands, subCommands
}

// calculateCommandPermission returns the appropriate content permission level for a given command.
// If the command is marked as protected, the function returns the BotProtectedContentPermission constant.
// Otherwise, the function returns the BotDefaultContentPermission constant.
//
// Parameters:
// - cmdData: a pointer to a config.Command object representing the command to calculate permission for.
//
// Returns:
// - An int64 representing the calculated content permission level.
func calculateCommandPermission(cmdData yaml.Command) int64 {
	if cmdData.Protected {
		return permissionProtected
	}

	return permissionDefault
}

func (b *Bot) CleanApplicationCommands() error {
	emptyCmd := make([]*discordgo.ApplicationCommand, 0)
	_, err := b.session.ApplicationCommandBulkOverwrite(
		b.discordCfg.AppId,
		b.discordCfg.GuildId,
		emptyCmd,
	)
	if err != nil {
		return fmt.Errorf("unable to clean the app commands: %v", err)
	}

	b.logger.Info("cleaned application commands")

	return nil
}
