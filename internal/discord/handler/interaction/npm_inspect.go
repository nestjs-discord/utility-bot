package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/npm"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	npmAPI "github.com/nestjs-discord/utility-bot/internal/npm"
	"net/url"
	"strings"
)

func NpmInspectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := &npmAPI.InspectOptions{
		Version: "latest",
	}
	for _, parentOption := range i.ApplicationCommandData().Options {
		for _, childOption := range parentOption.Options {
			switch childOption.Name {
			case npm.InspectOptionPackage:
				options.Name = childOption.StringValue()
			case npm.InspectOptionVersion:
				options.Version = childOption.StringValue()
			}
		}
	}

	if err := npmAPI.IsNPMPackageNameValid(options.Name); err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		return
	}

	if err := npmAPI.IsNPMVersionValid(options.Version); err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		return
	}

	data, err := npmAPI.Inspect(options)
	if err != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "I couldn't find any package with the specified detail.\n\nIf you believe this is a mistake, contact moderators.",
			},
		})
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
		return
	}

	fields := generateFields(data)
	messageComponents := generateComponents(data)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeRich,
					Title:       data.Name,
					Fields:      fields,
					URL:         "https://npmjs.com/package/" + options.Name,
					Description: data.Description,
				},
			},
			Components: messageComponents,
		},
	})
	if err != nil {
		util.InteractionRespondError(err, s, i)
	}
}

func generateComponents(data *npmAPI.InspectResponse) (components []discordgo.MessageComponent) {
	firstRowComponents := generateFirstRowComponents(data)

	if len(firstRowComponents) != 0 {
		components = append(components, discordgo.ActionsRow{
			Components: firstRowComponents,
		})
	}

	if data.Funding.Url != "" {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "🤗",
					},
					Label: "Funding: " + data.Funding.Type,
					Style: discordgo.LinkButton,
					URL:   data.Funding.Url,
				},
			},
		})
	}

	return
}

func generateFirstRowComponents(data *npmAPI.InspectResponse) (firstRowComponents []discordgo.MessageComponent) {
	if data.Homepage != "" {
		firstRowComponents = append(firstRowComponents, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "📦",
			},
			Label: "Homepage",
			Style: discordgo.LinkButton,
			URL:   data.Homepage,
		})
	}

	if data.Repository.Type != "" && data.Repository.Url != "" {
		repositoryURL := strings.Replace(data.Repository.Url, data.Repository.Type+"+", "", 1)

		// Validate URL
		_, err := url.Parse(repositoryURL)
		if err == nil {
			firstRowComponents = append(firstRowComponents, discordgo.Button{
				Emoji: discordgo.ComponentEmoji{
					Name: "🔗",
				},
				Label: "Repository",
				Style: discordgo.LinkButton,
				URL:   repositoryURL,
			})
		}
	}

	if data.Bugs.Url != "" {
		firstRowComponents = append(firstRowComponents, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "🐞",
			},
			Label: "Bugs",
			Style: discordgo.LinkButton,
			URL:   data.Bugs.Url,
		})
	}

	if data.Dist.Tarball != "" {
		firstRowComponents = append(firstRowComponents, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "🔗",
			},
			Label: "Download",
			Style: discordgo.LinkButton,
			URL:   data.Dist.Tarball,
		})
	}

	return
}

func generateFields(data *npmAPI.InspectResponse) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField

	fields = generateHeaderFields(data, fields)

	fields = appendSpacer(fields)

	fields = generateDependenciesFields(data, fields)

	fields = appendSpacer(fields)

	if data.Dist.UnpackedSize != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Unpacked size",
			Value:  "`" + humanize.Bytes(data.Dist.UnpackedSize) + "`",
			Inline: true,
		})
	}

	if data.Dist.FileCount != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Files count (folders excluded)",
			Value:  "`" + humanize.Comma(data.Dist.FileCount) + "`",
			Inline: true,
		})
	}

	fields = appendSpacer(fields)

	if data.GitHead != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Git head",
			Value: "`" + data.GitHead + "`",
		})
	}

	if data.Dist.Integrity != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Integrity",
			Value: "`" + data.Dist.Integrity + "`",
		})
	}

	if len(data.Keywords) != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Keywords",
			Value: wrapKeywords(data.Keywords),
		})
	}

	return fields
}

func wrapKeywords(keywords []string) string {
	var wrappedKeywords []string
	for _, keyword := range keywords {
		wrappedKeyword := fmt.Sprintf("`%s`", keyword)
		wrappedKeywords = append(wrappedKeywords, wrappedKeyword)
	}
	return strings.Join(wrappedKeywords, ", ")
}

func generateHeaderFields(data *npmAPI.InspectResponse, fields []*discordgo.MessageEmbedField) []*discordgo.MessageEmbedField {
	if data.Version != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Version",
			Value:  "`" + data.Version + "`",
			Inline: true,
		})
	}

	if data.License != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "License",
			Value:  "`" + data.License + "`",
			Inline: true,
		})
	}

	if len(data.Engines) != 0 {
		v := ""
		for item, value := range data.Engines {
			v += fmt.Sprintf("`%v %v`\n", item, value)
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Engines",
			Value:  v,
			Inline: true,
		})
	}
	return fields
}

func generateDependenciesFields(data *npmAPI.InspectResponse, fields []*discordgo.MessageEmbedField) []*discordgo.MessageEmbedField {
	appendDependencyFields := func(deps map[string]string, name string, inline bool) []*discordgo.MessageEmbedField {
		if len(deps) == 0 {
			return fields
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  generateDependencyFieldsValue(deps),
			Inline: inline,
		})

		return fields
	}

	fields = appendDependencyFields(data.Dependencies, "Dependencies", true)
	fields = appendDependencyFields(data.DevDependencies, "Dev Dependencies", true)
	fields = appendDependencyFields(data.PeerDependencies, "Peer Dependencies", len(data.DevDependencies) == 0)

	return fields
}

func generateDependencyFieldsValue(deps map[string]string) string {
	v := ""
	count := 0
	for item, value := range deps {
		s := fmt.Sprintf("`%v %v`\n", item, value)
		// Discord limits each fields value to maximum 1000 characters
		if len(v)+len(s) > 900 {
			remaining := len(deps) - count
			if remaining > 0 {
				v += fmt.Sprintf("%v more...", remaining)
			}
			break
		}
		v += s
		count++
	}
	return v
}

func appendSpacer(fields []*discordgo.MessageEmbedField) []*discordgo.MessageEmbedField {
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "",
		Value: "",
	})
	return fields
}
