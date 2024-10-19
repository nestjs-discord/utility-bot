package forms

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
)

func (f *Forms) ModBanButtonClicked(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
	encodedCustomId, err := components.EncodeCustomId(&components.CustomID{
		Action: ModeratorBanConfirmModal,
		UserId: customId.UserId,
		FormId: customId.FormId,
	})
	if err != nil {
		return fmt.Errorf("encode custom id failed: %w", err)
	}

	comp := dgo.TextInput{
		CustomID:    encodedCustomId,
		Label:       "Are you sure?",
		Style:       dgo.TextInputParagraph,
		Placeholder: "Type 'yes' to confirm this action.",
		Required:    true,
		MaxLength:   3,
		MinLength:   3,
	}

	actionsRow := dgo.ActionsRow{
		Components: []dgo.MessageComponent{
			comp,
		},
	}

	err = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseModal,
		Data: &dgo.InteractionResponseData{
			CustomID: encodedCustomId,
			Title:    "Confirmation",
			Components: []dgo.MessageComponent{
				actionsRow,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error sending modal interaction response: %w", err)
	}

	return nil
}
