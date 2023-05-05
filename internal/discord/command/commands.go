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
	RegisteredCommands      []*discordgo.ApplicationCommand
	registeredCommandsSlice []string
)

func IsCommandRegistered(commandName ...string) bool {
	for _, c := range registeredCommandsSlice {
		if c == strings.Join(commandName, " ") {
			return true
		}
	}
	return false
}

func GenerateRegisteredCommandsSlice() {
	for _, cmd := range RegisteredCommands {
		registeredCommandsSlice = append(registeredCommandsSlice, cmd.Name)

		if len(cmd.Options) == 0 {
			continue
		}

		for _, opt := range cmd.Options {
			if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
				registeredCommandsSlice = append(registeredCommandsSlice, cmd.Name+" "+opt.Name)
			}
		}
	}

	if len(registeredCommandsSlice) != 0 {
		log.Debug().Str("list", strings.Join(registeredCommandsSlice, "|")).Msg("discord registered commands")
	}
}

func RegisterStaticCommands(s *discordgo.Session) {
	for _, cmd := range StaticCommands {

		if IsCommandRegistered(cmd.Name) {
			continue
		}

		c, err := s.ApplicationCommandCreate(config.GetAppID(), config.GetGuildID(), cmd)
		if err != nil {
			log.Error().Err(err).Str("name", cmd.Name).Msg("failed to create static slash command")
		}

		RegisteredCommands = append(RegisteredCommands, c)
		log.Debug().Str("name", c.Name).Msg("registered static slash command")
	}
}

func RegisterContentCommands(s *discordgo.Session) {
	for cmdName, cmdData := range config.GetConfig().Commands {
		cmdNameParts := strings.Split(cmdName, " ")
		if IsCommandRegistered(cmdNameParts...) {
			continue
		}
		permission := calculateCommandPermission(cmdData)

		if len(cmdNameParts) != 1 {
			// cmdName has spaces, so it's a sub-command
			cmd := &discordgo.ApplicationCommand{
				Name:        cmdNameParts[0],
				Description: cmdData.Description,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        cmdNameParts[1],
						Description: cmdData.Description,
					},
				},
				DefaultMemberPermissions: &permission,
			}

			c, err := s.ApplicationCommandCreate(config.GetAppID(), config.GetGuildID(), cmd)
			if err != nil {
				log.Fatal().Err(err).Str("name", cmdName).Msg("failed to create content sub-command")
				return
			}

			RegisteredCommands = append(RegisteredCommands, c)
			log.Debug().Str("name", cmdName).Msg("registered content sub-command")
			continue
		}

		// cmdNameParts len is equal to 1, which means the cmdName is not a sub-command
		c, err := s.ApplicationCommandCreate(config.GetAppID(), config.GetGuildID(), &discordgo.ApplicationCommand{
			Name:                     cmdName,
			Description:              cmdData.Description,
			DefaultMemberPermissions: &permission,
		})
		if err != nil {
			log.Fatal().Err(err).Str("name", cmdName).Msg("failed to create content slash command")
			return
		}

		RegisteredCommands = append(RegisteredCommands, c)
		log.Debug().Str("name", c.Name).Msg("registered content slash command")
	}
}

func calculateCommandPermission(cmdData *config.Command) int64 {
	if cmdData.Protected {
		return config.ProtectedContentPermission
	}

	return config.DefaultContentPermission
}
