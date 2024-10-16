package forms

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"time"
)

func (f *Forms) ModAcceptButtonClicked(s *discordgo.Session, i *discordgo.InteractionCreate, customId *components.CustomID) error {
	if f.RaceConditionCheck(i.Message.ID) {
		f.RaceConditionRespond(s, i)
		return nil
	}

	formId := customId.FormId
	form, err := f.getFormById(formId)
	if err != nil {
		return err
	}

	// receive the last message in the public channel
	messages, err := s.ChannelMessages(form.ChannelId, 1, "", "", "")
	if err != nil {
		return fmt.Errorf("error getting channel messages: %v", err)
	}

	// Send the embed data into the public channel
	sentMessage, err := s.ChannelMessageSendComplex(form.ChannelId, &discordgo.MessageSend{
		Embeds: i.Message.Embeds,
	})
	if err != nil {
		return fmt.Errorf("failed to send the embeded message into the public channel: %s", err)
	}

	// send the interactive form button again (since we deleted the last one)
	err = f.sendFormButton(s, form.ChannelId, formId, form.ButtonLabel)
	if err != nil {
		return fmt.Errorf("failed to send the interactive form button again (after deleting): %s", err)
	}

	// at the point, since we know a new interactive form button is sent into the public channel
	// so it is safe to delete the old message that has the interactive form button
	if len(messages) == 1 && f.doesHaveButtonComponentWithLabel(messages[0], form.ButtonLabel) {
		_ = s.ChannelMessageDelete(form.ChannelId, messages[0].ID)
	}

	content := fmt.Sprintf("Accepted by %s, <t:%d:R>",
		i.Member.User.Mention(),
		time.Now().UTC().Unix(),
	)

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
		return fmt.Errorf("failed to edit the current message (after all the steps): %s", err)
	}

	return nil
}
