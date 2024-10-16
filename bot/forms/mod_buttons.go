package forms

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
)

func (f *Forms) generateModComponents(formId string, userId string) (*dgo.ActionsRow, error) {
	acceptBtn, err := f.generateAcceptButton(formId)
	if err != nil {
		return nil, err
	}

	rejectBtn, err := f.generateRejectButton(formId)
	if err != nil {
		return nil, err
	}

	banBtn, err := f.generateBanButton(userId)
	if err != nil {
		return nil, err
	}

	row := &dgo.ActionsRow{
		Components: []dgo.MessageComponent{
			acceptBtn,
			rejectBtn,
			banBtn,
		},
	}

	return row, nil
}

func (f *Forms) generateAcceptButton(formId string) (*dgo.Button, error) {
	customId, err := components.EncodeCustomId(&components.CustomID{
		Action: ModeratorAcceptButton,
		FormId: formId,
	})
	if err != nil {
		return nil, err
	}

	btn := &dgo.Button{
		Label:    "Accept & Publish",
		Style:    dgo.SecondaryButton,
		Disabled: false,
		CustomID: customId,
		Emoji:    &dgo.ComponentEmoji{Name: "✅", Animated: false},
	}
	return btn, nil
}

func (f *Forms) generateRejectButton(formId string) (*dgo.Button, error) {
	customId, err := components.EncodeCustomId(&components.CustomID{
		Action: ModeratorRejectButton,
		FormId: formId,
	})
	if err != nil {
		return nil, err
	}

	btn := &dgo.Button{
		Label:    "Reject",
		Style:    dgo.SecondaryButton,
		Disabled: false,
		CustomID: customId,
		Emoji:    &dgo.ComponentEmoji{Name: "❌", Animated: false},
	}
	return btn, nil
}

func (f *Forms) generateBanButton(userId string) (*dgo.Button, error) {
	customId, err := components.EncodeCustomId(&components.CustomID{
		Action: ModeratorBanButton,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	btn := &dgo.Button{
		Label:    "Ban the user",
		Style:    dgo.SecondaryButton,
		Disabled: false,
		CustomID: customId,
		Emoji:    &dgo.ComponentEmoji{Name: "🔴", Animated: false},
	}
	return btn, nil
}
