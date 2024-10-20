package forms

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"time"
)

func (f *Forms) isTimeWithinDays(t time.Time, days int) bool {
	if t.IsZero() {
		return false
	}
	threshold := time.Now().AddDate(0, 0, -days)
	return t.After(threshold)
}

func (f *Forms) ModalOpenButtonClicked(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
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

		label := input.Label

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
