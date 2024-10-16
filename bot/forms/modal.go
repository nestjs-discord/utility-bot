package forms

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/markdown"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func (f *Forms) OpenModalButtonClicked(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
	formId := customId.FormId

	form, ok := f.cfg[formId]
	if !ok {
		return fmt.Errorf("form not found")
	}

	var msgComponents []dgo.MessageComponent
	for _, input := range form.Inputs {
		style := dgo.TextInputShort
		if input.Multiline {
			style = dgo.TextInputParagraph
		}

		label := lo.Capitalize(input.Id)

		comp := dgo.TextInput{
			CustomID:    input.Id,
			Label:       label,
			Style:       style,
			Placeholder: input.Placeholder,
			Required:    input.Required,
			MaxLength:   input.Max,
			MinLength:   input.Min,
		}

		actionsRow := dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				comp,
			},
		}

		msgComponents = append(msgComponents, actionsRow)
	}

	formCustomId, err := components.EncodeCustomId(&components.CustomID{
		Action: Modal,
		FormId: formId,
	})
	if err != nil {
		return err
	}

	err = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseModal,
		Data: &dgo.InteractionResponseData{
			CustomID:   formCustomId,
			Title:      form.Title,
			Components: msgComponents,
		},
	})
	if err != nil {
		return fmt.Errorf("error sending modal interaction response: %w", err)
	}
	return nil
}

func (f *Forms) ModalSubmitted(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
	formId := customId.FormId

	form, ok := f.cfg[formId]
	if !ok {
		return fmt.Errorf("form not found")
	}

	data := i.ModalSubmitData()

	// map of the 'input id' to the 'user given value'
	var userInput []UserInput
	for _, parentComp := range data.Components {
		row, isRow := parentComp.(*dgo.ActionsRow)
		if !isRow {
			break
		}

		for _, childComp := range row.Components {
			child, isText := childComp.(*dgo.TextInput)
			if !isText {
				break
			}

			val := strings.TrimSpace(child.Value)       // basic space trim
			val = strings.ReplaceAll(val, "\n\n", "\n") // remove double next lines
			val = strings.ReplaceAll(val, "\t", " ")    // replace the tab character
			val = strings.ReplaceAll(val, "  ", " ")    // remove double spaces
			val = markdown.ConvertLinksToHyperlinks(val)

			if val == "" {
				continue
			}

			userInput = append(userInput, UserInput{
				InputId: child.CustomID,
				Value:   val,
			})
		}
	}

	// safety check, in case Discord updated its response
	if len(userInput) == 0 {
		return errors.New("failed to extract the modal values! please report this issue")
	}

	// generate embed
	discordProfileEmbed := &dgo.MessageEmbed{
		Title: "Discord profile",
		Color: i.Member.User.AccentColor,
		Thumbnail: &dgo.MessageEmbedThumbnail{
			URL: i.Member.User.AvatarURL("100"),
		},
		Fields: []*dgo.MessageEmbedField{
			{
				Name:   "Account",
				Value:  i.Member.User.Mention(),
				Inline: true,
			},
			{
				Name:   "Name",
				Value:  i.Member.User.GlobalName,
				Inline: true,
			},
			{
				Name:   "Username",
				Value:  "`" + i.Member.User.String() + "`",
				Inline: false,
			},
		},
	}

	if authorAccCreatedAt, err := dgo.SnowflakeTimestamp(i.Member.User.ID); err == nil {
		discordProfileEmbed.Fields = append(discordProfileEmbed.Fields, &dgo.MessageEmbedField{
			Name:   "Created",
			Value:  fmt.Sprintf("<t:%d:R>", authorAccCreatedAt.UTC().Unix()),
			Inline: true,
		})
	}

	formDataEmbed := &dgo.MessageEmbed{
		Color: form.Color,
		Footer: &dgo.MessageEmbedFooter{
			Text:    form.Footer,
			IconURL: i.Message.Author.AvatarURL(""), // the interaction author is the bot itself!
		},
	}

	// append user inputs
	for _, inp := range userInput {
		formDataEmbed.Fields = append(formDataEmbed.Fields, &dgo.MessageEmbedField{
			Name:   lo.Capitalize(inp.InputId),
			Value:  inp.Value,
			Inline: false,
		})
	}

	channelId := form.ChannelId
	message := &dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			discordProfileEmbed,
			formDataEmbed,
		},
	}

	respContent := "Thank you for taking your time to fill this form. ✅"

	if !form.ModSkipApproval {
		channelId = form.ModChannelId // overwrite the public channel with the private one.

		// the user who fills the modal must know their request is going to be in a pending state.
		respContent += "\n\nModerators will review your request shortly. 🔎"

		message.Components = append(message.Components, GenerateModComponents(formId, i.Member.User.ID))
	}

	sentMessage, err := s.ChannelMessageSendComplex(channelId, message)
	if err != nil {
		return fmt.Errorf("failed to send form embed: %s", err)
	}

	// publish the message
	if form.ModSkipApproval {
		_, _ = s.ChannelMessageCrosspost(sentMessage.ChannelID, sentMessage.ID)
	}

	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: respContent,
			Flags:   dgo.MessageFlagsEphemeral,
		},
	})

	return nil
}
