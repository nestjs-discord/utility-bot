package handler

import (
	"fmt"
	commands "github.com/erosdesire/discord-nestjs-utility-bot/cmd/discord/command"
	"github.com/erosdesire/discord-nestjs-utility-bot/config"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Info().
		Str("id", s.State.User.ID).
		Str("username", fmt.Sprintf("%s#%s", m.User.Username, m.User.Discriminator)).
		Msg("logged in as")

	registerStaticSlashCommands(s)
	registerContentSlashCommands(s)

	if err := updateStatus(s); err != nil {
		log.Error().Err(err).Msg("failed to update status")
	} else {
		log.Debug().Msg("status updated")
	}

	log.Info().Msg("ready")
}

func registerStaticSlashCommands(s *discordgo.Session) {
	for _, cmd := range commands.StaticCommands {
		c, err := s.ApplicationCommandCreate(s.State.User.ID, config.GetConfig().GuildID, cmd)
		if err != nil {
			log.Error().Err(err).Str("name", cmd.Name).Msg("failed to create static slash command")
		}

		commands.RegisteredCommands = append(commands.RegisteredCommands, c)
		log.Debug().Str("name", c.Name).Msg("registered static slash command")
	}
}

func registerContentSlashCommands(s *discordgo.Session) {
	for cmd, cmdData := range config.GetConfig().Commands {
		c, err := s.ApplicationCommandCreate(s.State.User.ID, config.GetConfig().GuildID, &discordgo.ApplicationCommand{
			Name:        cmd,
			Description: cmdData.Description,
		})
		if err != nil {
			log.Error().Err(err).Str("name", cmd).Msg("failed to create content slash command")
		}

		commands.RegisteredCommands = append(commands.RegisteredCommands, c)
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
