package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/infra/config"
	"github.com/samber/lo"
	"log/slog"
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

	if data.CustomID == "" { // TODO: is this check necessary?
		return
	}

	customId, err := components.DecodeCustomId(data.CustomID)
	if err != nil {
		h.logger.Error("failed to decode custom id",
			slog.String("customId", data.CustomID),
			slog.Any("interaction", i),
		)
		return
	}

	switch customId.Action {
	case forms.ModeratorAcceptButton:
		err = h.forms.ModAcceptButtonClicked(s, i, customId)
		if err != nil {
			h.RespondError(err, s, i)
		}
		return

		// ...
	}

	// Mod -> Reject button
	if strings.HasPrefix(data.CustomID, forms.ModeratorRejectButton) {
		//formId := strings.TrimPrefix(data.CustomID, forms.ModeratorRejectButton)

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
	if strings.HasPrefix(data.CustomID, forms.ModeratorBanButton) {
		userIdToBan := strings.TrimPrefix(data.CustomID, forms.ModeratorBanButton)

		banReason := fmt.Sprintf("Banned by %s (%s)",
			i.Member.User.GlobalName,
			i.Member.User.Username,
		)

		err = s.GuildBanCreateWithReason(i.GuildID, userIdToBan, banReason, 7)
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
	if !strings.HasPrefix(data.CustomID, forms.OpenModalButton) {
		return // not a form interactive button, skip it
	}

	inputFormId := strings.TrimPrefix(data.CustomID, forms.OpenModalButton)

	form, ok := config.Yaml().Forms[inputFormId]
	if !ok {
		return // form does not exist in the YAML config
	}

	var msgComponents []discordgo.MessageComponent
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

		msgComponents = append(msgComponents, actionsRow)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			//CustomID:   forms.Modal + i.Interaction.Member.User.ID,
			CustomID:   forms.Modal + inputFormId,
			Title:      form.Title,
			Components: msgComponents,
		},
	})
	if err != nil {
		h.RespondError(err, s, i)
	}
}
