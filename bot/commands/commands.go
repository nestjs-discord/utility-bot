package commands

// TODO: refactor this file

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/archive"
	"github.com/nestjs-discord/utility-bot/bot/commands/common"
	"github.com/nestjs-discord/utility-bot/bot/commands/credits"
	"github.com/nestjs-discord/utility-bot/bot/commands/dont_ping_mods"
	"github.com/nestjs-discord/utility-bot/bot/commands/google_it"
	"github.com/nestjs-discord/utility-bot/bot/commands/reference"
	"github.com/nestjs-discord/utility-bot/bot/commands/solved"
	"github.com/nestjs-discord/utility-bot/bot/permissions"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
	"strings"
)

type Options struct {
	Commands     yaml.Commands
	Session      *dgo.Session
	DiscordCfg   *env.DiscordConfig
	Archive      *archive.Archive
	Credits      *credits.Credits
	DontPingMods *dont_ping_mods.DontPingMods
	GoogleIt     *google_it.GoogleIt
	Reference    *reference.Reference
	Solved       *solved.Solved
}

type Commands struct {
	logger *slog.Logger
	opts   Options
}

func NewCommands(opts Options) (*Commands, error) {
	c := &Commands{
		logger: logger.NewWithSubsystem("bot", "commands"),
		opts:   opts,
	}

	staticCommands := []*dgo.ApplicationCommand{
		opts.Archive.Command(),
		opts.DontPingMods.Command(),
		opts.Credits.Command(),
		opts.GoogleIt.Command(),
		opts.Reference.Command(),
		opts.Solved.Command(),
	}

	err := c.registerApplicationCommands(staticCommands)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type subCommandsType = map[string]yaml.Commands

func (c *Commands) registerApplicationCommands(staticCommands []*dgo.ApplicationCommand) error {
	normalCmd, subCmd := c.generateCommandsToRegister(c.opts.Commands)

	commands := staticCommands // TODO: refactor

	commands = append(commands, c.generateDynamicCommands(normalCmd)...)
	commands = append(commands, c.generateDynamicSubcommands(subCmd)...)

	//guildId := c.opts.DiscordCfg.GuildId.String()

	// Clear any registered guild commands (for backward compatibility)
	//_, err := c.opts.Session.ApplicationCommandBulkOverwrite(
	//	c.opts.DiscordCfg.AppId,
	//	guildId,
	//	make([]*dgo.ApplicationCommand, 0),
	//)
	//if err != nil {
	//	return fmt.Errorf("failed to bulk overwrite guild commands: %s", err)
	//}

	// Register global commands (we don't want to use guild commands anymore)
	guildId := ""
	_, err := c.opts.Session.ApplicationCommandBulkOverwrite(c.opts.DiscordCfg.AppId, guildId, commands)
	if err != nil {
		return fmt.Errorf("failed to bulk overwrite global commands: %s", err)
	}

	c.logger.Info("registered global commands",
		slog.Int("len", len(commands)),
	)

	return nil
}

func (c *Commands) generateDynamicCommands(normalCommands yaml.Commands) (commands []*dgo.ApplicationCommand) {
	for k, v := range normalCommands {
		perm := c.calculateCommandPermission(v)

		cmd := &dgo.ApplicationCommand{
			Name:                     k,
			Description:              v.Description,
			DefaultMemberPermissions: &perm,
			Options:                  c.getDefaultOptions(),
		}

		commands = append(commands, cmd)
	}

	return
}

func (c *Commands) getDefaultOptions() []*dgo.ApplicationCommandOption {
	return []*dgo.ApplicationCommandOption{
		common.TargetOption,
		common.HideOption,
	}
}

func (c *Commands) generateDynamicSubcommands(subCommands subCommandsType) (commands []*dgo.ApplicationCommand) {
	for k, v := range subCommands {
		var perm int64 = 0
		var options []*dgo.ApplicationCommandOption

		for s, sd := range v {
			if perm == 0 {
				perm = c.calculateCommandPermission(sd)
			}

			options = append(options, &dgo.ApplicationCommandOption{
				Type:        dgo.ApplicationCommandOptionSubCommand,
				Name:        s,
				Description: sd.Description,
				Options:     c.getDefaultOptions(),
			})
		}

		cmd := &dgo.ApplicationCommand{
			Name:                     k,
			Description:              "sub-commands related to " + k,
			Options:                  options,
			DefaultMemberPermissions: &perm,
		}

		commands = append(commands, cmd)
	}

	return
}

func (c *Commands) generateCommandsToRegister(cfgCommands yaml.Commands) (yaml.Commands, subCommandsType) {
	subCommands := subCommandsType{}
	normalCommands := yaml.Commands{}

	for cmdName, cmdData := range cfgCommands {
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
func (c *Commands) calculateCommandPermission(cmdData yaml.Command) int64 {
	if cmdData.Protected {
		return permissions.ProtectedCommands
	}

	return permissions.DefaultCommands
}

func CleanApplicationCommands(session *dgo.Session, appId string, guildId string) error {
	emptyCmd := make([]*dgo.ApplicationCommand, 0)
	_, err := session.ApplicationCommandBulkOverwrite(
		appId,
		guildId,
		emptyCmd,
	)
	if err != nil {
		return fmt.Errorf("unable to clean the app commands: %v", err)
	}

	return nil
}
