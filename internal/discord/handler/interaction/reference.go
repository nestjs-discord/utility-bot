package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/algolia"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/common"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/reference"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
)

func ReferenceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	for _, option := range options {
		app, ok := algolia.Apps[option.Name]
		if !ok {
			continue
		}

		flags, content := parseReferenceOptions(option)

		objectID, err := getStringValueByName(reference.Query, option.Options)
		if err != nil {
			util.InteractionRespondError(err, s, i)
			return
		}

		hit, err := algolia.GetObject(app, objectID)
		if err != nil {
			util.InteractionRespondError(err, s, i)
			return
		}

		content += "**" + algolia.ResolveHitToName(*hit) + "**\n"
		if hit.Content != "" {
			content += algolia.Truncate(hit.Content, 300, "") + "\n"
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    content,
				Flags:      flags,
				Components: generateReferenceComponents(hit),
			},
		})
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
	}
}

func parseReferenceOptions(option *discordgo.ApplicationCommandInteractionDataOption) (discordgo.MessageFlags, string) {
	var flags discordgo.MessageFlags
	var content string

	for _, opt := range option.Options {
		if opt.Name == common.OptionHide && opt.Value == true {
			flags = discordgo.MessageFlagsEphemeral
		} else if opt.Name == common.OptionTarget && opt.Value != "" {
			content += fmt.Sprintf("Suggestion for <@%v>:\n", opt.Value)
		}
	}

	return flags, content
}

func generateReferenceComponents(hit *algolia.Hit) []discordgo.MessageComponent {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Read more",
					URL:   hit.URL,
					Style: discordgo.LinkButton,
					Emoji: discordgo.ComponentEmoji{Name: "📖"},
				},
			},
		},
	}
	return components
}

func ReferenceAutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, option := range i.ApplicationCommandData().Options {
		app, ok := algolia.Apps[option.Name]
		if !ok {
			continue
		}

		var choices []*discordgo.ApplicationCommandOptionChoice

		query, err := getStringValueByName(reference.Query, option.Options)
		if err != nil {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})
			return
		}

		hits, err := algolia.Search(app, query)

		if err == nil {
			for _, hit := range hits {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  algolia.Truncate(algolia.ResolveHitToName(hit), 95, ""),
					Value: hit.ObjectID,
				})
			}
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	}
}

func getStringValueByName(name string, options []*discordgo.ApplicationCommandInteractionDataOption) (string, error) {
	for _, opt := range options {
		if opt.Name == name {
			return opt.StringValue(), nil
		}
	}

	return "", fmt.Errorf("couldn't get string value of %v", name)
}
