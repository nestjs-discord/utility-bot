package reference

import (
	"errors"
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/common"
	"github.com/nestjs-discord/utility-bot/infra/services/algolia"
	"strings"
)

var emojis = map[string]string{
	algolia.Discord.ToSlug():        "<:discord:1106968504877461616>",
	algolia.DiscordJSGuide.ToSlug(): "<:discordjs:1106968508950122637>",
	algolia.Fastify.ToSlug():        "<:fastify:1106968514109116486>",
	algolia.Necord.ToSlug():         "<:necord:1106968169580613723>",
	algolia.NestCommander.ToSlug():  "<:commander:1106968502432190484>",
	algolia.NestJS.ToSlug():         "<:nestjs:1106967607434817698>",
	algolia.Ogma.ToSlug():           "<:ogma:1106968518160814180>",
	algolia.TypeORM.ToSlug():        "<:typeorm:1106976838695264348>",
	algolia.TypeScript.ToSlug():     "<:typescript:1106968521692414043>",
}

func (r *Reference) Handler(s *dgo.Session, i *dgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options

	for _, option := range options {
		app, ok := algolia.Apps[option.Name]
		if !ok {
			continue
		}

		var content strings.Builder

		flags := parseReferenceOptions(option, &content)

		objectID, err := getStringValueByName(QueryOption, option.Options)
		if err != nil {
			return err
		}

		hit, err := algolia.GetObject(app, objectID)
		if err != nil {
			return err
		}

		// Add emoji
		emoji, ok := emojis[option.Name]
		if ok {
			content.WriteString(emoji)
			content.WriteString(" ")
		}

		// Add title
		content.WriteString("**")
		content.WriteString(algolia.GetFormattedHierarchy(*hit))
		content.WriteString("**\n")

		// Add description (if present)
		if hit.Content != "" {
			content.WriteString(algolia.Truncate(hit.Content, 350) + "\n")
		}

		err = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content:    content.String(),
				Flags:      flags | dgo.MessageFlagsSupressEmbeds,
				Components: generateReferenceComponents(hit),
			},
		})
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("reference handler failed")
}

func parseReferenceOptions(option *dgo.ApplicationCommandInteractionDataOption, content *strings.Builder) dgo.MessageFlags {
	var flags dgo.MessageFlags

	for _, opt := range option.Options {
		if opt.Name == common.OptionHide && opt.Value == true {
			flags = dgo.MessageFlagsEphemeral
		} else if opt.Name == common.OptionTarget && opt.Value != "" {
			content.WriteString(fmt.Sprintf("*Suggestion for <@%v>:*\n", opt.Value))
		}
	}

	return flags
}

func generateReferenceComponents(hit *algolia.Hit) []dgo.MessageComponent {
	components := []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				dgo.Button{
					Label: "Read more",
					URL:   hit.URL,
					Style: dgo.LinkButton,
					Emoji: &dgo.ComponentEmoji{Name: "📖"},
				},
			},
		},
	}
	return components
}

func getStringValueByName(name string, options []*dgo.ApplicationCommandInteractionDataOption) (string, error) {
	for _, opt := range options {
		if opt.Name == name {
			return opt.StringValue(), nil
		}
	}

	return "", fmt.Errorf("couldn't get string value of %v", name)
}
