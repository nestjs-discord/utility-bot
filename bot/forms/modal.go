package forms

import (
	"errors"
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/samber/lo"
	"strings"
	"time"
)

func (f *Forms) isTimeWithinDays(t time.Time, days int) bool {
	if t.IsZero() {
		return false
	}
	threshold := time.Now().AddDate(0, 0, -days)
	return t.After(threshold)
}

func (f *Forms) OpenModalButtonClicked(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
	formId := customId.FormId

	form, err := f.getFormById(formId)
	if err != nil {
		return err
	}

	accCreatedAt, err := dgo.SnowflakeTimestamp(i.Member.User.ID)
	if err != nil {
		return err
	}

	// Minimum age of the user's Discord account (in days) required open the modal
	minimumAccountAgeDays := form.MinimumAccountAgeDays

	// Minimum number of days the user must have been a member of the server to open the modal
	minimumServerJoinDays := form.MinimumServerJoinDays

	if f.isTimeWithinDays(accCreatedAt, minimumAccountAgeDays) {
		respond.InteractionWithEphemeralMessage(s, i,
			":warning: This action cannot be performed as your Discord account is too new.",
		)
		return nil
	}

	if f.isTimeWithinDays(i.Member.JoinedAt, minimumServerJoinDays) {
		respond.InteractionWithEphemeralMessage(s, i,
			fmt.Sprintf(
				":warning: In order to perform this action, you must have been a member of this server for at least %d days.",
				minimumServerJoinDays,
			),
		)
		return nil
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
			Title:      form.ModalTitle,
			Components: msgComponents,
		},
	})
	if err != nil {
		return fmt.Errorf("error sending modal interaction response: %w", err)
	}
	return nil
}

func (f *Forms) ModalSubmitted(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
	form, err := f.getFormById(customId.FormId)
	if err != nil {
		return err
	}

	data := i.ModalSubmitData()

	// map of the 'input id' to the 'user given value'
	var userInputs []userInput
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

			userInputs = append(userInputs, userInput{
				InputId: child.CustomID,
				Value:   val,
			})
		}
	}

	// safety check, in case Discord updated its response
	if len(userInputs) == 0 {
		return errors.New("failed to extract the modal values! please report this issue")
	}

	// generate embed
	discordProfileEmbed := &dgo.MessageEmbed{
		// Title: "Discord profile",
		Title: i.Member.User.GlobalName,
		// Color: i.Member.User.AccentColor,
		Color: form.Color,
		Thumbnail: &dgo.MessageEmbedThumbnail{
			URL: i.Member.User.AvatarURL("4096"),
		},
		Fields: []*dgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  "`" + i.Member.User.String() + "`",
				Inline: true,
			},
			{
				Name:   "Profile",
				Value:  i.Member.User.Mention(),
				Inline: false,
			},
		},
	}

	if authorAccCreatedAt, err := dgo.SnowflakeTimestamp(i.Member.User.ID); err == nil {
		discordProfileEmbed.Fields = append(discordProfileEmbed.Fields, &dgo.MessageEmbedField{
			Name:   "Account created",
			Value:  fmt.Sprintf("<t:%d:R>", authorAccCreatedAt.UTC().Unix()),
			Inline: true,
		})
	}

	if !i.Member.JoinedAt.IsZero() {
		discordProfileEmbed.Fields = append(discordProfileEmbed.Fields, &dgo.MessageEmbedField{
			Name:   "Joined the server",
			Value:  fmt.Sprintf("<t:%d:R>", i.Member.JoinedAt.UTC().Unix()),
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
	for _, inp := range userInputs {
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

		modComps, err := f.generateModComponents(customId.FormId, i.Member.User.ID)
		if err != nil {
			return fmt.Errorf("failed to generate mod components: %s", err)
		}

		message.Components = append(message.Components, modComps)
	}

	sentMessage, err := s.ChannelMessageSendComplex(channelId, message)
	if err != nil {
		return fmt.Errorf("failed to send form embed: %s", err)
	}

	// publish the message
	if form.ModSkipApproval {
		_, _ = s.ChannelMessageCrosspost(sentMessage.ChannelID, sentMessage.ID)
	}

	respond.InteractionWithEphemeralMessage(s, i, respContent)

	return nil
}
