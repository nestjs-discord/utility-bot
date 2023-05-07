package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/npm"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/stats"
	"github.com/rs/zerolog/log"
	"strings"
)

var (
	StaticCommands = []*discordgo.ApplicationCommand{
		npm.Subcommand,
		stats.Command,
	}
	// RegisteredCommands stores both static and dynamic commands
	// that can be easily configured after the bot is launched
	RegisteredCommands []*discordgo.ApplicationCommand
)

type subCommands = map[string]map[string]*config.Command
type normalCommands = map[string]*config.Command

func RegisterStaticCommands(s *discordgo.Session) {
	for _, cmd := range StaticCommands {

		c, err := s.ApplicationCommandCreate(config.GetAppID(), config.GetGuildID(), cmd)
		if err != nil {
			log.Error().Err(err).Str("name", cmd.Name).Msg("failed to create static slash command")
		}

		RegisteredCommands = append(RegisteredCommands, c)
		log.Info().Str("name", c.Name).Msg("registered static command")
	}
}

func RegisterDynamicCommands(s *discordgo.Session) {
	subCommands, normalCommands := generateCommandsToRegister()
	registerDynamicSubcommands(s, subCommands)
	registerDynamicCommands(s, normalCommands)
}

func registerDynamicCommands(s *discordgo.Session, normalCommands map[string]*config.Command) {
	for k, v := range normalCommands {
		permission := calculateCommandPermission(v)
		cmd := &discordgo.ApplicationCommand{
			Name:                     k,
			Description:              v.Description,
			DefaultMemberPermissions: &permission,
		}

		c, err := s.ApplicationCommandCreate(config.GetAppID(), config.GetGuildID(), cmd)
		if err != nil {
			log.Fatal().Err(err).Str("name", k).Msg("failed to create content slash command")
			return
		}

		RegisteredCommands = append(RegisteredCommands, c)
		log.Info().Str("name", c.Name).Msg("registered dynamic command")
	}
}

func registerDynamicSubcommands(s *discordgo.Session, subCommands subCommands) {
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
			})
		}

		cmd := &discordgo.ApplicationCommand{
			Name:                     k,
			Description:              "sub-commands related to " + k,
			Options:                  options,
			DefaultMemberPermissions: &permission,
		}

		c, err := s.ApplicationCommandCreate(config.GetAppID(), config.GetGuildID(), cmd)
		if err != nil {
			log.Fatal().Err(err).Str("name", k).Msg("failed to create content sub-command")
			return
		}
		RegisteredCommands = append(RegisteredCommands, c)
		log.Info().Str("name", k).Msg("registered dynamic sub-command")
	}
}

func generateCommandsToRegister() (subCommands, normalCommands) {
	subCommands := subCommands{}
	normalCommands := normalCommands{}

	for cmdName, cmdData := range config.GetConfig().Commands {
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

	return subCommands, normalCommands
}

// calculateCommandPermission returns the appropriate content permission level for a given command.
// If the command is marked as protected, the function returns the ProtectedContentPermission constant.
// Otherwise, the function returns the DefaultContentPermission constant.
//
// Parameters:
// - cmdData: a pointer to a config.Command object representing the command to calculate permission for.
//
// Returns:
// - An int64 representing the calculated content permission level.
func calculateCommandPermission(cmdData *config.Command) int64 {
	if cmdData.Protected {
		return config.ProtectedContentPermission
	}

	return config.DefaultContentPermission
}
