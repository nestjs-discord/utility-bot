package forms

import (
	"github.com/bwmarrin/discordgo"
)

func (f *Forms) doesHaveButtonComponentWithLabel(msg *discordgo.Message, buttonLabel string) bool {
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

func (f *Forms) sendFormButton(session *discordgo.Session, channelId string, formId string, btnLabel string) error {
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
