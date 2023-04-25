package handler

import (
	"fmt"
	"github.com/erosdesire/discord-nestjs-utility-bot/core/config"
	"github.com/erosdesire/discord-nestjs-utility-bot/internal/discord/command"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Info().
		Str("id", s.State.User.ID).
		Str("username", fmt.Sprintf("%s#%s", m.User.Username, m.User.Discriminator)).
		Msg("logged in as")

	c, err := s.ApplicationCommands(m.User.ID, config.GetGuildID())
	if err != nil {
		log.Panic().Err(err).Msg("failed to fetch registered application commands")
		return
	}
	command.RegisteredCommands = append(command.RegisteredCommands, c...)

	registerStaticSlashCommands(s)
	registerContentSlashCommands(s)

	if err := updateStatus(s); err != nil {
		log.Panic().Err(err).Msg("failed to update status")
		return
	}
	log.Debug().Msg("status updated")

	log.Info().Msg("ready")
}

func isCommandRegistered(commandName string) bool {
	for _, cmd := range command.RegisteredCommands {
		if cmd.Name == commandName {
			return true
		}
	}
	return false
}

func registerStaticSlashCommands(s *discordgo.Session) {
	for _, cmd := range command.StaticCommands {

		if isCommandRegistered(cmd.Name) {
			continue
		}

		c, err := s.ApplicationCommandCreate(s.State.User.ID, config.GetGuildID(), cmd)
		if err != nil {
			log.Error().Err(err).Str("name", cmd.Name).Msg("failed to create static slash command")
		}

		command.RegisteredCommands = append(command.RegisteredCommands, c)
		log.Debug().Str("name", c.Name).Msg("registered static slash command")
	}
}

func registerContentSlashCommands(s *discordgo.Session) {
	for cmd, cmdData := range config.GetConfig().Commands {
		if isCommandRegistered(cmd) {
			continue
		}

		c, err := s.ApplicationCommandCreate(s.State.User.ID, config.GetGuildID(), &discordgo.ApplicationCommand{
			Name:        cmd,
			Description: cmdData.Description,
		})
		if err != nil {
			log.Error().Err(err).Str("name", cmd).Msg("failed to create content slash command")
		}

		command.RegisteredCommands = append(command.RegisteredCommands, c)
		log.Debug().Str("name", c.Name).Msg("registered content slash command")
	}
}

func updateStatus(s *discordgo.Session) error {
	return s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "slash commands",
				Type: discordgo.ActivityTypeListening,
			},
		},
	})
}
