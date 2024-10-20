package forms

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
)

func (f *Forms) doesHaveButtonComponentWithLabel(msg *dgo.Message, openModalMessage yaml.FormOpenModalMessage) bool {
	// if the message does not have any component
	if len(msg.Components) == 0 {
		return false
	}

	actionsRow, ok := msg.Components[0].(*dgo.ActionsRow)
	if !ok {
		return false
	}

	if len(actionsRow.Components) == 0 {
		return false
	}

	btn, ok := actionsRow.Components[0].(*dgo.Button)
	if !ok {
		return false
	}
	if btn.Label != openModalMessage.ButtonLabel {
		return false
	}

	return true
}

func (f *Forms) sendOpenModalMessage(channelId string, formId string, openModalMessage yaml.FormOpenModalMessage) error {
	customId, err := components.EncodeCustomId(&components.CustomID{
		Action: OpenModalButton,
		FormId: formId,
	})
	if err != nil {
		return err
	}

	button := dgo.Button{
		Label:    openModalMessage.ButtonLabel,
		Style:    dgo.SuccessButton,
		Disabled: false,
		CustomID: customId,
	}

	embed := &dgo.MessageEmbed{
		Title:       openModalMessage.EmbedTitle,
		Color:       openModalMessage.EmbedColor,
		Description: openModalMessage.EmbedDescription,
	}

	row := dgo.ActionsRow{
		Components: []dgo.MessageComponent{
			button,
		},
	}

	messageData := &dgo.MessageSend{
		Content:    "‎", // empty character to space out the previous message https://emptycharacter.com/
		Embeds:     []*dgo.MessageEmbed{embed},
		Components: []dgo.MessageComponent{row},
	}

	_, err = f.opts.Session.ChannelMessageSendComplex(channelId, messageData)

	return err
}
