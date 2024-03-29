package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/common"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
)

func ContentHandler(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	name, options := normalizeInteractionData(i)

	// Resolve cached content by name
	cmd, cmdExist := config.GetYaml().Commands[name]
	if !cmdExist {
		return false
	}

	var flags discordgo.MessageFlags

	// Copy the content into a new variable to avoid pointer overwrite.
	content := cmd.Content

	for _, opt := range options {
		if opt.Name == common.OptionHide && opt.Value == true {
			flags = discordgo.MessageFlagsEphemeral
		} else if opt.Name == common.OptionTarget && opt.Value != "" {
			content = fmt.Sprintf("*Suggestion for <@%v>:*\n", opt.Value) + content
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: convertButtonsToMessageComponents(cmd.Buttons),
			Flags:      flags,
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

// normalizeInteractionData normalizes the interaction data received from a Discord interaction create event.
// It extracts the name and options from the interaction data, accounting for sub-commands if present.
func normalizeInteractionData(i *discordgo.InteractionCreate) (string, []*discordgo.ApplicationCommandInteractionDataOption) {
	name := i.ApplicationCommandData().Name
	options := i.ApplicationCommandData().Options

	// Overwrite the "name" and "options" variables if the incoming event is a type sub-command
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
			name += " " + opt.Name
			options = opt.Options
			break
		}
	}

	return name, options
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
