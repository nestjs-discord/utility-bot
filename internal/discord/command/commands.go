package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/npm"
	"github.com/rs/zerolog/log"
)

var (
	StaticCommands = []*discordgo.ApplicationCommand{
		npm.Subcommand,
	}
	// RegisteredCommands stores both static and dynamic commands
	// that can be easily configured after the bot is launched
	RegisteredCommands []*discordgo.ApplicationCommand
)

func IsCommandRegistered(commandName string) bool {
	for _, cmd := range RegisteredCommands {
		if cmd.Name == commandName {
			return true
		}
	}
	return false
}

func RegisterStaticSlashCommands(s *discordgo.Session) {
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

func RegisterContentSlashCommands(s *discordgo.Session) {
	for cmd, cmdData := range config.GetConfig().Commands {
		if IsCommandRegistered(cmd) {
			continue
		}

		c, err := s.ApplicationCommandCreate(config.GetAppID(), config.GetGuildID(), &discordgo.ApplicationCommand{
			Name:        cmd,
			Description: cmdData.Description,
		})
		if err != nil {
			log.Error().Err(err).Str("name", cmd).Msg("failed to create content slash command")
		}

		RegisteredCommands = append(RegisteredCommands, c)
		log.Debug().Str("name", c.Name).Msg("registered content slash command")
	}
}
