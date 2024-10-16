package forms

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
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
	customId, err := components.EncodeCustomId(&components.CustomID{
		Action: OpenModalButton,
		FormId: formId,
	})
	if err != nil {
		return err
	}

	button := discordgo.Button{
		Label:    btnLabel,
		Style:    discordgo.SuccessButton,
		Disabled: false,
		CustomID: customId,
	}

	messageData := &discordgo.MessageSend{
		Content: "", // TODO: add some set of rules (can used embed too)
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					button,
				},
			},
		},
	}

	_, err = session.ChannelMessageSendComplex(channelId, messageData)
	return err
}
