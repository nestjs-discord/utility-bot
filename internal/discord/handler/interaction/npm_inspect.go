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
	options := &npmAPI.InspectOptions{}
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

	var mainComponents []discordgo.MessageComponent
	if data.Homepage != "" {
		mainComponents = append(mainComponents, discordgo.Button{
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
			mainComponents = append(mainComponents, discordgo.Button{
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
		mainComponents = append(mainComponents, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "🐞",
			},
			Label: "Bugs",
			Style: discordgo.LinkButton,
			URL:   data.Bugs.Url,
		})
	}

	if data.Dist.Tarball != "" {
		mainComponents = append(mainComponents, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "🔗",
			},
			Label: "Download",
			Style: discordgo.LinkButton,
			URL:   data.Dist.Tarball,
		})
	}

	var messageComponents []discordgo.MessageComponent

	if len(mainComponents) != 0 {
		messageComponents = append(messageComponents, discordgo.ActionsRow{
			Components: mainComponents,
		})
	}

	if data.Funding.Url != "" {
		messageComponents = append(messageComponents, discordgo.ActionsRow{
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

func generateFields(data *npmAPI.InspectResponse) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField

	if data.Version != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Version",
			Value:  data.Version,
			Inline: true,
		})
	}

	if data.License != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "License",
			Value:  strings.TrimSpace(data.License),
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

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "",
		Value: "",
	})

	if len(data.Dependencies) != 0 {
		v := ""
		for item, value := range data.Dependencies {
			v += fmt.Sprintf("`%v %v`\n", item, value)
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Dependencies",
			Value:  v,
			Inline: true,
		})
	}

	if len(data.DevDependencies) != 0 {
		v := ""
		for item, value := range data.DevDependencies {
			v += fmt.Sprintf("`%v %v`\n", item, value)
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Dev Dependencies",
			Value:  v,
			Inline: true,
		})
	}

	if len(data.PeerDependencies) != 0 {
		v := ""
		for item, value := range data.PeerDependencies {
			v += fmt.Sprintf("`%v %v`\n", item, value)
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Peer Dependencies",
			Value: v,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "",
		Value: "",
	})

	if data.Dist.FileCount != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Files count (folder excluded)",
			Value:  humanize.Comma(data.Dist.FileCount),
			Inline: true,
		})
	}

	if data.Dist.UnpackedSize != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Unpacked size",
			Value:  humanize.Bytes(data.Dist.UnpackedSize),
			Inline: true,
		})
	}

	if data.Dist.Integrity != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Integrity",
			Value: "`" + data.Dist.Integrity + "`",
		})
	}

	if data.GitHead != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Git head",
			Value: "`" + data.GitHead + "`",
		})
	}

	if len(data.Keywords) != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Keywords",
			Value: strings.Join(data.Keywords, ", "),
		})
	}

	return fields
}
