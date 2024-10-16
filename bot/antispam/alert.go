package antispam

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/user"
	"strings"
)

func (a *Antispam) GenerateAlertMessage(i *discordgo.MessageCreate) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content:    "",
		Embed:      a.generateAlertEmbed(i),
		Components: a.generateAlertComponents(i),
		Files:      a.generateAlertFiles(i),
	}
}

func (a *Antispam) generateAlertFiles(i *discordgo.MessageCreate) []*discordgo.File {
	var files []*discordgo.File
	userUniqueMessages := a.GetUserUniqueMessages(UserId(i.Author.ID))
	for msg, msgId := range userUniqueMessages {
		fileName := fmt.Sprintf("msg-%s.txt", msgId)
		files = append(files, &discordgo.File{
			Name:        fileName,
			ContentType: "text/plain",
			Reader:      strings.NewReader(msg),
		})
	}
	return files
}

func (a *Antispam) generateAlertComponents(i *discordgo.MessageCreate) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Emoji: &discordgo.ComponentEmoji{
						Name: "🔗",
					},
					Label: "Jump to the message",
					Style: discordgo.LinkButton,
					URL:   fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.ID),
				},
			},
		},
	}
}

func (a *Antispam) generateAlertEmbed(i *discordgo.MessageCreate) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       "Spam alert! 🚨",
		Color:       0xff0000, // Red
		Description: a.generateAlertEmbedDescription(),
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Username",
		Value:  "`" + user.FormatUsername(i.Author) + "`",
		Inline: true,
	})

	authorAccCreatedAt, err := discordgo.SnowflakeTimestamp(i.Author.ID)
	if err == nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Account created",
			//Value: humanize.Time(authorAccCreatedAt),
			Value:  fmt.Sprintf("<t:%d:R>", authorAccCreatedAt.UTC().Unix()),
			Inline: true,
		})
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "Search query",
		Value: "`from: " + i.Author.ID + "`",
	})

	userUniqueMessages := a.GetUserUniqueMessages(UserId(i.Author.ID))
	for msg := range userUniqueMessages {
		// Sanitize userMsg to avoid breaking the code block
		msg = strings.ReplaceAll(msg, "```", "")

		// If sanitizedMsg is longer than 350 characters, truncate it and add three dots
		if len(msg) > 350 {
			msg = msg[:347] + "..."
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Value: "```text\n" + msg + "\n```",
		})
	}

	return embed
}

func (a *Antispam) generateAlertEmbedDescription() string {
	return fmt.Sprintf(
		"Member exceeded channel limit `%d` within `%d` seconds."+"\n"+
			"Added to the denied list for the next `%d` seconds.",
		a.cfg.MaxChannelsPerUser,
		a.cfg.MessageTTLSec,
		a.cfg.DenyTTLSec,
	)
}
