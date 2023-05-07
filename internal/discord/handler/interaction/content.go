package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
)

func ContentHandler(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	name := i.ApplicationCommandData().Name
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
			name += " " + opt.Name
		}
	}

	// Check cache if the tag exists
	cmd, cmdExist := config.GetConfig().Commands[name]
	if !cmdExist {
		return false
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    cmd.Content,
			Components: convertButtonsToMessageComponents(cmd.Buttons),
		},
	})
	if err == nil {
		return true
	}

	log.Error().
		Err(err).
		Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", i.Member.User.ID).
		Msg("failed to respond to interaction")

	util.InteractionRespondError(err, s, i)

	return false
}

func convertButtonsToMessageComponents(b [][]*config.Button) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent
	for _, row := range b {
		componentsInRow := make([]discordgo.MessageComponent, 0, len(row))
		for _, btn := range row {
			componentsInRow = append(componentsInRow, discordgo.Button{
				Label: btn.Label,
				URL:   btn.URL,
				Style: discordgo.LinkButton,
				Emoji: discordgo.ComponentEmoji{Name: btn.Emoji},
			})
		}
		components = append(components, discordgo.ActionsRow{
			Components: componentsInRow,
		})
	}

	return components
}
