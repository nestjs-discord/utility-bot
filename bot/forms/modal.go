package forms

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
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
