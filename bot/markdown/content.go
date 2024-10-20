package markdown

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/common"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
)

func (m *Markdown) ContentHandler(s *dgo.Session, i *dgo.InteractionCreate) error {
	name, options := m.normalizeInteractionData(i)

	cmd, cmdExist := m.opts.Commands[name]
	if !cmdExist {
		return nil // skip
	}

	var flags dgo.MessageFlags

	// Copy the content into a new variable to avoid pointer overwrite.
	content := cmd.Content

	dynamicComponents := m.convertButtonsToMessageComponents(cmd.Buttons)

	for _, opt := range options {
		if opt.Name == common.OptionHide && opt.Value == true {
			flags = dgo.MessageFlagsEphemeral
		} else if opt.Name == common.OptionTarget && opt.Value != "" {
			userIdToMention := opt.UserValue(s)
			if userIdToMention.Bot {
				respond.InteractionWithEphemeralMessage(s, i, "You are not allowed to mention bots.")
				return nil
			}
			content = fmt.Sprintf("*Suggestion for:* %s\n\n", userIdToMention.Mention()) + content

			if len(dynamicComponents) > 3 { // Discord limit
				continue
			}

			ackButton, err := m.generateAckButton(userIdToMention)
			if err != nil {
				return fmt.Errorf("generate ack button failed: %v", err)
			}

			dynamicComponents = append(dynamicComponents, dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					ackButton,
				},
			})
		}
	}

	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content:    content,
			Components: dynamicComponents,
			Flags:      flags,
		},
	})
	if err == nil {
		return nil
	}

	return fmt.Errorf("content handler failed: %s", err)
}

// normalizeInteractionData normalizes the interaction data received from a Discord interaction create event.
// It extracts the name and options from the interaction data, accounting for sub-commands if present.
func (m *Markdown) normalizeInteractionData(i *dgo.InteractionCreate) (string, []*dgo.ApplicationCommandInteractionDataOption) {
	name := i.ApplicationCommandData().Name
	options := i.ApplicationCommandData().Options

	// Overwrite the "name" and "options" variables if the incoming event is a type sub-command
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Type == dgo.ApplicationCommandOptionSubCommand {
			name += " " + opt.Name
			options = opt.Options
			break
		}
	}

	return name, options
}

func (m *Markdown) convertButtonsToMessageComponents(buttons yaml.CommandButtons) []dgo.MessageComponent {
	var components []dgo.MessageComponent
	for _, row := range buttons {
		componentsInRow := make([]dgo.MessageComponent, 0, len(row))
		for _, btn := range row {
			componentsInRow = append(componentsInRow, dgo.Button{
				Label: btn.Label,
				URL:   btn.URL,
				Style: dgo.LinkButton,
				Emoji: &dgo.ComponentEmoji{Name: btn.Emoji},
			})
		}
		components = append(components, dgo.ActionsRow{
			Components: componentsInRow,
		})
	}

	return components
}
