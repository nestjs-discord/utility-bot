package reference

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/common"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/nestjs-discord/utility-bot/pkg/algolia"
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

func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			util.InteractionRespondError(err, s, i)
			return
		}

		hit, err := algolia.GetObject(app, objectID)
		if err != nil {
			util.InteractionRespondError(err, s, i)
			return
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

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    content.String(),
				Flags:      flags | discordgo.MessageFlagsSupressEmbeds,
				Components: generateReferenceComponents(hit),
			},
		})
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
	}
}

func parseReferenceOptions(option *discordgo.ApplicationCommandInteractionDataOption, content *strings.Builder) discordgo.MessageFlags {
	var flags discordgo.MessageFlags

	for _, opt := range option.Options {
		if opt.Name == common.OptionHide && opt.Value == true {
			flags = discordgo.MessageFlagsEphemeral
		} else if opt.Name == common.OptionTarget && opt.Value != "" {
			content.WriteString(fmt.Sprintf("*Suggestion for <@%v>:*\n", opt.Value))
		}
	}

	return flags
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

func getStringValueByName(name string, options []*discordgo.ApplicationCommandInteractionDataOption) (string, error) {
	for _, opt := range options {
		if opt.Name == name {
			return opt.StringValue(), nil
		}
	}

	return "", fmt.Errorf("couldn't get string value of %v", name)
}
