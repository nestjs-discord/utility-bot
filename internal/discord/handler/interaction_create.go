package handler

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/infra/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/samber/lo"
	"strings"
)

func handleInteractionModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	modalCustomId := data.CustomID
	if !strings.HasPrefix(modalCustomId, forms.FormModalIdPrefix) {
		return // skip it
	}

	formId := strings.TrimPrefix(modalCustomId, forms.FormModalIdPrefix)
	form, ok := config.Yaml().Forms[formId]
	if !ok {
		return // skip invalid forms
	}

	// map of the 'input id' to the 'user given value'
	var userInput []forms.UserInput
	for _, parentComp := range data.Components {
		switch row := parentComp.(type) {
		case *discordgo.ActionsRow:
			for _, childComp := range row.Components {
				switch child := childComp.(type) {
				case *discordgo.TextInput:
					val := strings.TrimSpace(child.Value)       // basic space trim
					val = strings.ReplaceAll(val, "\n\n", "\n") // remove double next lines
					val = strings.ReplaceAll(val, "\t", " ")    // replace the tab character
					val = strings.ReplaceAll(val, "  ", " ")    // remove double spaces
					val = markdown.ConvertLinksToHyperlinks(val)

					if val == "" {
						continue
					}

					userInput = append(userInput, forms.UserInput{
						InputId: child.CustomID,
						Value:   val,
					})
				}
			}
		}
	}

	// safety check, in case Discord updated its response
	if len(userInput) == 0 {
		util.InteractionRespondError(errors.New("failed to extract the modal values! please report this issue"), s, i)
		return
	}

	// generate embed
	embed := &discordgo.MessageEmbed{
		Color: form.Color,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    form.Footer,
			IconURL: i.Message.Author.AvatarURL(""), // the interaction author is the bot itself!
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: i.Member.User.AvatarURL("100"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Discord Name",
				Value:  i.Member.User.GlobalName,
				Inline: true,
			},
			{
				Name:   "Username",
				Value:  "`" + util.FormatUsername(i.Member.User) + "`",
				Inline: true,
			},
		},
	}

	if authorAccCreatedAt, err := discordgo.SnowflakeTimestamp(i.Member.User.ID); err == nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Account Created",
			//Value: humanize.Time(authorAccCreatedAt),
			Value:  fmt.Sprintf("<t:%d:R>", authorAccCreatedAt.UTC().Unix()),
			Inline: false,
		})
	}

	// append user inputs
	for _, inp := range userInput {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   lo.Capitalize(inp.InputId),
			Value:  inp.Value,
			Inline: false,
		})
	}

	channelId := form.ChannelId
	message := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	}
	respContent := "Thank you for taking your time to fill this form. ✅"

	if !form.ModSkipApproval {
		channelId = form.ModChannelId // overwrite the public channel with the private one.

		// the user who fills the modal should know their request is going to be in a pending state.
		respContent += "\n\nModerators will review your request shortly. 🔎"

		message.Components = append(message.Components, forms.GenerateModComponents(formId, i.Member.User.ID))
	}

	sentMessage, err := s.ChannelMessageSendComplex(channelId, message)
	if err != nil {
		util.InteractionRespondError(fmt.Errorf("failed to send form embed: %s", err), s, i)
		return
	}

	// publish the message
	if form.ModSkipApproval {
		_, _ = s.ChannelMessageCrosspost(sentMessage.ChannelID, sentMessage.ID)
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: respContent,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
