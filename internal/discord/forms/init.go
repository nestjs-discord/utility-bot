package forms

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/config"
)

// Init ensure the form channels always have interactive buttons as the last message.
func Init(session *discordgo.Session, forms map[string]config.Form) error {
	limit := 1
	beforeId := ""
	afterId := ""
	aroundId := ""

	for formId, form := range forms {
		messages, err := session.ChannelMessages(form.ChannelId, limit, beforeId, afterId, aroundId)
		if err != nil {
			return err
		}

		if len(messages) != 0 && DoesHaveButtonComponentWithLabel(messages[0], form.ButtonLabel) {
			continue
		}

		label := form.ButtonLabel
		if err = SendFormButton(session, form.ChannelId, formId, label); err != nil {
			return err
		}
	}

	return nil
}

func DoesHaveButtonComponentWithLabel(msg *discordgo.Message, buttonLabel string) bool {
	// if the message does not have any component
	if len(msg.Components) == 0 {
		return false
	}

	actionsRow, ok := msg.Components[0].(*discordgo.ActionsRow)
	if !ok {
		return false
	}

	if len(actionsRow.Components) == 0 {
		return false
	}

	btn, ok := actionsRow.Components[0].(*discordgo.Button)
	if !ok {
		return false
	}
	if btn.Label != buttonLabel {
		return false
	}

	return true
}

func SendFormButton(session *discordgo.Session, channelId string, formId string, btnLabel string) error {
	messageData := &discordgo.MessageSend{
		Content: "",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    btnLabel,
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: FormButtonIdPrefix + formId,
					},
				},
			},
		},
	}

	_, err := session.ChannelMessageSendComplex(channelId, messageData)
	return err
}
