package commands

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

type Commands struct {
	logger     *slog.Logger
	discordCfg *env.DiscordConfig
	session    *dgo.Session
}

func NewCommands(
	session *dgo.Session,
	discordCfg *env.DiscordConfig,
	commands yaml.Commands,
	archive *archive.Archive,
	credits *credits.Credits,
	dontPingMods *dont_ping_mods.DontPingMods,
	googleIt *google_it.GoogleIt,
	reference *reference.Reference,
	solved *solved.Solved,
) (*Commands, error) {
	c := &Commands{
		logger:     logger.NewWithSubsystem("bot", "commands"),
		discordCfg: discordCfg,
		session:    session,
	}

	staticCommands := []*dgo.ApplicationCommand{
		archive.Command(),
		dontPingMods.Command(),
		credits.Command(),
		googleIt.Command(),
		reference.Command(),
		solved.Command(),
	}

	err := c.registerApplicationCommands(staticCommands, commands)
	if err != nil {
		return nil, err
	}

	return c, nil
}

var (
	defaultOptions = []*dgo.ApplicationCommandOption{ // TODO: can this be removed?
		common.TargetOption,
		common.HideOption,
	}
)

// TODO: refactor this file

type subCommandsType = map[string]yaml.Commands

func (c *Commands) registerApplicationCommands(staticCommands []*dgo.ApplicationCommand, cfgCommands yaml.Commands) error { // TODO: make this a private method
	normalCmd, subCmd := c.generateCommandsToRegister(cfgCommands)

	commands := staticCommands // TODO: refactor

	commands = append(commands, c.generateDynamicCommands(normalCmd)...)
	commands = append(commands, c.generateDynamicSubcommands(subCmd)...)

	_, err := c.session.ApplicationCommandBulkOverwrite(
		c.discordCfg.AppId,
		c.discordCfg.GuildId, // TODO: can we globally register the commands instead? (on production only)
		commands,
	)
	if err != nil {
		return fmt.Errorf("failed to bulk overwrite application commands: %s", err)
	}

	c.logger.Info("registered application commands",
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
			Options:                  defaultOptions,
		}

		commands = append(commands, cmd)
	}

	return
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
				Options:     defaultOptions,
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
