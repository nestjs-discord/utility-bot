package antispam

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"strings"
)

func (a *Antispam) GenerateAlertMessage(i *dgo.MessageCreate) *dgo.MessageSend {
	return &dgo.MessageSend{
		Content:    "",
		Embed:      a.generateAlertEmbed(i),
		Components: a.generateAlertComponents(i),
		Files:      a.generateAlertFiles(i),
	}
}

func (a *Antispam) generateAlertFiles(i *dgo.MessageCreate) []*dgo.File {
	var files []*dgo.File
	userUniqueMessages := a.GetUserUniqueMessages(userIdType(i.Author.ID))
	for msg, msgId := range userUniqueMessages {
		fileName := fmt.Sprintf("msg-%s.txt", msgId)
		files = append(files, &dgo.File{
			Name:        fileName,
			ContentType: "text/plain",
			Reader:      strings.NewReader(msg),
		})
	}
	return files
}

func (a *Antispam) generateAlertComponents(i *dgo.MessageCreate) []dgo.MessageComponent {
	return []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				dgo.Button{
					Emoji: &dgo.ComponentEmoji{
						Name: "🔗",
					},
					Label: "Jump to the message",
					Style: dgo.LinkButton,
					URL:   fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.ID),
				},
			},
		},
	}
}

func (a *Antispam) generateAlertEmbed(i *dgo.MessageCreate) *dgo.MessageEmbed {
	embed := &dgo.MessageEmbed{
		Title:       "Spam alert! 🚨",
		Color:       0xff0000, // Red
		Description: a.generateAlertEmbedDescription(),
	}

	embed.Fields = append(embed.Fields, &dgo.MessageEmbedField{
		Name:   "Username",
		Value:  "`" + i.Author.String() + "`",
		Inline: true,
	})

	authorAccCreatedAt, err := dgo.SnowflakeTimestamp(i.Author.ID)
	if err == nil {
		embed.Fields = append(embed.Fields, &dgo.MessageEmbedField{
			Name: "Account created",
			//Value: humanize.Time(authorAccCreatedAt),
			Value:  fmt.Sprintf("<t:%d:R>", authorAccCreatedAt.UTC().Unix()),
			Inline: true,
		})
	}

	embed.Fields = append(embed.Fields, &dgo.MessageEmbedField{
		Name:  "Search query",
		Value: "`from: " + i.Author.ID + "`",
	})

	userUniqueMessages := a.GetUserUniqueMessages(userIdType(i.Author.ID))
	for msg := range userUniqueMessages {
		// Sanitize userMsg to avoid breaking the code block
		msg = strings.ReplaceAll(msg, "```", "")

		// If sanitizedMsg is longer than 350 characters, truncate it and add three dots
		if len(msg) > 350 {
			msg = msg[:347] + "..."
		}

		embed.Fields = append(embed.Fields, &dgo.MessageEmbedField{
			Value: "```text\n" + msg + "\n```",
		})
	}

	return embed
}

func (a *Antispam) generateAlertEmbedDescription() string {
	return fmt.Sprintf(
		"Member exceeded channel limit `%d` within `%d` seconds."+"\n"+
			"Added to the denied list for the next `%d` seconds.",
		a.opts.Cfg.MaxChannelsPerUser,
		a.opts.Cfg.MessageTTLSec,
		a.opts.Cfg.DenyTTLSec,
	)
}
