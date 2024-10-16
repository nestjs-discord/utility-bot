package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/infra/config"
	"github.com/samber/lo"
	"strings"
	"time"
)

func (h *Handler) MessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: refactor this, it should only be called for the moderator actions, not the public interactive button!
	if h.forms.RaceConditionCheck(i.Message.ID) {
		h.forms.RaceConditionRespond(s, i)
		return
	}

	data := i.MessageComponentData()

	// Mod -> Accept button
	if strings.HasPrefix(data.CustomID, forms.ModAcceptBtnIdPrefix) {
		formId := strings.TrimPrefix(data.CustomID, forms.ModAcceptBtnIdPrefix)

		form, ok := config.Yaml().Forms[formId]
		if !ok {
			return
		}

		// receive the last message in the public channel
		messages, err := s.ChannelMessages(form.ChannelId, 1, "", "", "")
		if err != nil {
			h.RespondError(err, s, i)
			return
		}

		// Send the embed data into the public channel
		sentMessage, err := s.ChannelMessageSendComplex(form.ChannelId, &discordgo.MessageSend{
			Embeds: i.Message.Embeds,
		})
		if err != nil {
			msg := fmt.Errorf("failed to send the embeded message into the public channel: %s", err)
			h.RespondError(msg, s, i)
			return
		}

		// send the interactive form button again (since we deleted the last one)
		err = forms.SendFormButton(s, form.ChannelId, formId, form.ButtonLabel)
		if err != nil {
			msg := fmt.Errorf("failed to send the interactive form button again (after deleting): %s", err)
			h.RespondError(msg, s, i)
			return
		}

		// at the point, since we know a new interactive form button is sent into the public channel
		// so it is safe to delete the old message that has the interactive form button
		if len(messages) == 1 && forms.DoesHaveButtonComponentWithLabel(messages[0], form.ButtonLabel) {
			_ = s.ChannelMessageDelete(form.ChannelId, messages[0].ID)
		}

		userId := i.Member.User.ID
		content := fmt.Sprintf("Accepted by <@%s>, <t:%d:R>", userId, time.Now().UTC().Unix())

		_, err = s.ChannelMessageCrosspost(sentMessage.ChannelID, sentMessage.ID)
		if err != nil {
			content += "\nCross posting the message to the followers failed ❌: " + err.Error()
		} else {
			content += "\nCross posted the message to the followers ✅"
		}

		msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
		msgEdit.SetContent(content)

		// remove the message components
		emptyComponent := make([]discordgo.MessageComponent, 0)
		msgEdit.Components = &emptyComponent

		_, err = s.ChannelMessageEditComplex(msgEdit)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return
	}

	// Mod -> Reject button
	if strings.HasPrefix(data.CustomID, forms.ModRejectBtnIdPrefix) {
		//formId := strings.TrimPrefix(data.CustomID, forms.ModRejectBtnIdPrefix)

		msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
		userId := i.Member.User.ID
		content := fmt.Sprintf("Rejected by <@%s>, <t:%d:R>\n", userId, time.Now().UTC().Unix())
		msgEdit.SetContent(content)

		// remove the message components
		emptyComponent := make([]discordgo.MessageComponent, 0)
		msgEdit.Components = &emptyComponent

		_, err := s.ChannelMessageEditComplex(msgEdit)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return
	}

	// Mod -> Ban button
	if strings.HasPrefix(data.CustomID, forms.ModBanBtnIdPrefix) {
		userIdToBan := strings.TrimPrefix(data.CustomID, forms.ModBanBtnIdPrefix)

		banReason := fmt.Sprintf("Banned by %s (%s)",
			i.Member.User.GlobalName,
			i.Member.User.Username,
		)

		err := s.GuildBanCreateWithReason(i.GuildID, userIdToBan, banReason, 7)
		if err != nil {
			h.RespondError(fmt.Errorf("failed to ban the given user: %s", err), s, i)
			return
		}

		msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
		userId := i.Member.User.ID
		content := fmt.Sprintf("Banned by <@%s>, <t:%d:R>\n", userId, time.Now().UTC().Unix())
		msgEdit.SetContent(content)

		// remove the message components
		emptyComponent := make([]discordgo.MessageComponent, 0)
		msgEdit.Components = &emptyComponent

		_, err = s.ChannelMessageEditComplex(msgEdit)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return
	}

	// Form interactive button (to open modal)
	if !strings.HasPrefix(data.CustomID, forms.FormButtonIdPrefix) {
		return // not a form interactive button, skip it
	}

	inputFormId := strings.TrimPrefix(data.CustomID, forms.FormButtonIdPrefix)

	form, ok := config.Yaml().Forms[inputFormId]
	if !ok {
		return // form does not exist in the YAML config
	}

	var components []discordgo.MessageComponent
	for _, input := range form.Inputs {
		style := discordgo.TextInputShort
		if input.Multiline {
			style = discordgo.TextInputParagraph
		}

		label := lo.Capitalize(input.Id)

		comp := discordgo.TextInput{
			CustomID:    input.Id,
			Label:       label,
			Style:       style,
			Placeholder: input.Placeholder,
			Required:    input.Required,
			MaxLength:   input.Max,
			MinLength:   input.Min,
		}

		actionsRow := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				comp,
			},
		}

		components = append(components, actionsRow)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			//CustomID:   forms.FormModalIdPrefix + i.Interaction.Member.User.ID,
			CustomID:   forms.FormModalIdPrefix + inputFormId,
			Title:      form.Title,
			Components: components,
		},
	})
	if err != nil {
		h.RespondError(err, s, i)
	}
}
