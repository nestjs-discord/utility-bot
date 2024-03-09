package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/archive"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/common"
	dont_ping_mods "github.com/nestjs-discord/utility-bot/internal/discord/command/dont-ping-mods"
	google_it "github.com/nestjs-discord/utility-bot/internal/discord/command/google-it"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/reference"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/solved"
	"github.com/rs/zerolog/log"
	"strings"
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

type subCommands = map[string]map[string]*config.Command
type normalCommands = map[string]*config.Command

func RegisterApplicationCommands(s *discordgo.Session) {
	normalCommands, subCommands := generateCommandsToRegister()

	commands = append(commands, generateDynamicCommands(normalCommands)...)
	commands = append(commands, generateDynamicSubcommands(subCommands)...)

	_, err := s.ApplicationCommandBulkOverwrite(config.GetAppID(), config.GetGuildID(), commands)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bulk overwrite application commands")
		return
	}

	log.Info().Int("len", len(commands)).Msg("registered slash commands")
}

func generateDynamicCommands(normalCommands map[string]*config.Command) (commands []*discordgo.ApplicationCommand) {
	for k, v := range normalCommands {
		permission := calculateCommandPermission(v)

		cmd := &discordgo.ApplicationCommand{
			Name:                     k,
			Description:              v.Description,
			DefaultMemberPermissions: &permission,
			Options:                  defaultOptions,
		}

		commands = append(commands, cmd)
	}

	return
}

func generateDynamicSubcommands(subCommands subCommands) (commands []*discordgo.ApplicationCommand) {
	for k, v := range subCommands {
		var permission int64 = 0
		var options []*discordgo.ApplicationCommandOption

		for s, sd := range v {
			if permission == 0 {
				permission = calculateCommandPermission(sd)
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
			DefaultMemberPermissions: &permission,
		}

		commands = append(commands, cmd)
	}

	return
}

func generateCommandsToRegister() (normalCommands, subCommands) {
	subCommands := subCommands{}
	normalCommands := normalCommands{}

	for cmdName, cmdData := range config.GetYaml().Commands {
		if !strings.Contains(cmdName, " ") {
			normalCommands[cmdName] = cmdData
			continue
		}

		parts := strings.Split(cmdName, " ")
		root := parts[0]
		subCmd := parts[1]

		if subCommands[root] == nil {
			subCommands[root] = make(map[string]*config.Command, 0)
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
func calculateCommandPermission(cmdData *config.Command) int64 {
	if cmdData.Protected {
		return config.BotProtectedContentPermission
	}

	return config.BotDefaultContentPermission
}
